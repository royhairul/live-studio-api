package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/services/permission"
)

func PermissionIndex(c *gin.Context) {
	permissions, err := permission.GetAllPermission()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get permissions",
			"error":   err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Permissions retrieved successfully",
		"data":    permissions,
	})
}

func PermissionIndexGrouped(c *gin.Context) {
	permissions, err := permission.GetAllPermissionGrouped()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get permissions",
			"error":   err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Permissions retrieved successfully",
		"data":    permissions,
	})
}
