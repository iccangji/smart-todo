package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, m *Module) {
	api := r.Group("/api/auth")
	{
		api.POST("/register", m.Handler.Register)
		api.POST("/login", m.Handler.Login)
		api.POST("/refresh", m.Handler.Refresh)
	}
}
