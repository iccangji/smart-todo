package auth

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Module struct {
	Handler *Handler
}

func NewModule(db *mongo.Database) *Module {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	return &Module{
		Handler: handler,
	}
}
