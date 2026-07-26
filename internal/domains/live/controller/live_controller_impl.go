package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/royhairul/live-studio-api/database"
	shopeeparams "github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	"github.com/royhairul/live-studio-api/internal/domains/live/params"
	"github.com/royhairul/live-studio-api/internal/domains/live/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"
	"github.com/royhairul/live-studio-api/internal/pkg/utils"
)

type LiveControllerImpl struct {
	service service.LiveService
}

func NewLiveController(service service.LiveService) LiveController {
	return &LiveControllerImpl{service}
}

// --- WebSocket upgrader ---
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // ⚠️ Sesuaikan di production
	},
}

// ==============================
// 🔹 Helper: Setup WebSocket
// ==============================
func setupWebSocket(ctx *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, err
	}

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Println("🔌 client disconnected:", err)
				conn.Close()
				return
			}
		}
	}()

	return conn, nil
}

// ==============================
// 🔹 Helper streaming realtime
// ==============================
func streamRealtime[T any](
	ctx context.Context,
	conn *websocket.Conn,
	intervalFunc func() time.Duration,
	fetchFunc func() (T, error),
) {
	defer conn.Close()
	ticker := time.NewTicker(intervalFunc())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 realtime stream stopped:", ctx.Err())
			return
		case <-ticker.C:
			data, err := fetchFunc()
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				log.Println("⚠️ realtime fetch error:", err)
				continue
			}

			if err := conn.WriteJSON(data); err != nil {
				log.Println("❌ websocket write error:", err)
				return
			}
		}
	}
}

// =====================================
// 🔹 Handler: List Realtime Live
// =====================================
func (l *LiveControllerImpl) GetLive(ctx *gin.Context) {
	// --- Ambil token dari query ---
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	// --- Validasi token ---
	claims, err := utils.VerifyTokenJWT(token, os.Getenv("JWT_SECRET"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	// ==========================
	// 🔥 Inject TENANT ke context
	// ==========================
	tenantID := claims.TenantID
	if tenantID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "tenant_id missing in token"})
		return
	}

	log.Println("WS Authenticated:", claims.Name, " Tenant:", tenantID)

	// Simpan tenant di Gin context
	ctx.Set("tenant_id", tenantID)

	// Buat request context dengan tenant
	reqCtx := tenantdb.AttachTenant(ctx.Request.Context(), tenantID)
	ctx.Request = ctx.Request.WithContext(reqCtx)

	// Set DB tenant-aware
	ctx.Set("DB", database.DB.WithContext(reqCtx))

	// --- Upgrade ke WebSocket ---
	conn, err := setupWebSocket(ctx)
	if err != nil {
		log.Println("failed to upgrade websocket:", err)
		return
	}

	// --- Kirim data awal ---
	if data, err := l.service.GetLive(ctx.Request.Context()); err == nil {
		_ = conn.WriteJSON(data)
	}

	// --- Kirim realtime ---
	streamRealtime(ctx.Request.Context(), conn, func() time.Duration {
		return utils.RandomDuration(2, 10)
	}, func() (any, error) {
		return l.service.GetLive(ctx.Request.Context())
	})
}

// =====================================
// 🔹 Handler: Riwayat Live (bukan realtime)
// =====================================
// GetStoredHistory serves persisted sessions from the database. Shopee is never
// contacted here — run a sync to refresh the table.
func (l *LiveControllerImpl) GetStoredHistory(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid page", "page must be a positive number"))
		return
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid pageSize", "pageSize must be a positive number"))
		return
	}

	filter := params.LiveFilter{Page: page, PageSize: pageSize}

	// :id is present on /live/history/:id and absent on /live/history.
	if accountID := ctx.Param("id"); accountID != "" {
		filter.AccountID = &accountID
	} else if accountID := ctx.Query("account"); accountID != "" {
		filter.AccountID = &accountID
	}

	if studioID := ctx.Query("studio"); studioID != "" {
		filter.StudioID = &studioID
	}

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")
	if (startDate == "") != (endDate == "") {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid date range", "startDate and endDate must be provided together"))
		return
	}
	if startDate != "" {
		start, err := timehandler.ParseDate(startDate)
		if err != nil {
			errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid startDate", err.Error()))
			return
		}
		end, err := timehandler.ParseDate(endDate)
		if err != nil {
			errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid endDate", err.Error()))
			return
		}
		filter.StartTime, filter.EndTime = start, end
	}

	history, err := l.service.GetStoredHistory(ctx.Request.Context(), filter)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved live history successfully", history)
	ctx.JSON(http.StatusOK, resp)
}

// =====================================
// 🔹 Handler: Detail Sesi Riwayat (REST, sekali GET)
// =====================================
//
// GetStoredHistoryDetail is the plain-REST twin of GetLiveDetail. It returns the
// exact same LiveDetailResponse (overview + product breakdown) but as one JSON
// response instead of a WebSocket stream — used from the history page, where the
// session has already ended and there is nothing left to stream.
//
// Auth and tenant come from the route group middleware (Bearer), unlike the
// WebSocket twin which authenticates via ?token= because browsers cannot set
// headers on a handshake. It still reaches Shopee once, so the product data
// depends on the account cookie being valid.
func (l *LiveControllerImpl) GetStoredHistoryDetail(ctx *gin.Context) {
	accountID := ctx.Param("id")
	sessionID := ctx.Param("sessionId")

	if sessionID == "" {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid session", "session ID is required"))
		return
	}

	productPage := ctx.DefaultQuery("productPage", "1")
	productPageSize := ctx.DefaultQuery("productPageSize", "10")

	detail, err := l.service.GetLiveDetail(ctx.Request.Context(), accountID, sessionID, productPage, productPageSize)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved live detail successfully", detail)
	ctx.JSON(http.StatusOK, resp)
}

