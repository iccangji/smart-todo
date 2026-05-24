package todo

import (
	"backend/internal/ai"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
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

	return s.repository.Update(ctx, id, payload)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) BreakdownTask(
	ctx context.Context,
	todoID string,
	writer io.Writer,
	flusher http.Flusher,
) error {
	todo, err := s.repository.FindByID(ctx, todoID)
	if err != nil {
		return err
	}

	// Stream from database cache
	if len(todo.Breakdown) > 0 {
		for _, item := range todo.Breakdown {
			fmt.Fprintf(writer, "data: %s\n\n", item)
			flusher.Flush()
		}
		return nil
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

	todo.Breakdown = result

	_, err = s.repository.Update(ctx, todoID, bson.M{
		"breakdown": todo.Breakdown,
	})
	if err != nil {
		return err
	}
	return nil
}
