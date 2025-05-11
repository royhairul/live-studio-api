package controllers

import (
	"live-studio-api/database"
	"live-studio-api/dto"
	"live-studio-api/models"
	"live-studio-api/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET Host
func HostIndex(c *gin.Context) {
	var hosts []models.Host
	database.DB.Find(&hosts)

	// Check if host list is empty
	if len(hosts) == 0 {
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
			"errors":   err.Error(),
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

	// 3. Mapping DTO ke model
	host := models.Host{
		Name:     req.Name,
		Phone:    req.Phone,
		StudioID: req.StudioID,
	}

	// 4. Simpan ke DB dan tangani error
	if err := database.DB.Create(&host).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create host",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Host successfully created",
		"data":    host,
	})
}

// GET Host by ID
func HostShow(c *gin.Context) {
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

//  Genenrate a schedule for a host
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
