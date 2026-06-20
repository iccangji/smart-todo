package todo

import (
	"backend/internal/ai"
	"backend/internal/infra/cache"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Create(ctx context.Context, userID string, req CreateTodoRequest) (*Todo, error)
	GetAll(ctx context.Context, query GetTodosQuery) ([]Todo, int64, error)
	GetByID(ctx context.Context, id string) (*Todo, error)
	Update(ctx context.Context, id string, req UpdateTodoRequest) (*Todo, error)
	Delete(ctx context.Context, id string) error

	// AI Features
	BreakdownTask(
		ctx context.Context,
		todoID string,
		writer io.Writer,
		flusher http.Flusher,
	) error

	InvalidateCache(ctx context.Context, todoID *string)
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

func (s *service) Create(
	ctx context.Context,
	userID string,
	req CreateTodoRequest,
) (*Todo, error) {

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	todo := &Todo{
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		UserID:      userObjectID,
		Priority:    req.Priority,
	}

	s.InvalidateCache(ctx, nil)
	return s.repository.Create(ctx, todo)
}

func (s *service) GetAll(ctx context.Context, filter GetTodosQuery) ([]Todo, int64, error) {
	return s.repository.FindAll(ctx, filter)
}

func (s *service) GetByID(ctx context.Context, id string) (*Todo, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id string, req UpdateTodoRequest) (*Todo, error) {
	payload := utils.StructToBsonM(req)

	if len(payload) == 0 {
		return nil, errors.New("no fields to update")
	}
	if payload["priority"] != nil {
		switch payload["priority"] {
		case Low, Medium, High, Urgent:
		default:
			return nil, errors.New("invalid priority")
		}
	}

	if title, ok := utils.GetString(payload, "title"); ok {
		payload["title"] = title
	}

	if desc, ok := utils.GetString(payload, "description"); ok {
		payload["description"] = desc
	}

	if completed, ok := utils.GetBool(payload, "completed"); ok {
		payload["completed"] = completed
		if completed {
			payload["completed_at"] = time.Now()
		} else {
			payload["completed_at"] = nil
		}
	}

	s.InvalidateCache(ctx, &id)
	return s.repository.Update(ctx, id, payload)
}

func (s *service) Delete(ctx context.Context, id string) error {
	s.InvalidateCache(ctx, &id)
	return s.repository.Delete(ctx, id)
}

func (s *service) BreakdownTask(
	ctx context.Context,
	todoID string,
	writer io.Writer,
	flusher http.Flusher,
) error {
	var taskBreakdownResult []string
	if ok := s.cache.Get(ctx, fmt.Sprint(cache.TodoBreakdownCacheKey, todoID), &taskBreakdownResult); ok {
		for _, item := range taskBreakdownResult {
			fmt.Fprintf(writer, "data: %s\n\n", item)
			flusher.Flush()
		}
		return nil
	}

	todo, err := s.repository.FindByID(ctx, todoID)
	if err != nil {
		return err
	}

	jsonData, _ := json.MarshalIndent(todo, "", "  ")

	var result []string
	ai.StreamResponse(
		ai.Breakdown,
		jsonData,
		writer,
		flusher,
		func(content string) {
			result = append(result, strings.TrimSpace(content))
		},
	)

	s.cache.Set(ctx, fmt.Sprint(cache.TodoBreakdownCacheKey, todoID), result, 15*time.Minute)
	return nil
}

func (s *service) InvalidateCache(ctx context.Context, todoID *string) {
	if todoID != nil {
		s.cache.Delete(ctx, fmt.Sprint(cache.TodoBreakdownCacheKey, *todoID))
	}
	s.cache.Delete(ctx, cache.SummaryCacheKey)
	s.cache.Delete(ctx, cache.DashboardSummaryCacheKey)
	s.cache.Delete(ctx, cache.DashboardThisWeekTodosCacheKey)
}
