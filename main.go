package main

import (
	"balkantask/database"
	"balkantask/router"
	"balkantask/utils/schedulers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
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

	go schedulers.Scheduler()

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
