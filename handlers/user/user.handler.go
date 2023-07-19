package userHandler

import (
	userRepo "balkantask/database/userRepo"
	"balkantask/model"
	userSchema "balkantask/schemas/user"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	var input userSchema.SignUpInput
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	user := model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		// Set other fields accordingly
	}

	errors := model.ValidateStruct(user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Error",
			"status":  "error",
			"errors":  errors,
		})
	}

	// Check if email already exists
	existingUser, err := userRepo.FindUserByEmail(input.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	if existingUser.ID != uuid.Nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already in use",
			"status":  "error",
		})
	}

	createdUser, err := userRepo.CreateUser(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	response := userSchema.MapUserRecord(&createdUser)

	return c.Status(201).JSON(fiber.Map{
		"message": "Created",
		"status":  "success",
		"data":    response,
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

	updatedUser, err := userRepo.UpdateUser(user)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	response := userSchema.MapUserRecord(&updatedUser)

	return c.Status(200).JSON(fiber.Map{
		"message": "OK",
		"status":  "success",
		"data":    response,
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

	deletedUser, err := userRepo.DeleteUser(user)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	response := userSchema.MapUserRecord(&deletedUser)

	return c.Status(200).JSON(fiber.Map{
		"message": "OK",
		"status":  "success",
		"data":    response,
	})
}
