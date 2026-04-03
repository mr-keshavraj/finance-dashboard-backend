package controller

import (
	"strings"

	"finance-dashboard-backend/constants"
	"finance-dashboard-backend/repository"
	"finance-dashboard-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateRoleRequest struct {
	Role string `json:"role"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

func GetAllUsers(c *fiber.Ctx) error {
	users, err := repository.GetAllUsers()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch users")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Users fetched successfully", users)
}

func GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User ID is required")
	}

	user, err := repository.GetUserByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch user")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User fetched successfully", user)
}

func UpdateUserRole(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User ID is required")
	}

	var req UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	req.Role = strings.TrimSpace(strings.ToLower(req.Role))

	if req.Role == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Role is required")
	}

	if req.Role != constants.RoleAdmin && req.Role != constants.RoleAnalyst && req.Role != constants.RoleViewer {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role")
	}

	_, err := repository.GetUserByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch user")
	}

	if err := repository.UpdateUserRole(userID, req.Role); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user role")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User role updated successfully", fiber.Map{
		"user_id": userID,
		"role":    req.Role,
	})
}

func UpdateUserStatus(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User ID is required")
	}

	var req UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	req.Status = strings.TrimSpace(strings.ToLower(req.Status))

	if req.Status == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Status is required")
	}

	if req.Status != constants.StatusActive && req.Status != constants.StatusInactive {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid status")
	}

	_, err := repository.GetUserByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch user")
	}

	if err := repository.UpdateUserStatus(userID, req.Status); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user status")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User status updated successfully", fiber.Map{
		"user_id": userID,
		"status":  req.Status,
	})
}