package routes

import (
	rolesHandler "balkantask/handlers/roles"
	middleware "balkantask/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRolesRoutes(router fiber.Router) {
	roles := router.Group("/roles", middleware.CheckJWT)

	roles.Get("/", rolesHandler.GetAllRoles)
	roles.Get("/:id", rolesHandler.GetRoleById)
	roles.Post("/", middleware.CheckJWT, rolesHandler.CreateRole)
	roles.Delete("/:id", middleware.CheckJWT, rolesHandler.DeleteRole)
}
