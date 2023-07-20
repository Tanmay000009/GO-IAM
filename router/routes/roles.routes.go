package routes

import (
	rolesHandler "balkantask/handlers/roles"
	middleware "balkantask/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRolesRoutes(app *fiber.App) {
	roles := app.Group("/roles", middleware.CheckJWT)

	roles.Get("/", rolesHandler.GetAllRoles)
	roles.Get("/:id", rolesHandler.GetRoleById)
	roles.Post("/", rolesHandler.CreateRole)
	roles.Delete("/:id", rolesHandler.DeleteRole)
}
