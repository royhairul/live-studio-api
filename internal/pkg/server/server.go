package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/config"
)

func Start(router *gin.Engine, cfg *config.Config) {
	router.Run(fmt.Sprintf(":%s", cfg.ServerPort))
}
