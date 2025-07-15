package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/services/account"
)

func AccountIndex(c *gin.Context) {
	accounts, err := account.GetAccountAll()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Failed to get accounts",
			"data":    accounts,
		})
	}

	// Check if akun list is empty
	if len(*accounts) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "List Akun is empty.",
			"data":    accounts,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List All Akun",
		"data":    accounts,
	})
}

// Create or Update, Unique userID
func AccountCreate(c *gin.Context) {
	var req dto.CreateAccountCookiesDTO

	// Get cookies from request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"errors":  err.Error(),
		})
		return
	}

	// Create or Update
	if err := account.CreateOrUpdateAccount(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error to create or update account",
			"errors":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "create or update succesfully",
	})
}
