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

const summaryCacheKey = "todos-summary"

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

	s.cache.Delete(ctx, summaryCacheKey)
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

	s.cache.Delete(ctx, fmt.Sprintf("todo-%s", id))
	s.cache.Delete(ctx, summaryCacheKey)
	return s.repository.Update(ctx, id, payload)
}

func (s *service) Delete(ctx context.Context, id string) error {
	s.cache.Delete(ctx, fmt.Sprintf("todo-%s", id))
	s.cache.Delete(ctx, summaryCacheKey)
	return s.repository.Delete(ctx, id)
}

func (s *service) BreakdownTask(
	ctx context.Context,
	todoID string,
	writer io.Writer,
	flusher http.Flusher,
) error {
	if data, ok := s.cache.Get(ctx, fmt.Sprint("todo-", todoID)); ok {
		for _, item := range data {
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

	s.cache.Set(ctx, fmt.Sprint("todo-", todoID), result, 15*time.Minute)
	return nil
}
