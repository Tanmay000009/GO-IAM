package userHandler

import (
	groupRepo "balkantask/database/group"
	rolesRepo "balkantask/database/roles"
	userRepo "balkantask/database/user"
	"balkantask/model"
	orgSchema "balkantask/schemas/org"
	userSchema "balkantask/schemas/user"
	constants "balkantask/utils"
	"balkantask/utils/roles"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	pass "github.com/sethvargo/go-password/password"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUsers(c *fiber.Ctx) error {
	org, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK {
		return c.Status(400).JSON(fiber.Map{
			"message": "Unauthorized",
			"status":  "error",
		})
	}

	var users []userSchema.UserResponse
	var err error

	if orgOK && org.ID != uuid.Nil {
		users, err = userRepo.FindUsersByOrgId(org.ID)
	} else if userOK {
		if !roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.UserReadAccess, roles.OrgFullAccess, roles.OrgReadAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}) {
			return c.Status(403).JSON(fiber.Map{
				"message": "Forbidden",
				"status":  "error",
			})
		}

		users, err = userRepo.FindUsersByOrgId(user.OrgId)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid token",
			"status":  "error",
		})
	}

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
		if user.ID != id_uuid && !roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.UserReadAccess, roles.OrgFullAccess, roles.OrgReadAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}) {
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
		return c.Status(404).JSON(fiber.Map{
			"message": "User Not Found",
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

	org, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	var orgId uuid.UUID

	if org.ID != uuid.Nil {
		orgId = org.ID
	} else {
		orgId = user.OrgId
	}

	exisitingUser, err := userRepo.FindUserByOrgAndUsernameWithPassword(input.Username, orgId.String())
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
				"status":  "error",
			})
		}
	}

	if exisitingUser.ID != uuid.Nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Username already in use",
			"status":  "error",
		})
	}

	if input.Password != input.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Passwords do not match",
			"status":  "error",
		})
	}

	// If password is not provided, generate a random password
	if input.Password == "" {
		input.Password, err = pass.Generate(10, 4, 2, true, true)
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

	if org.ID != uuid.Nil {
		newUser.OrgID = org.ID
	} else {
		newUser.OrgID = user.OrgId
	}

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

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && user.ID == updatedUser.ID && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

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

	if !(orgOK || (userOK && roles.UserIsAuthorized(userLoggedIn.Roles, userLoggedIn.Groups, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	userToDelete, err := userRepo.FindUserByIdWithPassword(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Not Found",
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

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}))) {
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

	// Check if the input contains Role ID or Role Name
	var role model.Role
	if input.RoleId != uuid.Nil {
		role, err = rolesRepo.GetRoleById(input.RoleId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Role doesn't exist",
				"status":  "error",
			})
		}
	} else if input.RoleName != "" {
		role, err = rolesRepo.GetRoleByName(input.RoleName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Role doesn't exist",
				"status":  "error",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role ID or Role Name is required",
			"status":  "error",
		})
	}

	// Check if the user already has the role
	if roles.UserHasRole(user_.Roles, role) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already has the role",
			"status":  "error",
		})
	}

	user_, err = userRepo.AddRoleToUser(role, user_)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	mappedUser := userSchema.MapUserRecord(&user_)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role added to user",
		"status":  "success",
		"data":    mappedUser,
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

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}))) {
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

	// Check if the input contains Role ID or Role Name
	var role model.Role
	if input.RoleId != uuid.Nil {
		role, err = rolesRepo.GetRoleById(input.RoleId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Role doesn't exist",
				"status":  "error",
			})
		}
	} else if input.RoleName != "" {
		role, err = rolesRepo.GetRoleByName(input.RoleName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Role doesn't exist",
				"status":  "error",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role ID or Role Name is required",
			"status":  "error",
		})
	}

	// Check if the user has the role
	if !roles.UserHasRole(user_.Roles, role) {
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

	mappedUser := userSchema.MapUserRecord(&updatedUser)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role removed from user",
		"status":  "success",
		"data":    mappedUser,
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

	if !(orgOK || (userOK && roles.UserIsAuthorized(userLoggedIn.Roles, userLoggedIn.Groups, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}))) {
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

func ReactivateUser(c *fiber.Ctx) error {
	id_ := c.Params("id")
	id, err := uuid.Parse(id_)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	userLoggedIn, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(userLoggedIn.Roles, userLoggedIn.Groups, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	userToReactivate, err := userRepo.FindUserByIdWithPassword(id)

	if userToReactivate.AccountStatus != constants.DEACTIVATED {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "User Account is already activated",
			"status":  "error",
		})
	}

	if userToReactivate.ID == uuid.Nil {
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

	userToReactivate.AccountStatus = constants.ACTIVATED
	updatedUser, err := userRepo.UpdateUser(userToReactivate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User reactivated successfully.",
		"status":  "success",
		"data":    updatedUser,
	})
}

