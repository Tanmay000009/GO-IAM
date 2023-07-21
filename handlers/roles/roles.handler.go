package rolesHandler

import (
	rolesRepo "balkantask/database/roles"
	"balkantask/model"
	orgSchema "balkantask/schemas/org"
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
	var role model.Role

	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.HasAnyRole(user.Roles, []roles.Role{roles.RoleWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.RoleFullAccess, roles.OrgReadAccess, roles.RoleReadAccess}) {
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

	exisitingRole, err := rolesRepo.GetRoleByName(role.Name)
	if err == nil && exisitingRole.Name == role.Name {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role already exists",
			"status":  "error",
		})
	}

	createdRole, err := rolesRepo.CreateRole(role)
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

	if !orgOK && !userOK && !roles.HasAnyRole(user.Roles, []roles.Role{roles.RoleWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.RoleFullAccess}) {
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

	err = rolesRepo.DeleteRoleById(id)
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
