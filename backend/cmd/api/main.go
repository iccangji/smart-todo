package main

import (
	"os"

	"backend/internal/auth"
	"backend/internal/database"
	"backend/internal/todo"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectMongo()
	r := gin.Default()
	todo.RegisterRoutes(r)
	auth.RegisterRoutes(r)
	port := os.Getenv("APP_PORT")
	r.Run(":" + port)
}
