package middleware

import (
	"strings"

	"finance-dashboard-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authorization header is missing")
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid authorization format")
		}

		tokenString := headerParts[1]

		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}