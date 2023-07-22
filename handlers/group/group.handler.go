package groupHandler

import (
	groupRepo "balkantask/database/group"
	rolesRepo "balkantask/database/roles"
	"balkantask/model"
	groupSchema "balkantask/schemas/group"
	orgSchema "balkantask/schemas/org"
	userSchema "balkantask/schemas/user"
	"balkantask/utils/roles"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

func GetAllGroups(c *fiber.Ctx) error {
	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.GroupWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.GroupFullAccess, roles.OrgReadAccess, roles.GroupReadAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}

	groups, err := groupRepo.GetAllGroups()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "data": groups})
}

func GetGroupById(c *fiber.Ctx) error {
	id_ := c.Params("id")
	id, err := uuid.Parse(id_)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.GroupWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.GroupFullAccess, roles.OrgReadAccess, roles.GroupReadAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}

	group, err := groupRepo.GetGroupById(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "data": group})
}

func CreateGroup(c *fiber.Ctx) error {
	var group groupSchema.CreateGroup

	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.GroupWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.GroupFullAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}

	errors := model.ValidateStruct(group)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Error",
			"status":  "error",
			"errors":  errors,
		})
	}

	// Check if the roles (id) exist in the database
	rolesExist, err := rolesRepo.GetRolesByIds(group.RoleIds)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Role IDs",
			"status":  "error",
		})
	}

	// Check if the roles (name) exist in the database
	rolesExist2, err := rolesRepo.GetRolesByNames(group.RoleNames)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Role IDs",
			"status":  "error",
		})
	}

	rolesExist = append(rolesExist, rolesExist2...)
	rolesExist = roles.RemoveDuplicates(rolesExist)

	groupExists, err := groupRepo.GetGroupByName(group.Name)
	if err != nil && err.Error() != "record not found" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	if groupExists.ID != uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Group already exists",
			"status":  "error",
		})
	}

	// Verify if all roles were found
	if len(rolesExist) != len(group.RoleIds) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Role IDs",
			"status":  "error",
		})
	}

	newGroup := model.Group{
		Name:  group.Name,
		Roles: rolesExist,
	}

	createdGroup, err := groupRepo.CreateGroup(&newGroup)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   createdGroup,
	})
}

func DeleteGroupById(c *fiber.Ctx) error {

	id_ := c.Params("id")
	id, err := uuid.Parse(id_)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.GroupWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.GroupFullAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}

	groupExists, err := groupRepo.GetGroupById(id)

	if err != nil || groupExists.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Group Not Found",
			"status":  "false",
		})
	}

	err = groupRepo.DeleteGroup(&groupExists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   true,
	})
}

func AddRoleToGroup(c *fiber.Ctx) error {
	var input groupSchema.AddOrDeleteRole
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.OrgFullAccess, roles.GroupFullAccess, roles.OrgWriteAccess, roles.GroupWriteAccess}))) {
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

	// Check if the input contains Group ID or Group Name
	var group model.Group
	if input.GroupId != uuid.Nil {
		group, err = groupRepo.GetGroupById(input.GroupId)
	} else if input.GroupName != "" {
		group, err = groupRepo.GetGroupByName(input.GroupName)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Group ID or Group Name is required",
			"status":  "error",
		})
	}

	if err != nil || group.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Group Not Found",
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

	if roles.GroupHasRole(group.Roles, []model.Role{role}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Group already has the role",
			"status":  "error",
		})
	}

	group, err = groupRepo.AddRoleToGroup(group, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role added to group",
		"status":  "success",
		"data":    group,
	})
}

