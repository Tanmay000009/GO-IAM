package routes

import (
	groupHandler "balkantask/handlers/group"
	middleware "balkantask/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupGroupRoutes(router fiber.Router) {
	rolesRouter := router.Group("/group", middleware.CheckJWT)

	rolesRouter.Get("/", groupHandler.GetAllGroups)
	rolesRouter.Get("/:id", groupHandler.GetGroupById)
	rolesRouter.Post("/", groupHandler.CreateGroup)
	rolesRouter.Post("/excel", groupHandler.SeedGroupsFromExcel)
	rolesRouter.Post("/csv", groupHandler.SeedGroupsFromCSV)
	rolesRouter.Delete("/:id", groupHandler.DeleteGroupById)
	rolesRouter.Post("/role/add", groupHandler.AddRoleToGroup)
	rolesRouter.Delete("/role/remove", groupHandler.DeleteRoleFromGroup)
}
