package ds

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID      uint `json:"user_id"`
	IsModerator bool `json:"is_moderator"`
	jwt.RegisteredClaims
}