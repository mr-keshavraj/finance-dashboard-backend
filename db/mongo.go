package db

import (
	"context"
	"log"
	"time"

	"finance-dashboard-backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	MongoClient = client
	MongoDatabase = client.Database(config.AppConfig.DBName)

	log.Println("MongoDB connected successfully")
}

func GetCollection(collectionName string) *mongo.Collection {
	return MongoDatabase.Collection(collectionName)
}