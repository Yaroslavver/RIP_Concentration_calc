package main

import (
	"go_project2/internal/app/config"
	"go_project2/internal/app/dsn"
	"go_project2/internal/app/handler"
	"go_project2/internal/app/repository"
	"go_project2/internal/pkg"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	//"github.com/sirupsen/logrus"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("config error:", err)
	}

	dsnStr := dsn.FromEnv()
	repo, err := repository.New(dsnStr)
	if err != nil {
		log.Fatal("repo error:", err)
	}

	hand := handler.NewHandler(repo)
	router := gin.Default()
	app := pkg.NewApp(cfg, router, hand)
	app.Run()
}