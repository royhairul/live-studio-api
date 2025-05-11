package controllers

import (
	"live-studio-api/database"
	"live-studio-api/dto"
	"live-studio-api/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to login",
			"error":   err.Error(),
		})
		return
	}

	// Check if user exists in the database
	if err := database.DB.Where("username = ? AND password = ?", user.Username, user.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    user,
	})
}

func Register(c *gin.Context) {
	var user dto.CreateUserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to register",
			"error":   err.Error(),
		})
		return
	}

	// Check if user already exists in the database
	if err := database.DB.Where("username = ?", user.Username).First(&user).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "User already exists",
			"error":   "Username already taken",
		})
		return
	}

	database.DB.Create(&user)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User successfully registered",
		"data":    user,
	})
}