// =====================================
// 🔹 Handler: Simpan Riwayat Live ke DB
// =====================================
func (l *LiveControllerImpl) SyncHistory(ctx *gin.Context) {
	req, ok := bindHistoryRequest(ctx)
	if !ok {
		return
	}

	result, err := l.service.SyncHistory(ctx.Request.Context(), ctx.Param("id"), req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("synced live history successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (l *LiveControllerImpl) SyncAllHistory(ctx *gin.Context) {
	req, ok := bindHistoryRequest(ctx)
	if !ok {
		return
	}

	results, err := l.service.SyncAllHistory(ctx.Request.Context(), req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("synced live history for all accounts successfully", results)
	ctx.JSON(http.StatusOK, resp)
}

// SyncHistoryRange backfills one account over a multi-month range by chunking it
// into <=30-day windows — liveList/v2 returns nothing for a timeDim past 30 days.
// The span is given either as startDate+endDate, or as months back from endDate.
func (l *LiveControllerImpl) SyncHistoryRange(ctx *gin.Context) {
	base, ok := bindHistoryRequest(ctx)
	if !ok {
		return
	}

	// endDate is already defaulted to today by bindHistoryRequest.
	end, err := timehandler.ParseDate(base.EndDate)
	if err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid endDate", "endDate must be YYYY-MM-DD"))
		return
	}

	// start = explicit startDate, else <months> before endDate (default 3).
	var start *time.Time
	if s := ctx.Query("startDate"); s != "" {
		start, err = timehandler.ParseDate(s)
		if err != nil {
			errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid startDate", "startDate must be YYYY-MM-DD"))
			return
		}
	} else {
		months, err := strconv.Atoi(ctx.DefaultQuery("months", "3"))
		if err != nil || months < 1 || months > 12 {
			errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid months", "months must be a number between 1 and 12"))
			return
		}
		s := end.AddDate(0, -months, 0)
		start = &s
	}

	if start.After(*end) {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid range", "startDate must be on or before endDate"))
		return
	}

	result, err := l.service.SyncHistoryRange(ctx.Request.Context(), ctx.Param("id"), base, *start, *end)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("synced live history range successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

// bindHistoryRequest parses the liveList/v2 query parameters shared by the
// history and sync handlers. It writes the error response itself and reports
// false when the request is malformed.
func bindHistoryRequest(ctx *gin.Context) (shopeeparams.ShopeeLiveHistoryRequest, bool) {
	var req shopeeparams.ShopeeLiveHistoryRequest

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid page", "page must be a number"))
		return req, false
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid pageSize", "pageSize must be a number"))
		return req, false
	}

	req = shopeeparams.ShopeeLiveHistoryRequest{
		Page:     page,
		PageSize: pageSize,
		Name:     ctx.Query("name"),
		OrderBy:  ctx.Query("orderBy"),
		Sort:     ctx.Query("sort"),
		TimeDim:  ctx.DefaultQuery("timeDim", "1m"),
		EndDate:  ctx.DefaultQuery("endDate", time.Now().Format("2006-01-02")),
	}

	return req, true
}

// =====================================
// 🔹 Handler: Detail Realtime Session
// =====================================
func (l *LiveControllerImpl) GetLiveDetail(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	claims, err := utils.VerifyTokenJWT(token, os.Getenv("JWT_SECRET"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	tenantID := claims.TenantID
	if tenantID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "tenant_id missing in token"})
		return
	}

	log.Println("WS Auth:", claims.Name, "Tenant:", tenantID)

	ctx.Set("tenant_id", tenantID)

	reqCtx := tenantdb.AttachTenant(ctx.Request.Context(), tenantID)
	ctx.Request = ctx.Request.WithContext(reqCtx)
	ctx.Set("DB", database.DB.WithContext(reqCtx))

	accountID := ctx.Param("id")
	sessionID := ctx.Param("sessionId")

	if sessionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "session ID is required"})
		return
	}

	productPage := ctx.DefaultQuery("productPage", "1")
	productPageSize := ctx.DefaultQuery("productPageSize", "10")

	conn, err := setupWebSocket(ctx)
	if err != nil {
		log.Println("failed to upgrade websocket:", err)
		return
	}

	// Kirim data awal
	if data, err := l.service.GetLiveDetail(ctx.Request.Context(), accountID, sessionID, productPage, productPageSize); err == nil {
		_ = conn.WriteJSON(data)
	}

	// Kirim realtime
	streamRealtime(ctx.Request.Context(), conn, func() time.Duration {
		return utils.RandomDuration(3, 8)
	}, func() (any, error) {
		return l.service.GetLiveDetail(ctx.Request.Context(), accountID, sessionID, productPage, productPageSize)
	})
}
