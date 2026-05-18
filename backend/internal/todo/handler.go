package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	todo, err := h.service.Create(c.Request.Context(), req)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithData(c, http.StatusCreated, todo)
}

func (h *Handler) GetAll(c *gin.Context) {
	todos, err := h.service.GetAll(c.Request.Context())

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithData(c, http.StatusOK, todos)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	todo, err := h.service.GetByID(c.Request.Context(), id)

	if err != nil {
		response.Error(c, http.StatusNotFound, "Todo not found")
		return
	}

	response.SuccessWithData(c, http.StatusOK, todo)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.service.Update(c.Request.Context(), id, req)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithData(c, http.StatusOK, updated)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.service.Delete(c.Request.Context(), id)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Todo deleted")
}
