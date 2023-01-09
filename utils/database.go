package utils

import (
	"fmt"
	"log"
	"os"
	"pendekin/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(model.Url{}, model.User{}, model.Session{}); err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PWD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	dbcon, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "pgx",
		DSN: dsn,
	}), &gorm.Config{})
	
	if err != nil {
		return nil, err
	}

	return dbcon, nil
}

func DisconnectDB(db *gorm.DB) {
	dbConn, err := db.DB()
	if err != nil {
		log.Fatal("Failed to disconnect database!")
	}
	dbConn.Close()
}