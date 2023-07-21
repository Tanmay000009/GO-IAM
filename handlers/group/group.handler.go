package groupHandler

import (
	groupRepo "balkantask/database/group"
	rolesRepo "balkantask/database/roles"
	"balkantask/model"
	groupSchema "balkantask/schemas/group"
	orgSchema "balkantask/schemas/org"
	userSchema "balkantask/schemas/user"
	"balkantask/utils/roles"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllGroups(c *fiber.Ctx) error {
	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.HasAnyRole(user.Roles, []roles.Role{roles.GroupWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.GroupFullAccess, roles.OrgReadAccess, roles.GroupReadAccess}) {
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

	if !orgOK && !userOK && !roles.HasAnyRole(user.Roles, []roles.Role{roles.GroupWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.GroupFullAccess, roles.OrgReadAccess, roles.GroupReadAccess}) {
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

	if !orgOK && !userOK && !roles.HasAnyRole(user.Roles, []roles.Role{roles.GroupWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.GroupFullAccess}) {
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

	// Check if the roles exist in the database
	rolesExist, err := rolesRepo.GetRolesByIds(group.RoleIds)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Role IDs",
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

	// Create an instance of model.Group
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

	if !orgOK && !userOK && !roles.HasAnyRole(user.Roles, []roles.Role{roles.GroupWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.GroupFullAccess}) {
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

	if !orgOK && !userOK && !roles.HasAnyRole(user.Roles, []roles.Role{roles.OrgFullAccess, roles.GroupFullAccess, roles.OrgWriteAccess, roles.GroupWriteAccess}) {
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

	group, err := groupRepo.GetGroupById(input.GroupId)

	if err != nil || group.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Group Not Found",
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

	if roles.HasAnyRole(group.Roles, []roles.Role{roles.Role(role.Name)}) {
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
		"message": "Role added to user",
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

	role, err := rolesRepo.GetRoleById(input.RoleId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role doesn't exist",
			"status":  "error",
		})
	}

	group, err := groupRepo.GetGroupById(input.GroupId)

	if err != nil || group.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Group Not Found",
			"status":  "false",
		})
	}

	// Check if the user has the role
	if !roles.HasAnyRole(group.Roles, []roles.Role{roles.Role(role.Name)}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Group does not have the role",
			"status":  "error",
		})
	}

	// Perform the role deletion logic here (if any)

	// Delete the role from the user
	group, err = groupRepo.RemoveRoleFromGroup(group, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role removed from user",
		"status":  "success",
		"data":    group,
	})
}
