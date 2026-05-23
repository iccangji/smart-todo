package dashboard

type SummaryResponse struct {
	Total          int32   `json:"total"`
	Completed      int32   `json:"completed"`
	Pending        int32   `json:"pending"`
	CompletionRate float64 `json:"completion_rate"`
}

type TodosPerDayResponse struct {
	Date  string `json:"date"`
	Count int32  `json:"count"`
}
