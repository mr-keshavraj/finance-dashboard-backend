package controller

import (
	"strings"
	"time"

	"finance-dashboard-backend/constants"
	"finance-dashboard-backend/models"
	"finance-dashboard-backend/repository"
	"finance-dashboard-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)
	req.Role = strings.TrimSpace(strings.ToLower(req.Role))

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Name, email and password are required")
	}

	if len(req.Password) < 6 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Password must be at least 6 characters")
	}

	if req.Role == "" {
		req.Role = constants.RoleViewer
	}

	if req.Role != constants.RoleAdmin && req.Role != constants.RoleAnalyst && req.Role != constants.RoleViewer {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role")
	}

	existingUser, err := repository.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return utils.ErrorResponse(c, fiber.StatusConflict, "User with this email already exists")
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check existing user")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
	}

	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      req.Role,
		Status:    constants.StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repository.CreateUser(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "User registered successfully", fiber.Map{
		"name":   user.Name,
		"email":  user.Email,
		"role":   user.Role,
		"status": user.Status,
	})
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email and password are required")
	}

	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to login")
	}

	if user.Status != constants.StatusActive {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "User is inactive")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
	}

	token, err := utils.GenerateJWT(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Login successful", fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":     user.ID.Hex(),
			"name":   user.Name,
			"email":  user.Email,
			"role":   user.Role,
			"status": user.Status,
		},
	})
}