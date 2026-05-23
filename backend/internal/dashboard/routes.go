package dashboard

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	repository := NewRepository()
	service := NewService(repository)
	handler := NewHandler(service)

	api := r.Group("/api/dashboard")
	{
		api.GET("/summary", handler.GetSummary)
		api.GET("/todos-per-day", handler.GetTodosPerDay)
	}
}
