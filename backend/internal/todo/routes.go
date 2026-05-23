package todo

import (
	"backend/internal/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	repository := NewRepository()
	service := NewService(repository)
	handler := NewHandler(service)

	api := r.Group("/api/todos")
	api.Use(auth.AuthMiddleware())
	{
		api.POST("", handler.Create)
		api.GET("", handler.GetAll)
		api.GET("/:id", handler.GetByID)
		api.GET("/:id/breakdown", handler.Breakdown)
		api.PUT("/:id", handler.Update)
		api.DELETE("/:id", handler.Delete)
	}

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "todos api is healthy",
		})
	})
}
