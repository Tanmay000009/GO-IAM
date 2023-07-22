package routes

import (
	rolesHandler "balkantask/handlers/roles"
	middleware "balkantask/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRolesRoutes(router fiber.Router) {
	roles := router.Group("/roles")

	roles.Get("/", rolesHandler.GetAllRoles)
	roles.Get("/:id", rolesHandler.GetRoleById)
	roles.Post("/", middleware.CheckJWT, rolesHandler.CreateRole)
	roles.Post("/test", middleware.CheckJWT, rolesHandler.TestUserRole)
	roles.Post("/seed", middleware.CheckJWT, rolesHandler.SeedRoles)
	roles.Delete("/:id", middleware.CheckJWT, rolesHandler.DeleteRole)
}
