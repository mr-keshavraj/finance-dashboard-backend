package main

import (
	"log"

	"finance-dashboard-backend/config"
	"finance-dashboard-backend/db"
	"finance-dashboard-backend/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load environment variables and app configuration (port, DB URI, etc.)
	config.LoadConfig()

	// Initialize database connection (MongoDB)
	db.ConnectDB()

	// Create a new Fiber app instance
	app := fiber.New()

	// Register all application routes (API endpoints)
	router.SetupRoutes(app)

	// Start the server on the configured port
	// log.Fatal ensures the app crashes if the server fails to start
	log.Fatal(app.Listen(":" + config.AppConfig.Port))
}