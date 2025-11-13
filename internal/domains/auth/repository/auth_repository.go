package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/auth/entity"
	"github.com/royhairul/live-studio-api/internal/domains/auth/params"
)

type AuthRepository interface {
	ResetPassword(password params.ResetPasswordRequest) (params.ChangePasswordResponse, error)
	ForgotPassword(email params.ForgotPasswordRequest) (params.ForgotPasswordResponse, error)
	VerifyOtp(otp params.VerifyOTPRequest) (params.ChangePasswordResponse, error)
	getValidResetOtp(code string) (*entity.ResetPassword, error)
}
