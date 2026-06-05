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

const summaryCacheKey = "todos-summary"

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
	return s.repository.GetSummary(ctx)
}

func (s *service) GetThisWeekTodos(
	ctx context.Context,
) (*ThisWeekTodosResponse, error) {
	return s.repository.GetThisWeekTodos(ctx)
}

func (s *service) Summarize(
	ctx context.Context,
	writer io.Writer,
	flusher http.Flusher,
) error {
	if s.cache != nil {
		if data, ok := s.cache.Get(ctx, summaryCacheKey); ok {
			fmt.Println("Retrieved from cache", data)
			for _, item := range data {
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
	s.cache.Set(ctx, summaryCacheKey, result, 15*time.Minute)

	return nil
}
