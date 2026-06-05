package main

import (
	"context"
	"log"
	"os"
	"time"

	"backend/internal/health"
	"backend/internal/infra/cache"
	"backend/internal/infra/database"
	"backend/internal/modules/auth"
	"backend/internal/modules/dashboard"
	"backend/internal/modules/todo"

	"github.com/gin-gonic/gin"
)

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

	// Init Cache
	redisAddr := os.Getenv("REDIS_HOST")
	redisDB := 0
	cache := cache.ConnectRedis(ctx, redisAddr, redisDB)

	r := gin.Default()

	todoModule := todo.NewModule(db, cache)
	todo.RegisterRoutes(r, todoModule)
	dashboardModule := dashboard.NewModule(db, cache)
	dashboard.RegisterRoutes(r, dashboardModule)
	authModule := auth.NewModule(db)
	auth.RegisterRoutes(r, authModule)

	healthService := health.NewService(db, cache)
	healthHandler := health.NewHandler(healthService)
	health.RegisterRoutes(r, healthHandler)

	port := os.Getenv("APP_PORT")
	r.Run(":" + port)
}
