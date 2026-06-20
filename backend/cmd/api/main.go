package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		<-stop
		log.Println("shutdown signal received")

		cancel()
	}()

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
	dashboard.StartScheduler(ctx, dashboardModule)

	authModule := auth.NewModule(db)
	auth.RegisterRoutes(r, authModule)

	healthService := health.NewService(db, cache)
	healthHandler := health.NewHandler(healthService)
	health.RegisterRoutes(r, healthHandler)

	port := os.Getenv("APP_PORT")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		<-ctx.Done()

		log.Println("shutting down http server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Println("server shutdown error:", err)
		}
	}()

	log.Println("server running on :" + port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	log.Println("server exited")
}
