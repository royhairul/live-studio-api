package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/models"
	"github.com/royhairul/live-studio-api/services/schedule"
)

func HostScheduleIndex(c *gin.Context) {
	data, err := schedule.GetHostScheduleAll()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to retrieve host schedules",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Host schedules retrieved successfully",
		"data":    data,
		"total":   len(data),
	})
}

func HostScheduleCreate(c *gin.Context) {
	var request dto.CreateHostScheduleDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	_, err := schedule.CreateHostSchedule(&request)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to create host schedule",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Host schedule created successfully",
	})
}

func HostScheduleShow(c *gin.Context) {
	id := c.Param("id")
	data, err := schedule.GetHostScheduleByID(id)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Host schedule not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Host schedule retrieved successfully",
		"data":    data,
	})
}

func HostScheduleUpdate(c *gin.Context) {
	var request dto.UpdateHostScheduleDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	id := c.Param("id")
	err := schedule.UpdateHostSchedule(id, &request)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to update host schedule",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Host schedule updated successfully",
	})
}

func HostScheduleDelete(c *gin.Context) {
	id := c.Param("id")
	err := schedule.DeleteHostSchedule(id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to delete host schedule",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Host schedule deleted successfully",
	})
}

func HostScheduleByHostID(c *gin.Context) {
	id := c.Param("id")
	data, err := schedule.GetHostScheduleByHostID(id)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Host schedule not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Host schedule retrieved successfully",
		"data":    data,
	})
}

func HostScheduleScheduled(ctx *gin.Context) {
	dateStr := ctx.Query("date")
	shiftID := ctx.Query("shift_id")

	if dateStr == "" || shiftID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "date and shift_id are required"})
		return
	}

	// Convert date string to time.Time if necessary
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	var schedules []models.Schedule
	err = database.DB.Preload("Host").Preload("Host.Studio").
		Where("date = ? AND shift_id = ?", date, shiftID).
		Find(&schedules).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, schedule := range schedules {
		result = append(result, gin.H{
			"host_id":     schedule.Host.ID,
			"name":        schedule.Host.Name,
			"studio_name": schedule.Host.Studio.Name,
		})
	}

	if len(result) == 0 {
		result = make([]gin.H, 0)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Scheduled hosts retrieved successfully",
		"shift":   shiftID,
		"date":    dateStr,
		"data":    result,
	})
}

func HostScheduleSwitch(c *gin.Context) {
	var request dto.SwitchHostScheduleDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	err := schedule.SwitchHostSchedule(&request)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to switch host schedule",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Host schedule switched successfully",
	})
}
