package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"lab/internal/app/ds"
)

type registerRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type authResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body registerRequest true "User credentials"
// @Success      201 {object} map[string]interface{} "message: user registered"
// @Failure      400 {object} map[string]interface{} "error"
// @Failure      409 {object} map[string]interface{} "error: user already exists"
// @Router       /register [post]
func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Проверяем, существует ли уже пользователь
	_, err := h.Repo.GetUserByLogin(req.Login)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}
	if err := h.Repo.CreateUser(req.Login, req.Password); err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticates user and returns JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body loginRequest true "User credentials"
// @Success      200 {object} authResponse
// @Failure      400 {object} map[string]interface{} "error"
// @Failure      401 {object} map[string]interface{} "error: invalid credentials"
// @Router       /login [post]
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.Repo.GetUserByLogin(req.Login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !h.Repo.CheckPassword(user, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	// Создаём JWT
	expiresIn := h.Config.JWT.ExpiresIn
	now := time.Now()
	claims := ds.JWTClaims{
		UserID:      user.ID,
		IsModerator: user.IsModerator,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "electrolyte-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.Config.JWT.TokenSecret))
	if err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, authResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   int64(expiresIn.Seconds()),
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Adds current JWT token to blacklist
// @Tags         Auth
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "message: logged out"
// @Failure      400 {object} map[string]interface{} "error"
// @Failure      401 {object} map[string]interface{} "error"
// @Router       /logout [post]
func (h *Handler) Logout(c *gin.Context) {
	// Получаем токен из заголовка
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization header"})
		return
	}
	tokenString := authHeader[7:]
	// Добавляем в черный список на оставшееся время жизни
	claims := &ds.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.TokenSecret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		c.JSON(http.StatusOK, gin.H{"message": "already expired"})
		return
	}
	if err := h.Redis.AddToBlacklist(c.Request.Context(), tokenString, ttl); err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}