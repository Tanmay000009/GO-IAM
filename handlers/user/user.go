package userHandler

import (
	userRepo "balkantask/database/userRepo"
	"balkantask/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUsers(c *fiber.Ctx) error {

	users, err := userRepo.FindUsers()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "OK",
		"status":  "success",
		"data":    users,
	})
}

func GetUserById(c *fiber.Ctx) error {

	id := c.Params("id")

	// validate if id is valid uuid
	_, err := uuid.Parse(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID",
			"status":  "error",
		})
	}

	user, err := userRepo.FindUserById(id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "OK",
		"status":  "success",
		"data":    user,
	})
}

func CreateUser(c *fiber.Ctx) error {

	var user model.User

	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	user, err = userRepo.CreateUser(user)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Created",
		"status":  "success",
		"data":    user,
	})
}

func UpdateUser(c *fiber.Ctx) error {

	var user model.User

	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	user, err = userRepo.UpdateUser(user)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "OK",
		"status":  "success",
		"data":    user,
	})
}

func DeleteUser(c *fiber.Ctx) error {

	var user model.User

	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	user, err = userRepo.DeleteUser(user)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "OK",
		"status":  "success",
		"data":    user,
	})
}