func SeedUsersFromExcel(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid file",
			"status":  "error",
		})
	}

	uploadedFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to read uploaded file",
			"status":  "error",
		})
	}
	// Close the file after the function returns
	defer uploadedFile.Close()

	org, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	var orgId uuid.UUID

	if org.ID != uuid.Nil {
		orgId = org.ID
	} else {
		orgId = user.OrgId
	}

	// Create a temporary file to save the uploaded content
	tempFile, err := os.CreateTemp("", "upload-*.xlsx")
	// CreateTemp function, it generates a unique temporary file name by replacing the asterisk (*) with a random string.
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create temporary file",
			"status":  "error",
		})
	}
	defer os.Remove(tempFile.Name())

	// Save the uploaded content into the temporary file
	_, err = io.Copy(tempFile, uploadedFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save uploaded file",
			"status":  "error",
		})
	}

	// Open the temporary file using excelize
	xlsx, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to read Excel file",
			"status":  "error",
		})
	}

	// Define the columns to read from the Excel file (adjust the column numbers accordingly)
	usernameCol := 1
	passwordCol := 2

	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to read Excel file",
			"status":  "error",
		})
	}

	var seededUsers []userSchema.CreateUserResponse

	for rowIndex, row := range rows {

		if rowIndex == 0 {
			continue
		}

		// Check if the row has enough columns, if not, set an empty password
		if len(row) < passwordCol {
			row = append(row, "")
		}

		username := row[usernameCol-1]
		password := row[passwordCol-1]

		excelUser := model.User{
			Username: username,
			Password: password,
		}
		errors := model.ValidateStruct(excelUser)
		if errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Validation error in row %d", rowIndex+1),
				"status":  "error",
				"errors":  errors,
			})
		}

		existingUser, err := userRepo.FindUserByOrgAndUsernameWithPassword(username, orgId.String())
		if err != nil && err != gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"status":  "error",
			})
		}

		if existingUser.ID != uuid.Nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("User with username %s already exists", username),
				"status":  "error",
			})
		}

		if excelUser.Password == "" {
			excelUser.Password, err = pass.Generate(10, 4, 2, true, true)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Internal Server Error",
					"status":  "error",
				})
			}
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(excelUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"status":  "error",
			})
		}

		newUser := model.User{
			Username:      excelUser.Username,
			Password:      string(hashedPassword),
			AccountStatus: constants.ACTIVATED,
			OrgID:         orgId,
		}

		createdUser, err := userRepo.CreateUser(newUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("Failed to seed user in row %d", rowIndex+1),
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
			Passcode:      excelUser.Password,
		}
		seededUsers = append(seededUsers, resData)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Users seeded successfully",
		"status":  "success",
		"data":    seededUsers,
	})
}

