package main

import (
	"lab/internal/app/config"
	"lab/internal/app/dsn"
	"lab/internal/app/handler"
	"lab/internal/app/redis"
	"lab/internal/app/repository"
	"lab/internal/pkg"
	"log"
	"os"
	//"context"
	_ "lab/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// @title           Electrolyte Concentration API
// @version         1.0
// @description     API for calculating ion concentration in electrolyte mixtures.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

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

	// Подключение к Redis
	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		log.Fatal("redis error:", err)
	}
	defer redisClient.Close()

	// Создаем хендлер с конфигурацией и Redis
	hand := handler.NewHandler(repo, cfg, redisClient)

	// Настраиваем роутер
	router := gin.Default()
	hand.RegisterRoutes(router)

	// Запускаем приложение
	app := pkg.NewApp(cfg, router, hand)
	logrus.Info("Server starting...")
	app.Run()
}