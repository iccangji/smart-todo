package main

import (
	"context"
	"log"
	"time"

	"backend/internal/database"
	"backend/internal/todo"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	database.ConnectMongo()

	collection := database.DB.Collection("todos")

	var documents []interface{}

	for i := 0; i < 100; i++ {
		documents = append(documents, todo.Todo{
			Title:       gofakeit.Sentence(3),
			Description: gofakeit.Paragraph(1, 3, 5, " "),
			Completed:   gofakeit.Bool(),
			CreatedAt:   randomTime(),
			UpdatedAt:   time.Now(),
		})
	}

	ctx := context.Background()

	result, err := collection.InsertMany(ctx, documents)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Inserted %d todos\n", len(result.InsertedIDs))

	count, _ := collection.CountDocuments(ctx, bson.M{})

	log.Printf("Total todos in DB: %d\n", count)
}

func randomTime() time.Time {
	now := time.Now()

	return gofakeit.DateRange(
		now.AddDate(0, -6, 0),
		now,
	)
}
