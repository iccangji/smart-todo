package auth

import (
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
	FindByID(ctx context.Context, userID string) (*User, error)
}

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) collection() string {
	return "users"
}

func (r *repository) Create(
	ctx context.Context,
	user *User,
) error {
	user.CreatedAt = time.Now()

	result, err := r.db.
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
	err := r.db.
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

func (r *repository) FindByID(
	ctx context.Context,
	userID string,
) (*User, error) {
	var user User
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	err = r.db.
		Collection(r.collection()).
		FindOne(ctx, bson.M{
			"_id": objID,
		}).
		Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
