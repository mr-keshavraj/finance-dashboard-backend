package controller

import (
	"strconv"

	"finance-dashboard-backend/repository"
	"finance-dashboard-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func GetDashboardSummary(c *fiber.Ctx) error {
	totals, err := repository.GetDashboardTotals()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch dashboard summary")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Dashboard summary fetched successfully", totals)
}

func GetCategoryBreakdown(c *fiber.Ctx) error {
	data, err := repository.GetCategoryBreakdown()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch category breakdown")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Category breakdown fetched successfully", data)
}

func GetMonthlyTrends(c *fiber.Ctx) error {
	data, err := repository.GetMonthlyTrends()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch monthly trends")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Monthly trends fetched successfully", data)
}

func GetRecentActivity(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "5")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = 5
	}

	data, err := repository.GetRecentActivity(limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch recent activity")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Recent activity fetched successfully", data)
}