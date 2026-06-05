package dashboard

import "time"

type Todo struct {
	ID        string    `bson:"id" json:"id"`
	Title     string    `bson:"title" json:"title"`
	Completed bool      `bson:"completed" json:"completed"`
	Priority  int16     `bson:"priority" json:"priority"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
