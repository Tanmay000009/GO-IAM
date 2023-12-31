package middleware

import (
	orgrepository "balkantask/database/org"
	userRepo "balkantask/database/user"
	orgSchema "balkantask/schemas/org"
	userSchema "balkantask/schemas/user"
	constants "balkantask/utils"

	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CheckJWT(c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if len(authorization) > 0 {
		tokenString = authorization
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "false", "message": "You are not logged in"})
	}

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "false", "message": fmt.Sprintf("Invalid token: %v", err)})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "false", "message": "Invalid token"})

	}

	id_uuid, err := uuid.Parse(fmt.Sprint(claims["sub"]))

	if err != nil {

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Something Went Wrong"})
	}

	user, err := userRepo.FindUserByIdWithPassword(id_uuid)
	org, orgErr := orgrepository.FindOrgById(id_uuid)
	if err != nil && orgErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": "Invalid token"})
	}

	if (user.Org != nil && user.Org.AccountStatus == constants.DELETED) || org.AccountStatus == constants.DELETED {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "false", "message": "Account does not exist"})
	}

	if user.AccountStatus == constants.DEACTIVATED {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "false", "message": "Account deactivated"})
	}

	if user.ID.String() == claims["sub"] {
		c.Locals("user", userSchema.MapUserRecord(&user))
	} else if org.ID.String() == claims["sub"] {
		c.Locals("org", orgSchema.MapOrgRecord(&org))
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "false", "message": "Invalid token"})
	}

	return c.Next()
}
