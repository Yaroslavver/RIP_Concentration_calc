package main

import (
	"go_project2/internal/app/ds"
	"go_project2/internal/app/dsn"
	"log"

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
	log.Println("Migration completed")
}