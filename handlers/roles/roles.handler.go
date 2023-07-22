package rolesHandler

import (
	rolesRepo "balkantask/database/roles"
	"balkantask/model"
	orgSchema "balkantask/schemas/org"
	roleSchema "balkantask/schemas/role"
	userSchema "balkantask/schemas/user"
	"balkantask/utils/roles"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllRoles(c *fiber.Ctx) error {
	roles, err := rolesRepo.GetAllRoles()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "data": roles})
}

func GetRoleById(c *fiber.Ctx) error {

	id_ := c.Params("id")
	id, err := uuid.Parse(id_)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID",
			"status":  "error",
		})
	}

	role, err := rolesRepo.GetRoleById(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "data": role})
}

func CreateRole(c *fiber.Ctx) error {
	var role roleSchema.AddOrDeleteRole

	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.RoleWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.RoleFullAccess, roles.OrgReadAccess, roles.RoleReadAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}

	errors := model.ValidateStruct(role)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Error",
			"status":  "error",
			"errors":  errors,
		})
	}

	exisitingRole, err := rolesRepo.GetRoleByName(role.RoleName)
	if err == nil && exisitingRole.Name == role.RoleName {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role already exists",
			"status":  "error",
		})
	}

	newRole := model.Role{
		Name: role.RoleName,
		Type: role.Type,
	}

	createdRole, err := rolesRepo.CreateRole(newRole)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   createdRole,
	})
}

func DeleteRole(c *fiber.Ctx) error {
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

	if !(orgOK || (userOK && roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.RoleWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.RoleFullAccess}))) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}

	roleExists, err := rolesRepo.GetRoleById(id)

	if err != nil || roleExists.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Role Not Found",
			"status":  "false",
		})
	}

	err = rolesRepo.DeleteRole(&roleExists)
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

func TestUserRole(c *fiber.Ctx) error {

	var role roleSchema.TestRole

	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	user, userOK := c.Locals("user").(userSchema.UserResponse)

	var roleExists model.Role
	var err error

	if role.RoleId != uuid.Nil {
		roleExists, err = rolesRepo.GetRoleById(role.RoleId)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Role Not Found",
				"status":  "false",
			})
		}
	} else if role.RoleName != "" {
		roleExists, err = rolesRepo.GetRoleByName(role.RoleName)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Role Not Found",
				"status":  "false",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Role",
			"status":  "error",
		})
	}

	if !userOK || !roles.UserHasRole(user.Roles, roleExists) {
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
