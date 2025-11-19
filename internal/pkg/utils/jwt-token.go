package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/royhairul/live-studio-api/internal/domains/user/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
)

type JWTPayloadDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role uint   `json:"role"`

	tenantdb.TenantBase
	jwt.RegisteredClaims
}

func SignToken(payload *JWTPayloadDTO, SECRET_KEY string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

func GenerateTokenJWT(user *entity.User) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", fmt.Errorf("JWT secret not set")
	}

	accessExpiryStr := os.Getenv("JWT_EXPIRED_AT")
	if accessExpiryStr == "" {
		accessExpiryStr = "24h"
	}

	accessToken, _ := time.ParseDuration(accessExpiryStr)

	payload := JWTPayloadDTO{
		ID:   user.ID.String(),
		Name: user.Name,
		Role: user.RoleID,
		TenantBase: tenantdb.TenantBase{
			TenantID: user.TenantID,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessToken)),
		},
	}

	tokenString, err := SignToken(&payload, secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyTokenJWT(tokenStr string, secret string) (*JWTPayloadDTO, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTPayloadDTO{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTPayloadDTO)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token claims")
	}

	return claims, nil
}
