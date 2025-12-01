package controller

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/internal/domains/live/service"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
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
	conn *websocket.Conn,
	intervalFunc func() time.Duration,
	fetchFunc func() (T, error),
) {
	defer conn.Close()
	ticker := time.NewTicker(intervalFunc())
	defer ticker.Stop()

	for range ticker.C {
		data, err := fetchFunc()
		if err != nil {
			log.Println("⚠️ realtime fetch error:", err)
			continue
		}

		if err := conn.WriteJSON(data); err != nil {
			log.Println("❌ websocket write error:", err)
			return
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
	streamRealtime(conn, func() time.Duration {
		return utils.RandomDuration(2, 10)
	}, func() (any, error) {
		return l.service.GetLive(ctx.Request.Context())
	})
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
	streamRealtime(conn, func() time.Duration {
		return utils.RandomDuration(3, 8)
	}, func() (any, error) {
		return l.service.GetLiveDetail(ctx.Request.Context(), accountID, sessionID, productPage, productPageSize)
	})
}
