package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Check(c *gin.Context) {
	result := h.service.Check(c.Request.Context())

	httpStatus := http.StatusOK

	for _, v := range result {
		if v != "up" {
			httpStatus = http.StatusServiceUnavailable
		}
	}

	c.JSON(httpStatus, gin.H{
		"status": result,
	})
}
