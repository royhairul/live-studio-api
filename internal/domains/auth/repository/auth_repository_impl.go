package repository

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/auth/entity"
	"github.com/royhairul/live-studio-api/internal/domains/auth/params"
	userEntity "github.com/royhairul/live-studio-api/internal/domains/user/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	DB *gorm.DB
}

// reusable: untuk verifikasi dan mengambil OTP valid
func (a *AuthRepositoryImpl) getValidResetOtp(code string) (*entity.ResetPassword, error) {
	var resetOtp entity.ResetPassword
	if err := a.DB.Where("otp = ?", code).First(&resetOtp).Error; err != nil {
		return nil, errors.New("invalid token")
	}

	if resetOtp.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("your token is expired")
	}

	return &resetOtp, nil
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{DB: db}
}

func (a *AuthRepositoryImpl) ResetPassword(req params.ResetPasswordRequest) (params.ChangePasswordResponse, error) {
	resetOtp, err := a.getValidResetOtp(req.OTP)
	if err != nil {
		return params.ChangePasswordResponse{}, err
	}
	// Validasi password
	if req.Password != req.ConfirmPassword {
		return params.ChangePasswordResponse{}, errors.New("passwords do not match")
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return params.ChangePasswordResponse{}, fmt.Errorf("failed to hash password: %v", err)
	}

	// Ambil user berdasarkan user_id dari resetOtp
	var user userEntity.User
	if err := a.DB.First(&user, "email = ?", resetOtp.Email).Error; err != nil {
		return params.ChangePasswordResponse{}, fmt.Errorf("user not found: %v", err)
	}

	// Update password user
	user.Password = string(hashedPassword)
	if err := a.DB.Save(&user).Error; err != nil {
		return params.ChangePasswordResponse{}, fmt.Errorf("failed to update password: %v", err)
	}

	// Hapus OTP setelah digunakan (opsional tapi disarankan)
	_ = a.DB.Delete(&resetOtp)

	return params.ChangePasswordResponse{
		Message: "password berhasil diubah",
	}, nil
}

func (a *AuthRepositoryImpl) VerifyOtp(req params.VerifyOTPRequest) (params.ChangePasswordResponse, error) {
	_, err := a.getValidResetOtp(req.OTP)
	if err != nil {
		return params.ChangePasswordResponse{}, err
	}

	return params.ChangePasswordResponse{
		Message: "otp valid",
	}, nil
}

func (a *AuthRepositoryImpl) ForgotPassword(req params.ForgotPasswordRequest) (params.ForgotPasswordResponse, error) {
	_ = a.DB.Where("email = ?", req.Email).Delete(&entity.ResetPassword{})

	otp := fmt.Sprintf("%06d", rand.Intn(900000)+100000)

	reset := entity.ResetPassword{
		Email:     req.Email,
		Otp:       otp,
		ExpiredAt: time.Now().Add(20 * time.Minute),
	}

	if err := a.DB.Create(&reset).Error; err != nil {
		return params.ForgotPasswordResponse{}, errors.New("failed to create reset token")
	}

	return params.ForgotPasswordResponse{
		Email: req.Email,
		Otp:   otp,
	}, nil
}
