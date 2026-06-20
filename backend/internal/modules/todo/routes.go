package todo

import (
	"backend/internal/modules/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, m *Module) {
	api := r.Group("/api/todos")
	api.Use(auth.AuthMiddleware())
	{
		api.POST("", m.Handler.Create)
		api.GET("", m.Handler.GetAll)
		api.GET("/:id", m.Handler.GetByID)
		api.GET("/:id/breakdown", m.Handler.Breakdown)
		api.PUT("/:id", m.Handler.Update)
		api.DELETE("/:id", m.Handler.Delete)
	}
}
