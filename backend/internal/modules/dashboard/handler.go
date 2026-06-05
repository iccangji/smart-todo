package dashboard

import (
	"backend/internal/response"
	"fmt"
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

func (h *Handler) Summarize(c *gin.Context) {
	c.Writer.Header().Set(
		"Content-Type",
		"text/event-stream",
	)

	c.Writer.Header().Set(
		"Cache-Control",
		"no-cache",
	)

	c.Writer.Header().Set(
		"Connection",
		"keep-alive",
	)

	flusher, ok := c.Writer.(http.Flusher)

	if !ok {
		response.Error(c, http.StatusInternalServerError, "stream unsupported")
		return
	}

	err := h.service.Summarize(
		c.Request.Context(),
		c.Writer,
		flusher,
	)

	if err != nil {
		fmt.Fprintf(
			c.Writer,
			"data: error: %s\n\n",
			err.Error(),
		)

		flusher.Flush()
	}
}
