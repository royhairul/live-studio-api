package controllers

import (
	"live-studio-api/database"
	"live-studio-api/dto"
	"live-studio-api/models"
	"live-studio-api/services/shopee"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AkunIndex(c *gin.Context) {
	var accounts []models.Akun
	database.DB.Find(&accounts)
	// Check if akun list is empty
	if len(accounts) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "List Akun is empty.",
			"data":    accounts,
		})
		return // Prevent sending multiple responses
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "List All Akun",
		"data":    accounts,
	})
}

func AkunCreate(c *gin.Context) {
	var akun dto.CreateAkunCookiesDTO

	// Get cookies from request body
	if err := c.ShouldBindJSON(&akun); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"errors":  err.Error(),
		})
		return
	}

	// Get Akun from cookies
	res, err := shopee.GetAccount(akun.Cookies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get akun from cookies",
			"errors":  err.Error(),
		})
		return
	}

	account := models.Akun{
		Name:     res["name"].(string),
		Username: res["username"].(string),
		Email:    res["email"].(string),
		Platform: "Shopee",
		Cookies:  akun.Cookies,
	}

	// Save to database
	if err := database.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create host",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List All Akun",
		"data":    account,
	})
}