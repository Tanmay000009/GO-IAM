package rolesHandler

import (
	rolesRepo "balkantask/database/roles"
	"balkantask/model"
	orgSchema "balkantask/schemas/org"
	userSchema "balkantask/schemas/user"
	roles "balkantask/utils"

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

	if !orgOK && !userOK && !roles.HasAnyRole(user.Roles, roles.RoleWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.RoleFullAccess) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
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
	var role model.Role

	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.HasAnyRole(user.Roles, roles.RoleWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.RoleFullAccess) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}

	err := rolesRepo.DeleteRoleById(role.ID)
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
