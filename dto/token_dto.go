package dto

import "github.com/golang-jwt/jwt/v5"

type JWTPayloadDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Valid implements jwt.Claims.
func (j *JWTPayloadDTO) Valid() error {
	panic("unimplemented")
}
