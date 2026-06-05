package todo

type CreateTodoRequest struct {
	Title       string   `json:"title" binding:"required,min=3"`
	Description string   `json:"description" binding:"required"`
	Priority    Priority `json:"priority" binding:"min=0,max=3"`
}

type UpdateTodoRequest struct {
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Completed   *bool     `json:"completed"`
	Priority    *Priority `json:"priority"`
}
