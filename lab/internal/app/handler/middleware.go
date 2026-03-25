package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"lab/internal/app/ds"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}
		tokenString := parts[1]

		// Проверяем черный список
		blacklisted, err := h.Redis.IsBlacklisted(c.Request.Context(), tokenString)
		if err != nil {
			h.errorJSON(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		if blacklisted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
			return
		}

		// Парсим и валидируем токен
		claims := &ds.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWT.TokenSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Сохраняем информацию о пользователе в контексте
		c.Set("userID", claims.UserID)
		c.Set("isModerator", claims.IsModerator)
		c.Next()
	}
}

// RequireModerator проверяет, что пользователь модератор
func RequireModerator() gin.HandlerFunc {
	return func(c *gin.Context) {
		isMod, exists := c.Get("isModerator")
		if !exists || !isMod.(bool) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "moderator role required"})
			return
		}
		c.Next()
	}
}