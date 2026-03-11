package main

import (
	"lab/internal/app/config"
	"lab/internal/app/dsn"
	"lab/internal/app/handler"
	"lab/internal/app/repository"
	"lab/internal/pkg"
	"log"
	"os"

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
	if dsnStr == "" {
		log.Fatal("DSN empty")
	}

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	minioBucket := os.Getenv("MINIO_BUCKET")
	minioUseSSL := os.Getenv("MINIO_USE_SSL") == "true"

	repo, err := repository.New(&repository.Settings{
		PostgresDSN:    dsnStr,
		MinioEndpoint:  minioEndpoint,
		MinioAccessKey: minioAccessKey,
		MinioSecretKey: minioSecretKey,
		MinioBucket:    minioBucket,
		MinioUseSSL:    minioUseSSL,
	})
	
	if err != nil {
		log.Fatal("repo error:", err)
	}

	hand := handler.NewHandler(repo)
	router := gin.Default()
	app := pkg.NewApp(cfg, router, hand)
	app.Run()
}