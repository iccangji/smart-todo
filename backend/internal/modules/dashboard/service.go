package dashboard

import (
	"backend/internal/ai"
	"backend/internal/infra/cache"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Service interface {
	GetSummary(ctx context.Context) (*SummaryResponse, error)
	GetThisWeekTodos(ctx context.Context) (*ThisWeekTodosResponse, error)
	Summarize(ctx context.Context, writer io.Writer, flusher http.Flusher) error
}

type service struct {
	repository Repository
	cache      cache.Cache
}

func NewService(repository Repository, cache cache.Cache) Service {
	return &service{
		repository: repository,
		cache:      cache,
	}
}

func (s *service) GetSummary(
	ctx context.Context,
) (*SummaryResponse, error) {
	var summaryDashboard SummaryResponse
	if ok := s.cache.Get(ctx, cache.DashboardSummaryCacheKey, &summaryDashboard); ok {
		return &summaryDashboard, nil
	}
	summary, err := s.repository.GetSummary(ctx)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, cache.DashboardSummaryCacheKey, summary, 15*time.Minute)
	return summary, nil

}

func (s *service) GetThisWeekTodos(
	ctx context.Context,
) (*ThisWeekTodosResponse, error) {
	var thisWeekTodosDashboard ThisWeekTodosResponse
	if ok := s.cache.Get(ctx, cache.DashboardThisWeekTodosCacheKey, &thisWeekTodosDashboard); ok {
		return &thisWeekTodosDashboard, nil
	}
	thisWeekTodos, err := s.repository.GetThisWeekTodos(ctx)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, cache.DashboardThisWeekTodosCacheKey, thisWeekTodos, 15*time.Minute)
	return thisWeekTodos, nil
}

func (s *service) Summarize(
	ctx context.Context,
	writer io.Writer,
	flusher http.Flusher,
) error {
	var summarizeResult []string
	if s.cache != nil {
		if ok := s.cache.Get(ctx, cache.SummaryCacheKey, &summarizeResult); ok {
			fmt.Println("Retrieved from cache", cache.SummaryCacheKey)
			for _, item := range summarizeResult {
				fmt.Fprintf(writer, "data: %s\n\n", item)
				flusher.Flush()
			}
			return nil
		}
	}

	summary, err := s.repository.GetSummary(ctx)
	if err != nil {
		return err
	}
	jsonData, _ := json.MarshalIndent(summary, "", "  ")

	var result []string
	ai.StreamResponse(
		ai.Summarize,
		jsonData,
		writer,
		flusher,
		func(content string) {
			result = append(result, content)
		},
	)
	s.cache.Set(ctx, cache.SummaryCacheKey, result, 15*time.Minute)

	return nil
}
