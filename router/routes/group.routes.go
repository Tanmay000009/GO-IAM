package routes

import (
	groupHandler "balkantask/handlers/group"
	middleware "balkantask/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupGroupRoutes(router fiber.Router) {
	groupRouter := router.Group("/group", middleware.CheckJWT)

	groupRouter.Get("/", groupHandler.GetAllGroups)
	groupRouter.Get("/:id", groupHandler.GetGroupById)
	groupRouter.Post("/", groupHandler.CreateGroup)
	groupRouter.Post("/test", groupHandler.TestUserGroup)
	groupRouter.Post("/excel", groupHandler.SeedGroupsFromExcel)
	groupRouter.Post("/csv", groupHandler.SeedGroupsFromCSV)
	groupRouter.Delete("/:id", groupHandler.DeleteGroupById)
	groupRouter.Post("/role/add", groupHandler.AddRoleToGroup)
	groupRouter.Delete("/role/remove", groupHandler.DeleteRoleFromGroup)
}
