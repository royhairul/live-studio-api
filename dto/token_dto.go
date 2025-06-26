package dto

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTPayloadDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (j *JWTPayloadDTO) Valid() error {
	now := time.Now()

	if j.NotBefore != nil && now.Before(j.NotBefore.Time) {
		return fmt.Errorf("token not active yet")
	}

	if j.ExpiresAt != nil && now.After(j.ExpiresAt.Time) {
		return fmt.Errorf("token expired")
	}

	return nil
}
