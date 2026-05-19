package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	repository := NewRepository()
	service := NewService(repository)
	handler := NewHandler(service)

	api := r.Group("/api/auth")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
		api.POST("/refresh", handler.Refresh)
	}
}
