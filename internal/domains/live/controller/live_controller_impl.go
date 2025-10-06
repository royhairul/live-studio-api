package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/royhairul/live-studio-api/helpers"
	"github.com/royhairul/live-studio-api/helpers/errorhandler"
	"github.com/royhairul/live-studio-api/internal/domains/live/service"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type LiveControllerImpl struct {
	service service.LiveService
}

func NewLiveController(service service.LiveService) LiveController {
	return &LiveControllerImpl{service}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Sesuaikan policy production
	},
}

func (l *LiveControllerImpl) GetLive(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("Websocket client disconnected", err)
				conn.Close()
				return
			}
		}
	}()

	realtimeData, err := l.service.GetLive()
	if err == nil {
		if err := conn.WriteJSON(realtimeData); err != nil {
			log.Println("WebSocket initial write error:", err)
			return
		}
	} else {
		log.Println("Initial realtime fetch error:", err)
	}

	ticker := time.NewTicker(helpers.RandomDuration(2, 10))
	defer ticker.Stop()

	for range ticker.C {
		realtimeData, err := l.service.GetLive()
		if err != nil {
			log.Println("Realtime fetch error:", err)
			continue
		}

		if err := conn.WriteJSON(realtimeData); err != nil {
			log.Println("Websocket write error:", err)
			return
		}
	}
}

// GetLiveDetail implements LiveController.
func (l *LiveControllerImpl) GetLiveDetail(ctx *gin.Context) {
	accountID := ctx.Param("id")
	sessionID := ctx.Query("sessionId")

	if sessionID == "" {
		ctx.JSON(http.StatusBadRequest, response.NewBaseResponse("session ID is required", nil))
		return
	}

	live, err := l.service.GetLiveDetail(accountID, sessionID)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("live founded", live)
	ctx.JSON(http.StatusOK, resp)
}
