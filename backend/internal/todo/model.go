package todo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Priority int16

const (
	Low    Priority = iota // 0
	Medium                 // 1
	High                   // 2
	Urgent                 // 3
)

type Todo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Priority    Priority           `bson:"priority" json:"priority"`
	Completed   bool               `bson:"completed" json:"completed"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Breakdown   []string           `bson:"breakdown,omitempty" json:"breakdown,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	CompletedAt *time.Time         `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
}
