package auth

import (
	"backend/internal/database"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) collection() string {
	return "users"
}

func (r *repository) Create(
	ctx context.Context,
	user *User,
) error {
	user.CreatedAt = time.Now()

	result, err := database.DB.
		Collection(r.collection()).
		InsertOne(ctx, user)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("email already exists")
		}
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return err
}

func (r *repository) FindByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	var user User
	err := database.DB.
		Collection(r.collection()).
		FindOne(ctx, bson.M{
			"email": email,
		}).
		Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