func SeedUsersFromCSV(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid file",
			"status":  "error",
		})
	}

	uploadedFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to read uploaded file",
			"status":  "error",
		})
	}
	// Close the file after the function returns
	defer uploadedFile.Close()

	org, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	var orgId uuid.UUID

	if org.ID != uuid.Nil {
		orgId = org.ID
	} else {
		orgId = user.OrgId
	}

	// Create a temporary file to save the uploaded content
	tempFile, err := os.CreateTemp("", "upload-*.csv")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create temporary file",
			"status":  "error",
		})
	}
	defer os.Remove(tempFile.Name())

	// Save the uploaded content into the temporary file
	_, err = io.Copy(tempFile, uploadedFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save uploaded file",
			"status":  "error",
		})
	}

	// Open the temporary file using os
	csvFile, err := os.Open(tempFile.Name())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to read CSV file",
			"status":  "error",
		})
	}
	defer csvFile.Close()

	// Create a new CSV reader
	reader := csv.NewReader(csvFile)

	// Skip the header row
	if _, err := reader.Read(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to read CSV file",
			"status":  "error",
		})
	}

	var seededUsers []userSchema.CreateUserResponse

	for rowIndex := 1; ; rowIndex++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Failed to read CSV file",
				"status":  "error",
			})
		}

		// Check if the row has enough columns, if not, set an empty password
		if len(row) < 2 {
			row = append(row, "")
		}

		username := row[0]
		password := row[1]

		csvUser := model.User{
			Username: username,
			Password: password,
		}
		errors := model.ValidateStruct(csvUser)
		if errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Validation error in row %d", rowIndex),
				"status":  "error",
				"errors":  errors,
			})
		}

		existingUser, err := userRepo.FindUserByOrgAndUsernameWithPassword(username, orgId.String())
		if err != nil && err != gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"status":  "error",
			})
		}

		if existingUser.ID != uuid.Nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("User with username %s already exists", username),
				"status":  "error",
			})
		}

		if csvUser.Password == "" {
			csvUser.Password, err = pass.Generate(10, 4, 2, true, true)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Internal Server Error",
					"status":  "error",
				})
			}
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(csvUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"status":  "error",
			})
		}

		newUser := model.User{
			Username:      csvUser.Username,
			Password:      string(hashedPassword),
			AccountStatus: constants.ACTIVATED,
			OrgID:         orgId,
		}

		createdUser, err := userRepo.CreateUser(newUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("Failed to seed user in row %d", rowIndex),
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
			Passcode:      csvUser.Password,
		}
		seededUsers = append(seededUsers, resData)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Users seeded successfully",
		"status":  "success",
		"data":    seededUsers,
	})
}

func ChangePassword(c *fiber.Ctx) error {
	var input userSchema.UpdatePassword
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

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)
	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}))) {
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

	// check if new password is the same as the old password
	err = bcrypt.CompareHashAndPassword([]byte(user_.Password), []byte(input.Password))
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

	user_.Password = string(hashedPassword)

	updatedUser, err := userRepo.UpdateUser(user_)
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

func AddGroupToUser(c *fiber.Ctx) error {
	var input userSchema.AddOrDeleteGroup
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}))) {
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

	// Check if the input contains Group ID or Group Name
	var group model.Group
	if input.GroupId != uuid.Nil {
		group, err = groupRepo.GetGroupById(input.GroupId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Group doesn't exist",
				"status":  "error",
			})
		}
	} else if input.GroupName != "" {
		group, err = groupRepo.GetGroupByName(input.GroupName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Group doesn't exist",
				"status":  "error",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Group ID or Group Name is required",
			"status":  "error",
		})
	}

	// Check if the user already has the group
	if roles.UserHasGroup(user_.Groups, []model.Group{group}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already has the group",
			"status":  "error",
		})
	}

	user_, err = userRepo.AddGroupToUser(group, user_)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	mappedUser := userSchema.MapUserRecord(&user_)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Group added to user",
		"status":  "success",
		"data":    mappedUser,
	})
}

func DeleteGroupFromUser(c *fiber.Ctx) error {
	var input userSchema.AddOrDeleteGroup
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.OrgFullAccess, roles.UserFullAccess, roles.OrgWriteAccess, roles.UserWriteAccess}))) {
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

	// Check if the input contains Group ID or Group Name
	var group model.Group
	if input.GroupId != uuid.Nil {
		group, err = groupRepo.GetGroupById(input.GroupId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Group doesn't exist",
				"status":  "error",
			})
		}
	} else if input.GroupName != "" {
		group, err = groupRepo.GetGroupByName(input.GroupName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Group doesn't exist",
				"status":  "error",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Group ID or Group Name is required",
			"status":  "error",
		})
	}

	// Check if the user has the group
	if !roles.UserHasGroup(user_.Groups, []model.Group{group}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User does not have the group",
			"status":  "error",
		})
	}

	updatedUser, err := userRepo.DeleteGroupFromUser(group, user_)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Group removed from user",
		"status":  "success",
		"data":    updatedUser,
	})
}
