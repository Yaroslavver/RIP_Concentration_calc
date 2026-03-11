package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab/internal/app/repository"
)

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{Repo: repo}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*") // для возможных HTML-страниц
	router.Static("/static", "./static")

	api := router.Group("/api")
	{
		// Электролиты
		api.GET("/electrolytes", h.GetElectrolytes)
		api.GET("/electrolytes/:id", h.GetElectrolyte)
		api.POST("/electrolytes", h.CreateElectrolyte)

		// Корзина
		api.GET("/cart", h.GetCartInfo)

		// Заявки
		api.GET("/concentrations", h.GetConcentrations)
		api.GET("/concentrations/:id", h.GetConcentration)
		api.PUT("/concentrations/:id", h.UpdateConcentration)
		api.PUT("/concentrations/:id/formed", h.FormConcentration)
		api.PUT("/concentrations/:id/finish", h.FinishConcentration)
		api.PUT("/concentrations/:id/reject", h.RejectConcentration)
		api.DELETE("/concentrations/:id", h.DeleteConcentration)

		// Элементы
		api.POST("/cart/items", h.AddItem)
		api.PUT("/cart/items/:id", h.UpdateItem)
		api.DELETE("/cart/items/:id", h.DeleteItem)

		// Пользователи
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
		api.POST("/logout", h.Logout)
	}
}

func (h *Handler) errorJSON(c *gin.Context, code int, err error) {
	logrus.Error(err)
	c.JSON(code, gin.H{"error": err.Error()})
}