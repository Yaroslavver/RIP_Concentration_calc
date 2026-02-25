package pkg

import (
	"fmt"
	"go_project2/internal/app/config"
	"go_project2/internal/app/handler"
	//"go_project2/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type App struct {
	Config  *config.Config
	Router  *gin.Engine
	Handler *handler.Handler
}

func NewApp(cfg *config.Config, router *gin.Engine, h *handler.Handler) *App {
	return &App{Config: cfg, Router: router, Handler: h}
}

func (a *App) Run() {
	a.Handler.RegisterRoutes(a.Router)
	addr := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
	logrus.Info("Server started at ", addr)
	if err := a.Router.Run(addr); err != nil {
		logrus.Fatal(err)
	}
}