package utils

import (
	"log"
	"os"

	"github.com/fluffy-octo/ik-reddit-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitialiseDB() {
	var err error
	log.Print("Initialising Database...")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		DB, err = gorm.Open(sqlite.Open("ikreddit.db"), &gorm.Config{})
	} else {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		log.Panic(err)
	}

	log.Print("Successfully connected!")

	setupModels(
		&models.Admin{},
		&models.User{},
		&models.Loan{},
		&models.Payments{},
		&models.Agent{},
		models.UserDetails{},
	)
}

func setupModels(models ...interface{}) {
	err := DB.AutoMigrate(models...)
	if err != nil {
		panic(err)
	}
}