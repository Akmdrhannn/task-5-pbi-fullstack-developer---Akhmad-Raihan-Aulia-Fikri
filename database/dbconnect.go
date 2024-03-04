package database

import (
	"btpn-golang/models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	globalInstance *gorm.DB
)

func Connect() (*gorm.DB, error) {

	if globalInstance != nil {
		return globalInstance, nil
	}

	err_load := godotenv.Load()

	if err_load != nil {
		println("Tidak bisa membuka file env")
		return nil, err_load
	}

	connection := os.Getenv("DB_CONNECTION")

	db, err := gorm.Open(mysql.Open(connection), &gorm.Config{})
	if err != nil {
		println("Tidak bisa terhubung ke database")
		return nil, err
	}

	db.AutoMigrate(&models.User{}, &models.Photo{})

	globalInstance = db
	return db, nil
}
