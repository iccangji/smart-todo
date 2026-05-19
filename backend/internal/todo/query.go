package todo

type GetTodosQuery struct {
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	Search    string `form:"search"`
	Completed *bool  `form:"completed"`
	Sort      string `form:"sort"`
	Order     string `form:"order"`
}

func (q *GetTodosQuery) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}

	if q.Limit <= 0 {
		q.Limit = 10
	}

	if q.Limit > 100 {
		q.Limit = 100
	}

	if q.Sort == "" {
		q.Sort = "created_at"
	}

	if q.Order == "" {
		q.Order = "desc"
	}
}
