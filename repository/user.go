package repository

import (
	"context"
	"time"

	"finance-dashboard-backend/db"
	"finance-dashboard-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUserCollection() *mongo.Collection {
	return db.GetCollection("users")
}

func CreateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := getUserCollection().InsertOne(ctx, user)
	return err
}

func GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := getUserCollection().FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id) // convert string → ObjectID
	if err != nil {
		return nil, err
	}

	var user models.User
	err = getUserCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := getUserCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func UpdateUserRole(userID string, role string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = getUserCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"role":       role,
				"updated_at": time.Now(), // keep track of updates
			},
		},
	)

	return err
}

func UpdateUserStatus(userID string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = getUserCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	)

	return err
}