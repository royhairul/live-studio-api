package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/services/schedule"
)

func ScheduleShiftIndex(c *gin.Context) {
	// Get all shifts
	shifts, err := schedule.GetShiftAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve shifts",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shifts retrieved successfully",
		"data":    shifts,
	})
}

func ScheduleShiftCreate(c *gin.Context) {
	var req dto.CreateShiftDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"errors":  err.Error(),
		})
		return
	}

	if err := schedule.CreateShift(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create shift",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Shift created successfully",
	})
}

func ScheduleShiftShow(c *gin.Context) {
	id := c.Param("id")

	shift, err := schedule.GetShiftByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get shift",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shift retrieved successfully",
		"data":    shift,
	})
}

func ScheduleShiftUpdate(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateShiftDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	if err := schedule.UpdateShift(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update shift",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shift updated successfully",
	})
}

func ScheduleShiftDelete(c *gin.Context) {
	id := c.Param("id")

	if err := schedule.DeleteShift(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete shift",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shift deleted successfully",
	})
}
