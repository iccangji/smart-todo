package health

import (
	"backend/internal/infra/cache"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	mongo *mongo.Database
	redis *cache.RedisCache
}

func NewService(mongo *mongo.Database, cache *cache.RedisCache) *Service {
	return &Service{
		mongo: mongo,
		redis: cache,
	}
}

func (s *Service) Check(ctx context.Context) map[string]string {

	status := map[string]string{
		"mongo": "down",
		"redis": "down",
	}

	if err := s.mongo.Client().Ping(ctx, nil); err == nil {
		status["mongo"] = "up"
	}

	if err := s.redis.Ping(ctx); err == nil {
		status["redis"] = "up"
	}

	return status
}
