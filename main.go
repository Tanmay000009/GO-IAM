package main

import (
	"balkantask/database"
	"balkantask/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()

	dsn := "host=localhost user=postgres password=postgres dbname=balkantask port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
			"status":  "error",
		})
	})

	app.Get("api/healtcheck", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "OK",
			"status":  "success",
		})
	})

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
