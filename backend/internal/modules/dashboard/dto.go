package dashboard

type Pending struct {
	Low    int32 `json:"low"`
	Medium int32 `json:"medium"`
	High   int32 `json:"high"`
	Urgent int32 `json:"urgent"`
}
type SummaryResponse struct {
	Total                int32   `json:"total"`
	CompletedCount       int32   `json:"completed_count"`
	PendingCount         int32   `json:"pending_count"`
	PendingPriorityCount Pending `json:"pending_priority_count"`
	CompletionRate       float64 `json:"completion_rate"`
	CompletedToday       int32   `json:"completed_today"`
	CompletedThisWeek    int32   `json:"completed_this_week"`
}

type DayTodo struct {
	Date  string `json:"date"`
	Todos []Todo `json:"todos"`
}

type ThisWeekTodosResponse struct {
	Days []DayTodo `json:"days"`
}
