package auth

import (
	"errors"
	"fmt"
	"live-studio-api/database"
	"live-studio-api/dto"
	"live-studio-api/models"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(userLogin dto.LoginDTO) (string, error) {
	var user models.User
	
	if err := database.DB.Where("email = ?", userLogin.Email).First(&user).Error; err != nil {
		return "", fmt.Errorf("Email not registered.")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		return "", fmt.Errorf("Invalid email or password")
	}

	accessExpiryStr := os.Getenv("JWT_EXPIRED_AT")
	if accessExpiryStr == "" {
		accessExpiryStr = "24h"
	}

	accessToken, _ := time.ParseDuration(accessExpiryStr)

	payload := dto.JWTPayloadDTO{
		ID: user.ID,
		Username: user.Username,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessToken)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := GenerateTokenJWT(payload)
	if err != nil {
		return "", err
	}

	return token,nil
}

func Register(userRegister dto.RegisterDTO) error {
	var existingUser models.User

	if err := database.DB.Where("username = ?", userRegister.Username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("Username already exists.")
	}

	// Melakukan Hashing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)

	newUser := models.User{
		Username: userRegister.Username,
		Password: string(hashedPassword),
		Email: userRegister.Email,
		Role: userRegister.Role,
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		return fmt.Errorf("Failed to create user: %v", err)
	}

	return nil
}

func ForgotPassword(email string) error {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return fmt.Errorf("Email not registered")
	}

	// Hapus Token Lama
	_ = database.DB.Where("email = ?", email).Delete(&models.ResetPassword{})
	
	// Generate OTP
	otp := fmt.Sprintf("%06d", rand.Intn(900000)+100000)
	
	reset := models.ResetPassword{
		Email: email,
		Otp: otp,
		ExpiredAt: time.Now().Add(20 * time.Minute),
	}

	if err := database.DB.Create(&reset).Error; err != nil {
		return fmt.Errorf("Failed to create reset token")
	}
	
	return nil;
}

func VerifyOtp(otp string) (*models.ResetPassword, error) {
	var resetOtp models.ResetPassword
	if err := database.DB.Where("otp = ?", otp).First(&resetOtp).Error; err != nil {
		return nil, fmt.Errorf("Invalid token")
	}
	
	if resetOtp.ExpiredAt.Before(time.Now()) {
		return nil, fmt.Errorf("Your token is expired")
	}

	return &resetOtp, nil
}

func ResetPassword(resetPassword dto.ResetPassword) error {
	resetUser, err := VerifyOtp(resetPassword.Otp);
	
	if err != nil {
		return err
	}

	// Serach user by email
	var user models.User
	if err := database.DB.Where("email = ?", resetUser.Email).First(&user).Error; err != nil {
		return fmt.Errorf("User not found")
	}

	// Validasi
	if resetPassword.Password != resetPassword.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}
	
	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetPassword.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user.Password = string(hashedPassword)
	if err := database.DB.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	return nil
}

func GenerateTokenJWT(payload dto.JWTPayloadDTO) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT secret not set")
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

func VerifyTokenJWT(tokenStr string, secret string) (*dto.JWTPayloadDTO, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &dto.JWTPayloadDTO{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*dto.JWTPayloadDTO)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token claims")
	}

	return claims, nil
}