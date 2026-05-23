package dashboard

import "context"

type Service interface {
	GetSummary(ctx context.Context) (*SummaryResponse, error)
	GetTodosPerDay(ctx context.Context) ([]TodosPerDayResponse, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
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
