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
	router.GET("/calculation/:id", h.GetCalculation)
	router.POST("/add-to-calculation", h.AddToCalculation)
	router.POST("/delete-calculation", h.DeleteCalculation)
}

func (h *Handler) errorJSON(c *gin.Context, code int, err error) {
	logrus.Error(err)
	c.JSON(code, gin.H{"error": err.Error()})
}