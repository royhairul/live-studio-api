package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/services/auth"
	"github.com/royhairul/live-studio-api/services/role"
)

func Login(c *gin.Context) {
	var req dto.LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to login",
			"error":   err.Error(),
		})
		return
	}

	// Check if user exists in the database
	token, err := auth.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data": gin.H{
			"access_token": token,
		},
	})
}

func Register(c *gin.Context) {
	var req dto.RegisterDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
			"error":   err.Error(),
		})
		return
	}

	err := auth.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to register",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User successfully registered",
	})
}

func Me(c *gin.Context) {
	username := c.GetString("username")
	roleName := c.GetString("role")
	name := c.GetString("name")

	id, exists := c.Get("superadmin_id")
	if !exists {
		id, exists = c.Get("user_id")
	}

	role, err := role.GetRoleByName(roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get role",
			"error":   err.Error(),
			"role":    roleName,
		})
		return
	}

	permissions := []string{}
	for _, p := range role.Permissions {
		permissions = append(permissions, p.Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Welcome to the admin",
		"id":          id,
		"name":        name,
		"username":    username,
		"role":        roleName,
		"permissions": permissions,
	})
}

func ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
			"error":   err.Error(),
		})
		return
	}

	// Forgot Password
	err := auth.ForgotPassword(req.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP has been sent to your email",
		"email":   req.Email,
	})
}

func VerifyOtp(c *gin.Context) {
	var req dto.VerifyOtp
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
			"error":   err.Error(),
		})
		return
	}

	// Verify Token
	_, err := auth.VerifyOtp(req.Otp)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP is valid",
	})
}

func ResetPassword(c *gin.Context) {
	var req dto.ResetPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
			"error":   err.Error(),
		})
		return
	}

	err := auth.ResetPassword(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reset Password Succesfully",
	})
}
