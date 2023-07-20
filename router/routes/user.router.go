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
	userRouter.Put("/:id", userHandler.UpdateUser)
	userRouter.Delete("/:id", userHandler.DeleteUser)
}
