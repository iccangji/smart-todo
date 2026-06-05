package main

import (
	"context"
	"log"
	"os"
	"time"

	"backend/internal/infra/database"
	"backend/internal/modules/todo"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getUsers(ctx context.Context) ([]primitive.ObjectID, error) {
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	db, err := database.ConnectMongo(ctx, uri, dbName)
	if err != nil {
		log.Fatal(err)
	}
	collection := db.Collection("users")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var userIDs []primitive.ObjectID
	for cursor.Next(ctx) {
		var u struct {
			ID primitive.ObjectID `bson:"_id"`
		}

		if err := cursor.Decode(&u); err != nil {
			return nil, err
		}

		userIDs = append(userIDs, u.ID)
	}
	return userIDs, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Init Database
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	db, err := database.ConnectMongo(ctx, uri, dbName)
	if err != nil {
		log.Fatal(err)
	}

	userIDs, err := getUsers(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if len(userIDs) == 0 {
		log.Fatal("no users found, create users first")
	}

	collection := db.Collection("todos")

	var docs []interface{}

	// Seed for 100 todos
	for i := 0; i < 20; i++ {

		randomUser := userIDs[gofakeit.Number(0, len(userIDs)-1)]

		docs = append(docs, todo.Todo{
			Title:       gofakeit.Sentence(3),
			Description: gofakeit.Paragraph(1, 3, 5, " "),
			Completed:   false,
			Priority:    todo.Priority((gofakeit.Number(0, 3))),
			UserID:      randomUser,

			CreatedAt: randomTime(),
			UpdatedAt: time.Now(),
		})
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Inserted %d todos\n", len(result.InsertedIDs))
}

func randomTime() time.Time {
	now := time.Now()

	return gofakeit.DateRange(
		now.AddDate(0, -6, 0),
		now,
	)
}
