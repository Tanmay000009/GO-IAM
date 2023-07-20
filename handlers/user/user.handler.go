package userHandler

import (
	rolesRepo "balkantask/database/roles"
	userRepo "balkantask/database/user"
	"balkantask/model"
	orgSchema "balkantask/schemas/org"
	userSchema "balkantask/schemas/user"
	roles "balkantask/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *fiber.Ctx) error {
	// Check for org and user in locals
	org, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	// Check for unauthorized access
	if !orgOK && !userOK {
		return c.Status(400).JSON(fiber.Map{
			"message": "Unauthorized",
			"status":  "error",
		})
	}

	// Get users based on the context (organization or user)
	var users []userSchema.UserResponse
	var err error

	if orgOK && org.ID != uuid.Nil {
		// Fetch users based on the organization ID
		users, err = userRepo.FindUsersByOrgId(org.ID)
	} else if userOK {
		// Check if the user has the required role
		if !roles.HasAnyRole(user.Roles, []roles.Role{roles.UserReadAccess, roles.OrgFullAccess, roles.OrgReadAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}) {
			return c.Status(403).JSON(fiber.Map{
				"message": "Forbidden",
				"status":  "error",
			})
		}
		// Fetch users based on the user's OrgID
		users, err = userRepo.FindUsersByOrgId(user.OrgId)
	} else {
		// Handle the case when neither org nor user is present or of the correct type
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid token",
			"status":  "error",
		})
	}

	if err != nil {
		// Handle internal server errors
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	// Return the users data
	return c.Status(200).JSON(fiber.Map{
		"message": "OK",
		"status":  "success",
		"data":    users,
	})
}

func GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if userOK {
		if user.ID != id_uuid && !roles.HasAnyRole(user.Roles, []roles.Role{roles.UserReadAccess, roles.OrgFullAccess, roles.OrgReadAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}) {
			return c.Status(403).JSON(fiber.Map{
				"message": "Forbidden",
				"status":  "error",
			})
		}
	} else if !orgOK {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid token",
			"status":  "error",
		})
	}

	user_, err := userRepo.FindUserWithOrgById(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "OK",
		"status":  "success",
		"data":    user_,
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

	// Check if the user is an organization or has the necessary permission
	org, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.HasAnyRole(user.Roles, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}) {
		// If neither org nor user is present or not of the correct type, or the user doesn't have the necessary permission, return an error
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	// Create the new user
	newUser := model.User{
		Username: input.Username,
		// Set other fields accordingly
	}

	// Validate the new user data
	errors := model.ValidateStruct(newUser)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Error",
			"status":  "error",
			"errors":  errors,
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
	}

	newUser.Password = string(hashedPassword)
	newUser.OrgID = org.ID // Set the organization ID if the user is created by an organization

	createdUser, err := userRepo.CreateUser(newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	// Return the created user data
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Created",
		"status":  "success",
		"data":    createdUser,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	id_ := c.Params("id")
	id, err := uuid.Parse(id_)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID",
			"status":  "error",
		})
	}
	var updatedUser model.User

	err = c.BodyParser(&fiber.Map{
		"username": updatedUser.Username,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	// Check if the user is an organization or the user themselves, or has the necessary permission
	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && user.ID != updatedUser.ID && !roles.HasAnyRole(user.Roles, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	// Fetch the existing user from the database
	_, err = userRepo.FindUserById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	// Validate the updated user data
	errors := model.ValidateStruct(updatedUser)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Error",
			"status":  "error",
			"errors":  errors,
		})
	}

	updatedUser_, err := userRepo.UpdateUser(updatedUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OK",
		"status":  "success",
		"data":    updatedUser_,
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

	userDeleted, err := userRepo.DeleteUser(user)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "OK",
		"status":  "success",
		"data":    userDeleted,
	})
}

func AddRoleToUser(c *fiber.Ctx) error {
	var input userSchema.AddRoleToUser
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.HasAnyRole(user.Roles, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	errors := model.ValidateStruct(input)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Error",
			"status":  "error",
			"errors":  errors,
		})
	}

	user_, err := userRepo.FindUserByIdWithPassword(input.UserId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	role, err := rolesRepo.GetRoleById(input.RoleId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role doesn't exist",
			"status":  "error",
		})
	}

	// Check if the user already has the role
	if roles.HasAnyRole(user_.Roles, []roles.Role{roles.Role(role.Name)}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already has the role",
			"status":  "error",
		})
	}

	user_.Roles = append(user_.Roles, role)

	updatedUser, err := userRepo.UpdateUser(user_)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role added to user",
		"status":  "success",
		"data":    updatedUser,
	})
}
