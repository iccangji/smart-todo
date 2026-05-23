package todo

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required,min=3"`
	Description string `json:"description" binding:"required"`
}

type UpdateTodoRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

type BreakdownResponse struct {
	Title string   `json:"title"`
	Steps []string `json:"steps"`
}
