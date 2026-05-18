package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectMongo() {
	err := godotenv.Load()

	if err != nil {
		log.Println(".env file not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")),
	)

	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database(os.Getenv("MONGO_DB"))

	log.Println("MongoDB connected")
}
