package authHandler

import (
	userRepo "balkantask/database/user"
	"balkantask/model"
	authSchema "balkantask/schemas/auth"
	userSchema "balkantask/schemas/user"
	"fmt"
	"os"
	"time"

	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignInUser(c *fiber.Ctx) error {
	var payload *userSchema.SignInInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": err.Error()})
	}

	errors := model.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	var user model.User
	_, err := userRepo.FindUserByEmail(strings.ToLower(payload.Email))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": "Invalid email or Password"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": "Invalid email or Password"})
	}

	// Create a new JWT token with a custom expiration time
	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	expirationTime := now.Add(time.Hour * 24) // Token will expire in 24 hours

	claims := tokenByte.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID
	claims["exp"] = expirationTime.Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	config := os.Getenv("JWT_SECRET")
	tokenString, err := tokenByte.SignedString([]byte(config))
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "false", "message": fmt.Sprintf("generating JWT Token failed: %v", err)})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   86400, // Token will expire in 24 hours (in seconds)
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "token": tokenString})
}

func LogoutUser(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expired,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func GetMe(c *fiber.Ctx) error {
	if user, ok := c.Locals("user").(userSchema.UserResponse); ok {
		// The value is not nil and has the correct type (userSchema.UserResponse)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
	}

	// Handle the case when the value is nil or not of the correct type
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid user data"})
}

func SignUpOrg(c *fiber.Ctx) error {
	var input authSchema.SignupInput
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
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

	user := model.User{
		Username: input.Username,
		Email:    input.Email,
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "Internal Server Error", "message": err.Error()})
	}

	user.Password = string(hashedPassword)

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
