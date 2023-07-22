package authHandler

import (
	orgRepo "balkantask/database/org"
	userRepo "balkantask/database/user"
	"balkantask/model"
	orgSchema "balkantask/schemas/org"
	userSchema "balkantask/schemas/user"
	constants "balkantask/utils"
	"balkantask/utils/roles"
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
	user, err := userRepo.FindUserByOrgAndUsernameWithPassword(strings.ToLower(payload.Username), payload.AccountId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": "Invalid username or Password"})
	}

	if user.AccountStatus == constants.DELETED {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": "Invalid username or Password"})
	}

	if user.AccountStatus == constants.DEACTIVATED {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": "Account deactivated. Contact your admin"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": "Invalid username or Password"})
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
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "false", "message": "Internal Server Error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "token": tokenString})
}

func SignInOrg(c *fiber.Ctx) error {
	var payload *orgSchema.SignInInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": err.Error()})
	}

	errors := model.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	var org model.Org
	org, err := orgRepo.FindOrgByEmail(strings.ToLower(payload.Email))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": "Invalid email or Password"})
	}

	if org.AccountStatus == constants.DELETED {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": "Invalid email or Password"})
	}

	if org.AccountStatus == constants.DEACTIVATED {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": "Account deactivated. Contact support."})
	}

	err = bcrypt.CompareHashAndPassword([]byte(org.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": "Invalid email or Password"})
	}

	// Create a new JWT token with a custom expiration time
	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	expirationTime := now.Add(time.Hour * 24) // Token will expire in 24 hours

	claims := tokenByte.Claims.(jwt.MapClaims)
	claims["sub"] = org.ID
	claims["exp"] = expirationTime.Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	config := os.Getenv("JWT_SECRET")
	tokenString, err := tokenByte.SignedString([]byte(config))
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "false", "message": fmt.Sprintf("generating JWT Token failed: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "token": tokenString})
}

func GetMe(c *fiber.Ctx) error {
	if user, ok := c.Locals("user").(userSchema.UserResponse); ok {

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
	}

	if org, ok := c.Locals("org").(orgSchema.OrgResponse); ok {

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"org": org}})
	}

	// Handle the case when the value is nil or not of the correct type
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid token"})
}

func SignUpOrg(c *fiber.Ctx) error {
	var input orgSchema.SignupInput
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	if input.Password != input.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password and password confirmation do not match",
			"status":  "error",
		})
	}

	if !input.ValidatePassword() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number and one special character.",
			"status":  "error",
		})
	}

	// Check if email already exists
	exisitingOrg, err := orgRepo.FindOrgByEmail(input.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	if exisitingOrg.AccountStatus == constants.DELETED {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email cannot be used. Contact support.",
			"status":  "error",
		})
	}

	if exisitingOrg.ID != uuid.Nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already in use",
			"status":  "error",
		})
	}

	org := model.Org{
		Username: input.Username,
		Email:    input.Email,
	}

	errors := model.ValidateStruct(org)
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

	org.Password = string(hashedPassword)

	createdOrg, err := orgRepo.CreateOrg(org)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	response := orgSchema.MapOrgRecord(&createdOrg)

	return c.Status(201).JSON(fiber.Map{
		"message": "Created",
		"status":  "success",
		"data":    response,
	})
}

func DeleteAccount(c *fiber.Ctx) error {
	id_ := c.Params("id")
	id, err := uuid.Parse(id_)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	userLoggedIn, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.UserIsAuthorized(userLoggedIn.Roles, userLoggedIn.Groups, []roles.Role{roles.OrgFullAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	orgToDeactivate, err := orgRepo.FindOrgById(id)

	if orgToDeactivate.AccountStatus == constants.DELETED {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Account is deleted",
			"status":  "error",
		})
	}

	if orgToDeactivate.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Account Not Found",
			"status":  "false",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Account Not Found",
			"status":  "false",
		})
	}

	// Perform the user deactivation logic here
	// ...

	// Update the user with the new account status
	orgToDeactivate.AccountStatus = constants.DEACTIVATED
	updatedOrg, err := orgRepo.UpdateOrg(orgToDeactivate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Account deletion initiated. Your account is under review, and will be marked as deleted in 5 days. Your data will be deleted after 45 days.",
		"status":  "success",
		"data":    updatedOrg,
	})
}

func ChangePassword(c *fiber.Ctx) error {
	var input orgSchema.UpdatePassword
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	if input.Password != input.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password and password confirmation do not match",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)
	if !orgOK && !userOK && !roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.OrgFullAccess, roles.OrgWriteAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	// Validate the input fields
	errors := model.ValidateStruct(input)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Error",
			"status":  "error",
			"errors":  errors,
		})
	}

	org_, err := orgRepo.FindOrgById(input.OrgId)
	if org_.AccountStatus == constants.DEACTIVATED || org_.AccountStatus == constants.DELETED {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Account is deactivated",
			"status":  "error",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Account Not Found",
			"status":  "false",
		})
	}

	// check if new password is the same as the old password
	err = bcrypt.CompareHashAndPassword([]byte(org_.Password), []byte(input.Password))
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "New password cannot be the same as the old password",
			"status":  "error",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	org_.Password = string(hashedPassword)

	updatedUser, err := orgRepo.UpdateOrg(org_)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update password",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password updated successfully",
		"status":  "success",
		"data":    updatedUser,
	})
}
