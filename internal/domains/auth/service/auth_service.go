package service

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/auth/params"
)

type AuthService interface {
	Login(input params.LoginRequest) (params.LoginResponse, error)
	Register(input params.RegisterRequest) (params.RegisterResponse, error)
	ForgotPassword(input params.ForgotPasswordRequest) (email string, err error)
	ResetPassword(input params.ResetPasswordRequest) (params.ChangePasswordResponse, error)
	VerifyOtp(input params.VerifyOTPRequest) (params.ChangePasswordResponse, error)
	Me(ctx context.Context) (params.MeResponse, error)
}
