package dashboard

import (
	"backend/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetSummary(c *gin.Context) {
	data, err := h.service.GetSummary(
		c.Request.Context(),
	)

	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

	response.SuccessWithData(c, http.StatusOK, data)

}

func (h *Handler) GetTodosPerDay(c *gin.Context) {
	data, err := h.service.GetTodosPerDay(
		c.Request.Context(),
	)

	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

	response.SuccessWithData(c, http.StatusOK, data)
}
