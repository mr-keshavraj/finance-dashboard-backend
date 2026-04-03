package repository

import (
	"context"
	"fmt"
	"time"

	"finance-dashboard-backend/db"
	"finance-dashboard-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getRecordCollection() *mongo.Collection {
	return db.GetCollection("records")
}

func CreateRecord(record models.Record) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := getRecordCollection().InsertOne(ctx, record)
	if err != nil {
		return primitive.NilObjectID, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	return insertedID, nil
}

func GetRecords(filter bson.M, page int64, limit int64) ([]models.Record, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	skip := (page - 1) * limit

	findOptions := options.Find().
		SetSort(bson.D{
			{Key: "date", Value: -1},
			{Key: "created_at", Value: -1},
		}).
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := getRecordCollection().Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var records []models.Record
	for cursor.Next(ctx) {
		var raw bson.M
		if err := cursor.Decode(&raw); err != nil {
			return nil, 0, err
		}

		record, err := mapToRecord(raw)
		if err != nil {
			return nil, 0, err
		}

		records = append(records, record)
	}

	total, err := getRecordCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func GetRecordByID(id string) (*models.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var raw bson.M
	err = getRecordCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&raw)
	if err != nil {
		return nil, err
	}

	record, err := mapToRecord(raw)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func UpdateRecord(id string, updatedData bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updatedData["updated_at"] = time.Now()

	_, err = getRecordCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updatedData},
	)

	return err
}

func DeleteRecord(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = getRecordCollection().DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func GetDashboardTotals() (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{
			{Key: "$group", Value: bson.M{
				"_id": nil,
				"total_income": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$type", "income"}},
							"$amount",
							0,
						},
					},
				},
				"total_expense": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$type", "expense"}},
							"$amount",
							0,
						},
					},
				},
			}},
		},
	}

	cursor, err := getRecordCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return bson.M{
			"total_income":  0,
			"total_expense": 0,
			"net_balance":   0,
		}, nil
	}

	totalIncome := toFloat64(results[0]["total_income"])
	totalExpense := toFloat64(results[0]["total_expense"])

	return bson.M{
		"total_income":  totalIncome,
		"total_expense": totalExpense,
		"net_balance":   totalIncome - totalExpense,
	}, nil
}

func GetCategoryBreakdown() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{
			{Key: "$group", Value: bson.M{
				"_id": bson.M{
					"category": "$category",
					"type":     "$type",
				},
				"total": bson.M{"$sum": "$amount"},
			}},
		},
		{
			{Key: "$project", Value: bson.M{
				"_id":      0,
				"category": "$_id.category",
				"type":     "$_id.type",
				"total":    1,
			}},
		},
		{
			{Key: "$sort", Value: bson.D{
				{Key: "total", Value: -1},
			}},
		},
	}

	cursor, err := getRecordCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetMonthlyTrends() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{
			{Key: "$group", Value: bson.M{
				"_id": bson.M{
					"year":  bson.M{"$year": "$date"},
					"month": bson.M{"$month": "$date"},
					"type":  "$type",
				},
				"total": bson.M{"$sum": "$amount"},
			}},
		},
		{
			{Key: "$project", Value: bson.M{
				"_id":   0,
				"year":  "$_id.year",
				"month": "$_id.month",
				"type":  "$_id.type",
				"total": 1,
			}},
		},
		{
			{Key: "$sort", Value: bson.D{
				{Key: "year", Value: 1},
				{Key: "month", Value: 1},
			}},
		},
	}

	cursor, err := getRecordCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetRecentActivity(limit int64) ([]models.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find().
		SetSort(bson.D{
			{Key: "date", Value: -1},
			{Key: "created_at", Value: -1},
		}).
		SetLimit(limit)

	cursor, err := getRecordCollection().Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []models.Record
	for cursor.Next(ctx) {
		var raw bson.M
		if err := cursor.Decode(&raw); err != nil {
			return nil, err
		}

		record, err := mapToRecord(raw)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func mapToRecord(raw bson.M) (models.Record, error) {
	var record models.Record

	if v, ok := raw["_id"]; ok {
		switch id := v.(type) {
		case primitive.ObjectID:
			record.ID = id
		case string:
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return record, fmt.Errorf("invalid _id string: %w", err)
			}
			record.ID = objID
		}
	}

	record.Amount = toFloat64(raw["amount"])
	record.Type = toString(raw["type"])
	record.Category = toString(raw["category"])
	record.Notes = toString(raw["notes"])

	date, err := toTime(raw["date"])
	if err != nil {
		return record, fmt.Errorf("invalid date: %w", err)
	}
	record.Date = date

	if v, ok := raw["created_by"]; ok {
		switch createdBy := v.(type) {
		case primitive.ObjectID:
			record.CreatedBy = createdBy
		case string:
			objID, err := primitive.ObjectIDFromHex(createdBy)
			if err != nil {
				return record, fmt.Errorf("invalid created_by string: %w", err)
			}
			record.CreatedBy = objID
		default:
			return record, fmt.Errorf("unsupported created_by type: %T", v)
		}
	}

	createdAt, err := toTime(raw["created_at"])
	if err != nil {
		return record, fmt.Errorf("invalid created_at: %w", err)
	}
	record.CreatedAt = createdAt

	updatedAt, err := toTime(raw["updated_at"])
	if err != nil {
		return record, fmt.Errorf("invalid updated_at: %w", err)
	}
	record.UpdatedAt = updatedAt

	return record, nil
}

func toFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	default:
		return 0
	}
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	default:
		return ""
	}
}

func toTime(value interface{}) (time.Time, error) {
	switch v := value.(type) {
	case time.Time:
		return v, nil
	case primitive.DateTime:
		return v.Time(), nil
	case string:
		return time.Parse(time.RFC3339, v)
	default:
		return time.Time{}, fmt.Errorf("unsupported time type: %T", value)
	}
}