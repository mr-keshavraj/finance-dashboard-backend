package controller

import (
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"finance-dashboard-backend/models"
	"finance-dashboard-backend/repository"
	"finance-dashboard-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateRecordRequest struct {
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"`
	Category string  `json:"category"`
	Date     string  `json:"date"`
	Notes    string  `json:"notes"`
}

type UpdateRecordRequest struct {
	Amount   *float64 `json:"amount"`
	Type     *string  `json:"type"`
	Category *string  `json:"category"`
	Date     *string  `json:"date"`
	Notes    *string  `json:"notes"`
}

func CreateRecord(c *fiber.Ctx) error {
	var req CreateRecordRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	req.Type = strings.TrimSpace(strings.ToLower(req.Type))
	req.Category = strings.TrimSpace(req.Category)
	req.Notes = strings.TrimSpace(req.Notes)

	if req.Amount <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Amount must be greater than zero")
	}

	if req.Type != "income" && req.Type != "expense" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Type must be either income or expense")
	}

	if req.Category == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Category is required")
	}

	if req.Date == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Date is required")
	}

	parsedDate, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Date must be in RFC3339 format, e.g. 2026-04-03T00:00:00Z")
	}

	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User ID not found in token")
	}

	userObjectID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	record := models.Record{
		Amount:    req.Amount,
		Type:      req.Type,
		Category:  req.Category,
		Date:      parsedDate,
		Notes:     req.Notes,
		CreatedBy: userObjectID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	insertedID, err := repository.CreateRecord(record)
	if err != nil {
		log.Println("CreateRecord error:", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create record: "+err.Error())
	}

	record.ID = insertedID

	return utils.SuccessResponse(c, fiber.StatusCreated, "Record created successfully", record)
}

func GetAllRecords(c *fiber.Ctx) error {
	filter := bson.M{}

	recordType := strings.TrimSpace(strings.ToLower(c.Query("type")))
	category := strings.TrimSpace(c.Query("category"))
	startDate := strings.TrimSpace(c.Query("start_date"))
	endDate := strings.TrimSpace(c.Query("end_date"))

	page, err := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.ParseInt(c.Query("limit", "10"), 10, 64)
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	if recordType != "" {
		if recordType != "income" && recordType != "expense" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "type must be either income or expense")
		}
		filter["type"] = recordType
	}

	if category != "" {
		filter["category"] = category
	}

	if startDate != "" || endDate != "" {
		dateFilter := bson.M{}

		if startDate != "" {
			parsedStartDate, err := time.Parse(time.RFC3339, startDate)
			if err != nil {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "start_date must be in RFC3339 format")
			}
			dateFilter["$gte"] = parsedStartDate
		}

		if endDate != "" {
			parsedEndDate, err := time.Parse(time.RFC3339, endDate)
			if err != nil {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "end_date must be in RFC3339 format")
			}
			dateFilter["$lte"] = parsedEndDate
		}

		filter["date"] = dateFilter
	}

	records, total, err := repository.GetRecords(filter, page, limit)
	if err != nil {
		log.Println("GetAllRecords error:", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch records: "+err.Error())
	}

	totalPages := int64(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Records fetched successfully", fiber.Map{
		"items": records,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total_items": total,
			"total_pages": totalPages,
		},
	})
}

func GetRecordByID(c *fiber.Ctx) error {
	recordID := c.Params("id")
	if recordID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Record ID is required")
	}

	record, err := repository.GetRecordByID(recordID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Record not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch record")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Record fetched successfully", record)
}

func UpdateRecord(c *fiber.Ctx) error {
	recordID := c.Params("id")
	if recordID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Record ID is required")
	}

	_, err := repository.GetRecordByID(recordID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Record not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch record")
	}

	var req UpdateRecordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	updateData := bson.M{}

	if req.Amount != nil {
		if *req.Amount <= 0 {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Amount must be greater than zero")
		}
		updateData["amount"] = *req.Amount
	}

	if req.Type != nil {
		recordType := strings.TrimSpace(strings.ToLower(*req.Type))
		if recordType != "income" && recordType != "expense" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Type must be either income or expense")
		}
		updateData["type"] = recordType
	}

	if req.Category != nil {
		category := strings.TrimSpace(*req.Category)
		if category == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Category cannot be empty")
		}
		updateData["category"] = category
	}

	if req.Date != nil {
		parsedDate, err := time.Parse(time.RFC3339, *req.Date)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Date must be in RFC3339 format")
		}
		updateData["date"] = parsedDate
	}

	if req.Notes != nil {
		updateData["notes"] = strings.TrimSpace(*req.Notes)
	}

	if len(updateData) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No fields provided to update")
	}

	if err := repository.UpdateRecord(recordID, updateData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update record")
	}

	updatedRecord, err := repository.GetRecordByID(recordID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Record updated but failed to fetch updated record")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Record updated successfully", updatedRecord)
}

func DeleteRecord(c *fiber.Ctx) error {
	recordID := c.Params("id")
	if recordID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Record ID is required")
	}

	_, err := repository.GetRecordByID(recordID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Record not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch record")
	}

	if err := repository.DeleteRecord(recordID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete record")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Record deleted successfully", nil)
}