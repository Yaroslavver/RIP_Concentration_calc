package main

import (
	"lab/internal/app/ds"
	"lab/internal/app/dsn"
	"log"
	//"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	err = db.AutoMigrate(
		&ds.User{},
		&ds.Electrolyte{},
		&ds.Concentration{},
		&ds.ConcentrationItem{},
	)
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}

	// Создаём тестового пользователя, если его нет
	var count int64
	db.Model(&ds.User{}).Where("login = ?", "test").Count(&count)
	if count == 0 {
		testUser := ds.User{Login: "test", Password: "test123", IsModerator: true} // сделаем модератором для тестов
		if err := db.Create(&testUser).Error; err != nil {
			log.Fatal("failed to create test user:", err)
		}
		log.Println("Test user created")
	}


	log.Println("Migration completed")
}