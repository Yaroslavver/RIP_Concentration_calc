package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	"lab/internal/app/config"
	"lab/internal/app/redis"
	"lab/internal/app/repository"
)

type Handler struct {
	Repo   *repository.Repository
	Config *config.Config
	Redis  *redis.Client
}

func NewHandler(repo *repository.Repository, cfg *config.Config, r *redis.Client) *Handler {
	return &Handler{
		Repo:   repo,
		Config: cfg,
		Redis:  r,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Публичные маршруты (без авторизации)
	router.GET("/api/electrolytes", h.GetElectrolytes)
	router.GET("/api/electrolytes/:id", h.GetElectrolyte)

	// Регистрация и авторизация (без middleware)
	router.POST("/api/register", h.Register)
	router.POST("/api/login", h.Login)

	// Маршруты для авторизованных пользователей (создатель заявки)
	auth := router.Group("/api")
	auth.Use(h.AuthMiddleware()) // ← middleware применяется ко всей группе
	{
		// Logout внутри группы, чтобы использовать AuthMiddleware
		auth.POST("/logout", h.Logout)
		
		auth.GET("/cart", h.GetCartInfo)
		auth.POST("/cart/items", h.AddItem)
		auth.PUT("/cart/items/:id", h.UpdateItem)
		auth.DELETE("/cart/items/:id", h.DeleteItem)

		auth.GET("/concentrations", h.GetConcentrations)
		auth.GET("/concentrations/:id", h.GetConcentration)
		auth.PUT("/concentrations/:id", h.UpdateConcentration)
		auth.PUT("/concentrations/:id/formed", h.FormConcentration)
		auth.DELETE("/concentrations/:id", h.DeleteConcentration)

		// Маршруты для модератора (требуют дополнительной проверки)
		mod := auth.Group("/")
		mod.Use(RequireModerator())
		{
			mod.PUT("/concentrations/:id/finish", h.FinishConcentration)
			mod.PUT("/concentrations/:id/reject", h.RejectConcentration)
			mod.GET("/concentrations/all", h.GetAllConcentrationsForModerator)
		}
	}
}

func (h *Handler) errorJSON(c *gin.Context, code int, err error) {
	logrus.Error(err)
	c.JSON(code, gin.H{"error": err.Error()})
}