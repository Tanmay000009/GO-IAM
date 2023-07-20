package routes

import (
	userHandler "balkantask/handlers/user"
	middleware "balkantask/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(router fiber.Router) {

	userRouter := router.Group("/user")

	userRouter.Get("/", userHandler.GetUsers)
	userRouter.Get("/:id", userHandler.GetUserById)
	userRouter.Post("/", middleware.CheckJWT, userHandler.CreateUser)
	userRouter.Put("/:id", userHandler.UpdateUser)
	userRouter.Delete("/:id", userHandler.DeleteUser)
}
