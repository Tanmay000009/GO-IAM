package routes

import (
	taskHandler "balkantask/handlers/task"
	middleware "balkantask/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupTaskRoutes(router fiber.Router) {
	taskRouter := router.Group("/task", middleware.CheckJWT)

	taskRouter.Get("/", taskHandler.GetAllTasks)
	taskRouter.Get("/:id", taskHandler.GetTaskById)
	taskRouter.Post("/", taskHandler.CreateTask)
	taskRouter.Post("/test", taskHandler.TestUserTask)
	taskRouter.Post("/excel", taskHandler.SeedTasksFromExcel)
	taskRouter.Post("/csv", taskHandler.SeedTasksFromCSV)
	taskRouter.Delete("/:id", taskHandler.DeleteTaskById)
	taskRouter.Post("/role/add", taskHandler.AddRoleToTask)
	taskRouter.Delete("/role/remove", taskHandler.DeleteRoleFromTask)
}
