package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/models"
	"github.com/royhairul/live-studio-api/services/host"
	"github.com/royhairul/live-studio-api/validators"
)

// GET Host
func HostIndex(c *gin.Context) {
	hosts, err := host.GetHostAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch hosts",
			"error":   err.Error(),
		})
		return

	}

	// Check if host list is empty
	if len(*hosts) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "List Host is empty.",
			"data":    hosts,
		})
		return // Prevent sending multiple responses
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List All Host",
		"data":    hosts,
	})
}

// POST Create a Host
func HostCreate(c *gin.Context) {
	var req dto.CreateHostDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"errors":  err.Error(),
		})
		return
	}

	if err := validators.Validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation failed",
			"error":   err.Error(),
		})
		return
	}

	if err := host.CreateHost(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create Host",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Host successfully created",
	})
}

// GET Host by ID
func HostShow(c *gin.Context) {
	var host models.Host
	id := c.Param("id")
	if err := database.DB.Preload("Studio").Where("id = ?", id).First(&host).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Host not found",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Host found",
		"data":    host,
	})
}

func HostUpdate(c *gin.Context) {
	var host models.Host
	id := c.Param("id")
	if err := database.DB.Where("id = ?", id).First(&host).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Host not found",
			"error":   err.Error(),
		})
		return
	}

	var requestUpdate dto.UpdateHostDTO

	if err := c.ShouldBindJSON(&requestUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	database.DB.Model(&host).Updates(requestUpdate)

	c.JSON(http.StatusOK, gin.H{
		"message": "Host updated successfully",
		"data":    host,
	})
}

// DELETE Host by ID
func HostDelete(c *gin.Context) {
	var host models.Host
	id := c.Param("id")
	if err := database.DB.Where("id = ?", id).First(&host).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Host not found",
			"error":   err.Error(),
		})
		return
	}

	if err := database.DB.Delete(&host).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete host",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Host deleted successfully",
	})
}

// Genenrate a schedule for a host
func HostGenerateSchedule(c *gin.Context) {
	var host models.Host
	id := c.Param("id")
	if err := database.DB.Where("id = ?", id).First(&host).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Host not found",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Host found",
		"data":    host,
	})
}

func HostGroupedByStudio(c *gin.Context) {
	var hosts []models.Host

	// Ambil semua host dan preload Studio
	if err := database.DB.Preload("Studio").Find(&hosts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch hosts",
			"error":   err.Error(),
		})
		return
	}

	// Buat map untuk mengelompokkan host berdasarkan nama studio
	grouped := make(map[string][]dto.HostResponse)

	for _, h := range hosts {
		studioName := h.Studio.Name
		hostResp := dto.HostResponse{
			ID:         h.ID,
			Name:       h.Name,
			Phone:      h.Phone,
			StudioName: studioName,
		}
		grouped[studioName] = append(grouped[studioName], hostResp)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List of Hosts Grouped by Studio",
		"data":    grouped,
	})
}
