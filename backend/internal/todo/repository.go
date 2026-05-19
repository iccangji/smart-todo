package todo

import (
	"context"
	"strings"
	"time"

	"backend/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Create(ctx context.Context, todo *Todo) (*Todo, error)
	FindAll(ctx context.Context, query GetTodosQuery) ([]Todo, int64, error)
	FindByID(ctx context.Context, id string) (*Todo, error)
	Update(ctx context.Context, id string, payload bson.M) (*Todo, error)
	Delete(ctx context.Context, id string) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) collection() string {
	return "todos"
}

func (r *repository) Create(ctx context.Context, todo *Todo) (*Todo, error) {
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	result, err := database.DB.
		Collection(r.collection()).
		InsertOne(ctx, todo)

	todo.ID = result.InsertedID.(primitive.ObjectID)

	return todo, err
}

func (r *repository) FindAll(
	ctx context.Context,
	query GetTodosQuery,
) ([]Todo, int64, error) {
	var todos []Todo

	collection := database.DB.
		Collection(r.collection())

	filter := bson.M{}

	// Filter by search
	if query.Search != "" {
		filter["title"] = bson.M{
			"regex":    query.Search,
			"$options": "i",
		}
	}

	// Filter by completed
	if query.Completed != nil {
		filter["completed"] = *query.Completed
	}

	// Set order by filter
	skip := (query.Page - 1) * query.Limit
	sortOrder := -1
	if strings.ToLower(query.Order) == "asc" {
		sortOrder = 1
	}

	opts := options.Find()
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(query.Limit))
	opts.SetSort(bson.M{
		query.Sort: sortOrder,
	})

	// Get collections
	cursor, err := collection.Find(ctx, filter, opts)

	if err != nil {
		return nil, 0, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var todo Todo

		if err := cursor.Decode(&todo); err != nil {
			return nil, 0, err
		}

		todos = append(todos, todo)
	}

	total, err := collection.CountDocuments(ctx, filter)

	if err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Todo, error) {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	var todo Todo

	err = database.DB.
		Collection(r.collection()).
		FindOne(ctx, bson.M{
			"_id": objectID,
		}).
		Decode(&todo)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *repository) Update(ctx context.Context, id string, payload bson.M) (*Todo, error) {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	payload["updated_at"] = time.Now()
	after := options.After
	var updated Todo
	err = database.DB.
		Collection(r.collection()).
		FindOneAndUpdate(
			ctx,
			bson.M{
				"_id": objectID,
			},
			bson.M{
				"$set": payload,
			},
			&options.FindOneAndUpdateOptions{
				ReturnDocument: &after,
			},
		).
		Decode(&updated)

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	_, err = database.DB.
		Collection(r.collection()).
		DeleteOne(ctx, bson.M{
			"_id": objectID,
		})

	return err
}
