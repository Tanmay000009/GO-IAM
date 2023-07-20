package userHandler

import (
	userRepo "balkantask/database/user"
	"balkantask/model"
	orgSchema "balkantask/schemas/org"
	userSchema "balkantask/schemas/user"
	userroles "balkantask/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	var input userSchema.CreateUser
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	if input.Password != input.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password and password confirmation do not match",
			"status":  "error",
		})
	}

	// Check if email already exists
	existingUser, err := userRepo.FindUserByUsername(input.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	if existingUser.ID != uuid.Nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Username already in use",
			"status":  "error",
		})
	}

	// Fetch the org details from the middleware
	org, orgExists := c.Locals("org").(orgSchema.OrgResponse)

	if !orgExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unauthorized",
			"status":  "error",
		})
	}

	user := model.User{
		Username: input.Username,
	}

	errors := model.ValidateStruct(user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Error",
			"status":  "error",
			"errors":  errors,
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	user.Password = string(hashedPassword)
	user.OrgID = org.ID // Set the organization ID for the user
	user.Roles = []userroles.Role{}
	createdUser, err := userRepo.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	response := userSchema.MapUserRecord(&createdUser)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
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
