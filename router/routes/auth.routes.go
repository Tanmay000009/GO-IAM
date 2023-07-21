package routes

import (
	authHandler "balkantask/handlers/auth"
	middleware "balkantask/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router) {

	userRouter := router.Group("/auth")

	userRouter.Get("/me", middleware.CheckJWT, authHandler.GetMe)
	userRouter.Post("/login", authHandler.SignInUser)
	userRouter.Post("/login/root", authHandler.SignInOrg)
	userRouter.Post("/logout", authHandler.LogoutUser)
	userRouter.Post("/signup", authHandler.SignUpOrg)
	userRouter.Delete("/:id", middleware.CheckJWT, authHandler.DeleteAccount)
}
