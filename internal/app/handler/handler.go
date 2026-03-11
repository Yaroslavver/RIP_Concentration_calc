package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go_project2/internal/app/repository"
)

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{Repo: repo}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	router.GET("/", h.GetElectrolytes)
	router.GET("/electrolyte/:id", h.GetElectrolyte)
	router.GET("/concentration/:id", h.GetConcentration)
	router.POST("/add-to-concentration", h.AddToConcentration)
	router.POST("/delete-concentration", h.DeleteConcentration)
}

func (h *Handler) errorJSON(c *gin.Context, code int, err error) {
	logrus.Error(err)
	c.JSON(code, gin.H{"error": err.Error()})
}