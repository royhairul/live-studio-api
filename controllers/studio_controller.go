package controllers

import (
	"live-studio-api/database"
	"live-studio-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET Studio
func StudioIndex(c *gin.Context) {
	var studios []models.Studio
	database.DB.Find(&studios)

	// Check if studio list is empty
	if len(studios) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "List Studio is empty.",
			"data":    studios,
		})
		return // Prevent sending multiple responses
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List All Studio",
		"data":    studios,
	})
}

// POST Create a Studio
func StudioCreate(c *gin.Context) {
	var studio models.Studio
	if err := c.ShouldBindJSON(&studio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create studio",
			"error":   err.Error(),
		})
		return
	}
	database.DB.Create(&studio)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Studio successfully created",
		"data":    studio,
	})
}