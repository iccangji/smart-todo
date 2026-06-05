package dashboard

import (
	"backend/internal/infra/cache"

	"go.mongodb.org/mongo-driver/mongo"
)

type Module struct {
	Handler *Handler
}

func NewModule(db *mongo.Database, cache cache.Cache) *Module {

	repo := NewRepository(db)
	service := NewService(repo, cache)
	handler := NewHandler(service)

	return &Module{
		Handler: handler,
	}
}
