package dashboard

import (
	"backend/internal/ai"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Service interface {
	GetSummary(ctx context.Context) (*SummaryResponse, error)
	GetTodosPerDay(ctx context.Context) ([]TodosPerDayResponse, error)
	Summarize(ctx context.Context, cacheKey string, writer io.Writer, flusher http.Flusher) error
}

type service struct {
	repository Repository
	cache      Cache
}

func NewService(repository Repository, cache Cache) Service {
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

func (s *service) GetTodosPerDay(
	ctx context.Context,
) ([]TodosPerDayResponse, error) {
	return s.repository.GetTodosPerDay(ctx)
}

func (s *service) Summarize(
	ctx context.Context,
	cacheKey string,
	writer io.Writer,
	flusher http.Flusher,
) error {
	if s.cache != nil {
		if data, ok := s.cache.Get(cacheKey); ok {
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
	s.cache.Set(cacheKey, result, 15*time.Minute)

	return nil
}
