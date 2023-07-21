package routes

import (
	userHandler "balkantask/handlers/user"
	middleware "balkantask/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(router fiber.Router) {

	userRouter := router.Group("/user", middleware.CheckJWT)

	userRouter.Get("/", userHandler.GetUsers)
	userRouter.Get("/:id", userHandler.GetUserById)
	userRouter.Post("/", userHandler.CreateUser)
	userRouter.Post("/seed", userHandler.SeedUsersFromExcel)
	userRouter.Put("/:id", userHandler.UpdateUser)
	userRouter.Delete("/:id", userHandler.DeleteUser)
	userRouter.Post("/role/add", userHandler.AddRoleToUser)
	userRouter.Delete("/role/remove", userHandler.DeleteRoleFromUser)
	userRouter.Post("/group/add", userHandler.AddGroupToUser)
	userRouter.Delete("/group/remove", userHandler.DeleteGroupFromUser)
	userRouter.Put("/deactivate/:id", userHandler.DeactivateUser)
	userRouter.Put("/reactivate/:id", userHandler.ReactivateUser)
	userRouter.Put("/update/password", userHandler.ChangePassword)
}