func DeleteRoleFromGroup(c *fiber.Ctx) error {
	var input groupSchema.AddOrDeleteRole
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

	// Check if the input contains Group ID or Group Name
	var group model.Group
	if input.GroupId != uuid.Nil {
		group, err = groupRepo.GetGroupById(input.GroupId)
	} else if input.GroupName != "" {
		group, err = groupRepo.GetGroupByName(input.GroupName)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Group ID or Group Name is required",
			"status":  "error",
		})
	}

	if err != nil || group.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Group Not Found",
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

	if !roles.GroupHasRole(group.Roles, []model.Role{role}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Group does not have the role",
			"status":  "error",
		})
	}

	// Delete the role from the user
	group, err = groupRepo.RemoveRoleFromGroup(group, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role removed from group",
		"status":  "success",
		"data":    group,
	})
}

func SeedGroupsFromExcel(c *fiber.Ctx) error {
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

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.GroupWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.GroupFullAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
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
	groupNameCol := 1
	roleNamesCol := 2

	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to read Excel file",
			"status":  "error",
		})
	}

	var createdGroups []model.Group

	for rowIndex, row := range rows {
		if rowIndex == 0 {
			continue
		}

		// Check if the row has enough columns
		if len(row) < roleNamesCol {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Insufficient columns in row %d", rowIndex+1),
				"status":  "error",
			})
		}

		groupName := row[groupNameCol-1]
		roleNames := strings.Split(row[roleNamesCol-1], ",")

		// Trim spaces from role names
		for i := range roleNames {
			roleNames[i] = strings.TrimSpace(roleNames[i])
		}

		// Retrieve the roles from the database based on role names
		rolesExist, err := rolesRepo.GetRolesByNames(roleNames)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Invalid role names in row %d", rowIndex+1),
				"status":  "error",
			})
		}

		newGroup := model.Group{
			Name:  groupName,
			Roles: rolesExist,
		}

		createdGroup, err := groupRepo.CreateGroup(&newGroup)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("Failed to create group in row %d", rowIndex+1),
			})
		}

		createdGroups = append(createdGroups, *createdGroup)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Groups created successfully",
		"status":  "success",
		"data":    createdGroups,
	})
}

func SeedGroupsFromCSV(c *fiber.Ctx) error {
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

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.GroupWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.GroupFullAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
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
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to read CSV file",
			"status":  "error",
		})
	}

	var createdGroups []model.Group

	for rowIndex := 1; ; rowIndex++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Failed to read CSV file",
				"status":  "error",
			})
		}

		// Check if the row has enough columns
		if len(row) < 2 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Insufficient columns in row %d", rowIndex),
				"status":  "error",
			})
		}

		groupName := row[0]
		roleNames := strings.Split(row[1], " ")
		// Trim spaces from role names
		for i := range roleNames {
			roleNames[i] = strings.TrimSpace(roleNames[i])
		}

		// Retrieve the roles from the database based on role names
		rolesExist, err := rolesRepo.GetRolesByNames(roleNames)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Invalid role names in row %d", rowIndex),
				"status":  "error",
			})
		}

		newGroup := model.Group{
			Name:  groupName,
			Roles: rolesExist,
		}

		createdGroup, err := groupRepo.CreateGroup(&newGroup)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("Failed to create group in row %d", rowIndex),
			})
		}

		createdGroups = append(createdGroups, *createdGroup)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Groups created successfully",
		"status":  "success",
		"data":    createdGroups,
	})
}

func TestUserGroup(c *fiber.Ctx) error {
	var group groupSchema.TestGroup

	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	user, userOK := c.Locals("user").(userSchema.UserResponse)

	var groupExists model.Group
	var err error

	if group.GroupId != uuid.Nil {
		groupExists, err = groupRepo.GetGroupById(group.GroupId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid Group ID",
				"status":  "error",
			})
		}
	} else if group.GroupName != "" {
		groupExists, err = groupRepo.GetGroupByName(group.GroupName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid Group Name",
				"status":  "error",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Group ID or Group Name is required",
			"status":  "error",
		})
	}

	if !userOK && !roles.UserHasGroup(user.Groups, []model.Group{groupExists}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "User has role",
		"data":   true,
	})
}
