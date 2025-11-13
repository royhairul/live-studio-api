package controller

import "github.com/gin-gonic/gin"

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	ForgotPassword(c *gin.Context)
	VerifyOtp(c *gin.Context)
	ResetPassword(c *gin.Context)
	Me(c *gin.Context)
}
