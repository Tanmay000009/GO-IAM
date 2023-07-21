package database

import (
	"balkantask/model"
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	port_ := os.Getenv("DB_PORT")
	// Parse port to int
	port, err := strconv.ParseUint(port_, 10, 32)
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to the database.\n", err)
		os.Exit(2)
	}

	log.Println("Running database migrations")
	err = db.AutoMigrate(&model.User{}, &model.Org{}, &model.Role{}, &model.Group{})
	if err != nil {
		log.Fatal("Migration failed.\n", err)
		os.Exit(1)
	}

	DB = db
	log.Println("Connected successfully to the database")
}
