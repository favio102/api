package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func LoadEnv() error {
	env := os.Getenv("ENVIRONMENT")
	if env == "development" {
		err := godotenv.Load()
		if err != nil {
			return fmt.Errorf("error loading .env file: %w", err)
		}
		fmt.Println("Loaded environment variables successfully from .env file.")
	}
	fmt.Println("Environment variables loaded successfully.")
	return nil
}

func ConnectDB() *mongo.Client {
	// Load environment variables
	err := LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	connectionString := os.Getenv("MONGODB_URL")
	if connectionString == "" {
		log.Fatal("MONGODB_URL is not set in the environment variables")
	}

	clientOptions := options.Client().ApplyURI(connectionString)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB: ", err)
	}

	fmt.Println("Connected to MongoDB Atlas!")
	Client = client
	return client
}
