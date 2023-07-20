package routes

import (
	authHandler "balkantask/handlers/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router) {

	userRouter := router.Group("/auth")

	userRouter.Get("/me", authHandler.GetMe)
	userRouter.Post("/login", authHandler.SignInUser)
	userRouter.Post("/logout", authHandler.LogoutUser)
	userRouter.Post("/signup", authHandler.SignUpOrg)
}
