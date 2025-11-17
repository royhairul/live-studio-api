package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/auth/params"
	"github.com/royhairul/live-studio-api/internal/domains/auth/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type AuthControllerImpl struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) AuthController {
	return &AuthControllerImpl{service: service}
}

func (ctrl *AuthControllerImpl) Login(c *gin.Context) {
	var req params.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	token, err := ctrl.service.Login(req)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewBaseResponse("Login successful", token))
}

func (ctrl *AuthControllerImpl) Register(c *gin.Context) {
	var req params.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	createdUser, err := ctrl.service.Register(req)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.NewBaseResponse("User successfully registered", createdUser))
}

func (ctrl *AuthControllerImpl) ForgotPassword(c *gin.Context) {
	var req params.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	email, err := ctrl.service.ForgotPassword(req)
	log.Println(email)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewBaseResponse("OTP has been sent to your email", email))
}

func (ctrl *AuthControllerImpl) VerifyOtp(c *gin.Context) {
	var req params.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	msg, err := ctrl.service.VerifyOtp(req)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewBaseResponse("Verify OTP successfully", msg.Message))
}

func (ctrl *AuthControllerImpl) ResetPassword(c *gin.Context) {
	var req params.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	msg, err := ctrl.service.ResetPassword(req)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewBaseResponse("Reset Password Successfully", msg.Message))
}

// Me implements AuthController.
func (ctrl *AuthControllerImpl) Me(c *gin.Context) {
	me, err := ctrl.service.Me(c)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewBaseResponse(fmt.Sprintf("Welcome, %s", me.Role), me))
}
