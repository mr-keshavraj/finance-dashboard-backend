package utils

import "github.com/gofiber/fiber/v2"

func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	// Standard success response format
	return c.Status(statusCode).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	// Standard error response format
	return c.Status(statusCode).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}