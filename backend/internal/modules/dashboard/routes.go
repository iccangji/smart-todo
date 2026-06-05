package dashboard

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, m *Module) {
	api := r.Group("/api/dashboard")
	{
		api.GET("/summary", m.Handler.GetSummary)
		api.GET("/summary/ai", m.Handler.Summarize)
		api.GET("/todos-per-day", m.Handler.GetTodosPerDay)
	}
}
