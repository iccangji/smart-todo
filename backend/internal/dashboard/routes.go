package dashboard

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	repository := NewRepository()
	memCache := NewMemoryCache()
	service := NewService(repository, memCache)
	handler := NewHandler(service)

	api := r.Group("/api/dashboard")
	{
		api.GET("/summary", handler.GetSummary)
		api.GET("/summary/ai", handler.Summarize)
		api.GET("/todos-per-day", handler.GetTodosPerDay)
	}
}
