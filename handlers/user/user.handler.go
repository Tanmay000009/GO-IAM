package userHandler

import (
	rolesRepo "balkantask/database/roles"
	userRepo "balkantask/database/user"
	"balkantask/model"
	orgSchema "balkantask/schemas/org"
	userSchema "balkantask/schemas/user"
	constants "balkantask/utils"
	"balkantask/utils/roles"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

	exisitingUser, err := userRepo.FindUserByUsername(input.Username)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
				"status":  "error",
			})
		}
	}

	// Check the existence of the user.
	if exisitingUser.ID != uuid.Nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Username already in use",
			"status":  "error",
		})
	}

	if input.Password != input.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Passwords do not match",
			"status":  "error",
		})
	}

	// If password is not provided, generate a random password
	if input.Password == "" {
		input.Password, err = password.Generate(10, 4, 2, true, true)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"status":  "error",
			})
		}
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

	resData := userSchema.CreateUserResponse{
		ID:            createdUser.ID,
		Username:      createdUser.Username,
		CreatedAt:     createdUser.CreatedAt,
		UpdatedAt:     createdUser.UpdatedAt,
		Roles:         createdUser.Roles,
		OrgId:         createdUser.OrgId,
		AccountStatus: createdUser.AccountStatus,
		Passcode:      input.Password,
	}

	// Return the created user data
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Created",
		"status":  "success",
		"data":    resData,
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
	existingUser, err := userRepo.FindUserById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	if existingUser.AccountStatus == constants.DEACTIVATED || existingUser.AccountStatus == constants.DELETED {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Account is deactivated",
			"status":  "error",
		})
	}

	if existingUser.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Not Found",
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

	if !orgOK && !userOK && userLoggedIn.ID == id && !roles.HasAnyRole(userLoggedIn.Roles, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	userToDelete, err := userRepo.FindUserByIdWithPassword(id)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "false",
		})
	}

	if userToDelete.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Not Found",
			"status":  "false",
		})
	}

	userToDelete, err = userRepo.FindUserByIdWithPassword(id)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Not Found",
			"status":  "false",
		})
	}

	// Delete the user
	userDeleted, err := userRepo.DeleteUser(userToDelete)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
		"status":  "success",
		"data":    userDeleted,
	})
}

func AddRoleToUser(c *fiber.Ctx) error {
	var input userSchema.AddOrDeleteRole
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

	if user_.AccountStatus == constants.DEACTIVATED || user_.AccountStatus == constants.DELETED {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Account is deactivated",
			"status":  "error",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Not Found",
			"status":  "false",
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

	user_, err = userRepo.AddRoleToUser(role, user_)

	// updatedUser, err := userRepo.UpdateUser(user_)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role added to user",
		"status":  "success",
		"data":    user_,
	})
}

func DeleteRoleFromUser(c *fiber.Ctx) error {
	var input userSchema.AddOrDeleteRole
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
	if user_.AccountStatus == constants.DEACTIVATED || user_.AccountStatus == constants.DELETED {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Account is deactivated",
			"status":  "error",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Not Found",
			"status":  "false",
		})
	}

	role, err := rolesRepo.GetRoleById(input.RoleId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role doesn't exist",
			"status":  "error",
		})
	}

	// Check if the user has the role
	if !roles.HasAnyRole(user_.Roles, []roles.Role{roles.Role(role.Name)}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User does not have the role",
			"status":  "error",
		})
	}

	updatedUser, err := userRepo.DeleteRoleFromUser(role, user_)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role removed from user",
		"status":  "success",
		"data":    updatedUser,
	})
}

func DeactivateUser(c *fiber.Ctx) error {
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

	if !orgOK && !userOK && !roles.HasAnyRole(userLoggedIn.Roles, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	userToDeactivate, err := userRepo.FindUserByIdWithPassword(id)

	if userToDeactivate.AccountStatus == constants.DELETED {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "User is deleted",
			"status":  "error",
		})
	}

	if userToDeactivate.AccountStatus == constants.DEACTIVATED {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "User Account is already deactivated",
			"status":  "error",
		})
	}

	if userToDeactivate.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Not Found",
			"status":  "false",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Not Found",
			"status":  "false",
		})
	}

	// Perform the user deactivation logic here
	// ...

	// Update the user with the new account status
	userToDeactivate.AccountStatus = constants.DEACTIVATED
	updatedUser, err := userRepo.UpdateUser(userToDeactivate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deactivated successfully. User can be reactivated. User data will be deleted after 30 days.",
		"status":  "success",
		"data":    updatedUser,
	})
}
