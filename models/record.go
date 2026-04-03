package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Record represents a financial transaction (income or expense)
type Record struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"` // MongoDB unique ID
	Amount    float64            `bson:"amount" json:"amount"`   // Transaction amount
	Type      string             `bson:"type" json:"type"`       // income / expense
	Category  string             `bson:"category" json:"category"` // Category like food, salary, rent, etc.
	Date      time.Time          `bson:"date" json:"date"`       // Date of transaction
	Notes     string             `bson:"notes" json:"notes"`     // Optional description
	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by"` // Reference to user who created it
	CreatedAt time.Time          `bson:"created_at" json:"created_at"` // Record creation timestamp
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"` // Last update timestamp
}