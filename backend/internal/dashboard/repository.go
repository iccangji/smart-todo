package dashboard

import (
	"backend/internal/database"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	GetSummary(ctx context.Context) (*SummaryResponse, error)
	GetTodosPerDay(ctx context.Context) ([]TodosPerDayResponse, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetSummary(
	ctx context.Context,
) (*SummaryResponse, error) {
	collection := database.DB.Collection("todos")
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": nil,

				"total": bson.M{
					"$sum": 1,
				},

				"completed": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							"$completed",
							1,
							0,
						},
					},
				},
			},
		},
	}
	cursor, err := collection.Aggregate(
		ctx,
		pipeline,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var result []bson.M

	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return &SummaryResponse{}, nil
	}

	data := result[0]

	total := data["total"].(int32)
	completed := data["completed"].(int32)

	pending := total - completed

	var completionRate float64

	if total > 0 {
		completionRate = (float64(completed) / float64(total)) * 100
	}

	return &SummaryResponse{
		Total:          total,
		Completed:      completed,
		Pending:        pending,
		CompletionRate: completionRate,
	}, nil
}

func (r *repository) GetTodosPerDay(
	ctx context.Context,
) ([]TodosPerDayResponse, error) {
	collection := database.DB.Collection("todos")
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$created_at",
					},
				},

				"count": bson.M{
					"$sum": 1,
				},
			},
		},

		{
			"$sort": bson.M{
				"_id": 1,
			},
		},
	}
	cursor, err := collection.Aggregate(
		ctx,
		pipeline,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var rawResults []bson.M

	if err := cursor.All(ctx, &rawResults); err != nil {
		return nil, err
	}

	var results []TodosPerDayResponse

	for _, item := range rawResults {

		results = append(results, TodosPerDayResponse{
			Date:  item["_id"].(string),
			Count: item["count"].(int32),
		})
	}

	return results, nil
}
