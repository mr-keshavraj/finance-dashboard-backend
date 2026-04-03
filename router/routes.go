package router

import (
	"finance-dashboard-backend/constants"
	"finance-dashboard-backend/controller"
	"finance-dashboard-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Health check / base route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Finance Dashboard Backend is running",
		})
	})

	api := app.Group("/api")

	// Public auth routes
	auth := api.Group("/auth")
	auth.Post("/register", controller.Register)
	auth.Post("/login", controller.Login)

	// Protected routes (JWT required)
	protected := api.Group("/", middleware.Protected())

	protected.Get("/profile", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user_id": c.Locals("user_id"),
			"email":   c.Locals("email"),
			"role":    c.Locals("role"),
		})
	})

	// Example admin-only route
	protected.Get("/admin-only",
		middleware.Authorize(constants.RoleAdmin),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "Welcome Admin",
			})
		},
	)

	// User management (admin only)
	users := protected.Group("/users", middleware.Authorize(constants.RoleAdmin))
	users.Get("/", controller.GetAllUsers)
	users.Get("/:id", controller.GetUserByID)
	users.Patch("/:id/role", controller.UpdateUserRole)
	users.Patch("/:id/status", controller.UpdateUserStatus)

	// Records APIs (RBAC applied per route)
	records := protected.Group("/records")
	records.Post("/", middleware.Authorize(constants.RoleAdmin, constants.RoleAnalyst), controller.CreateRecord)
	records.Get("/", controller.GetAllRecords)
	records.Get("/:id", controller.GetRecordByID)
	records.Put("/:id", middleware.Authorize(constants.RoleAdmin, constants.RoleAnalyst), controller.UpdateRecord)
	records.Delete("/:id", middleware.Authorize(constants.RoleAdmin), controller.DeleteRecord)

	// Dashboard analytics endpoints
	dashboard := protected.Group("/dashboard")
	dashboard.Get("/summary", controller.GetDashboardSummary)
	dashboard.Get("/category", controller.GetCategoryBreakdown)
	dashboard.Get("/monthly", controller.GetMonthlyTrends)
	dashboard.Get("/recent", controller.GetRecentActivity)
}