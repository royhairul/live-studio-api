package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/models"
	"github.com/royhairul/live-studio-api/services/studio"
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
	var req dto.Studio
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create studio",
			"error":   err.Error(),
		})
		return
	}

	studio := models.Studio{
		Name:    req.Name,
		Address: req.Address,
	}

	if err := database.DB.Create(&studio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create studio in database",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Studio successfully created",
		"data":    studio,
	})
}

func StudioShow(c *gin.Context) {
	id := c.Param("id")

	studio, err := studio.GetStudioByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Studio not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Studio details",
		"data":    studio,
	})
}

func StudioUpdate(c *gin.Context) {
	id := c.Param("id")

	var req dto.Studio
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	err := studio.UpdateStudio(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update studio",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Studio successfully updated",
	})
}

func StudioDelete(c *gin.Context) {
	id := c.Param("id")
	var studio models.Studio

	if err := database.DB.First(&studio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Studio not found",
			"error":   err.Error(),
		})
		return
	}

	if err := database.DB.Delete(&studio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete studio",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Studio successfully deleted",
	})
}
