package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/royhairul/live-studio-api/internal/domains/live/service"
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
		return true // ⚠️ ubah sesuai origin policy di production
	},
}

// --- Helper: inisialisasi websocket & listener ---
func setupWebSocket(ctx *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, err
	}

	// Keep-alive
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Goroutine untuk handle disconnect
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

// --- Helper: kirim data realtime berkala ---
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

// --- Handler utama: daftar live (realtime) ---
func (l *LiveControllerImpl) GetLive(ctx *gin.Context) {
	conn, err := setupWebSocket(ctx)
	if err != nil {
		log.Println("failed to upgrade websocket:", err)
		return
	}

	// Kirim data awal
	if data, err := l.service.GetLive(); err == nil {
		_ = conn.WriteJSON(data)
	} else {
		log.Println("⚠️ initial live fetch error:", err)
	}

	// Jalankan streaming realtime
	streamRealtime(conn, func() time.Duration { return utils.RandomDuration(2, 10) }, l.service.GetLive)
}

// --- Handler detail live (realtime per sesi) ---
func (l *LiveControllerImpl) GetLiveDetail(ctx *gin.Context) {
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
	if data, err := l.service.GetLiveDetail(accountID, sessionID, productPage, productPageSize); err == nil {
		_ = conn.WriteJSON(data)
	} else {
		log.Println("⚠️ initial live detail fetch error:", err)
	}

	// Jalankan streaming realtime
	streamRealtime(conn, func() time.Duration { return utils.RandomDuration(3, 8) }, func() (any, error) {
		return l.service.GetLiveDetail(accountID, sessionID, productPage, productPageSize)
	})
}
