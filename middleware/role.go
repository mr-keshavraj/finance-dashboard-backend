package middleware

import (
	"finance-dashboard-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func Authorize(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok || userRole == "" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "User role not found")
		}

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied")
	}
}