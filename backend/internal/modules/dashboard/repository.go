package dashboard

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	GetSummary(ctx context.Context) (*SummaryResponse, error)
	GetTodosPerDay(ctx context.Context) ([]TodosPerDayResponse, error)
}

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetSummary(
	ctx context.Context,
) (*SummaryResponse, error) {
	now := time.Now()

	startOfDay := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	startOfWeek := startOfDay.AddDate(
		0,
		0,
		-int(now.Weekday()),
	)
	collection := r.db.Collection("todos")
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

				"pending": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$eq": []interface{}{
									"$completed",
									false,
								},
							},
							1,
							0,
						},
					},
				},
				"pending_low_priority": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []interface{}{
											"$completed",
											false,
										},
									},
									bson.M{
										"$eq": []interface{}{
											"$priority",
											0,
										},
									},
								},
							},
							1,
							0,
						},
					},
				},
				"pending_medium_priority": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []interface{}{
											"$completed",
											false,
										},
									},
									bson.M{
										"$eq": []interface{}{
											"$priority",
											1,
										},
									},
								},
							},
							1,
							0,
						},
					},
				},
				"pending_high_priority": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []interface{}{
											"$completed",
											false,
										},
									},
									bson.M{
										"$eq": []interface{}{
											"$priority",
											2,
										},
									},
								},
							},
							1,
							0,
						},
					},
				},
				"pending_urgent_priority": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []interface{}{
											"$completed",
											false,
										},
									},
									bson.M{
										"$eq": []interface{}{
											"$priority",
											3,
										},
									},
								},
							},
							1,
							0,
						},
					},
				},

				"completed_today": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []interface{}{
											"$completed",
											true,
										},
									},
									bson.M{
										"$gte": []interface{}{
											"$completed_at",
											startOfDay,
										},
									},
								},
							},
							1,
							0,
						},
					},
				},
				"completed_this_week": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []interface{}{
											"$completed",
											true,
										},
									},
									bson.M{
										"$gte": []interface{}{
											"$completed_at",
											startOfWeek,
										},
									},
								},
							},
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

	pendingCount := total - completed
	pendingTasks := Pending{
		Low:    data["pending_low_priority"].(int32),
		Medium: data["pending_medium_priority"].(int32),
		High:   data["pending_high_priority"].(int32),
		Urgent: data["pending_urgent_priority"].(int32),
	}

	completedToday := data["completed_today"].(int32)
	completedThisWeek := data["completed_this_week"].(int32)

	var completionRate float64

	if total > 0 {
		completionRate = (float64(completed) / float64(total)) * 100
	}

	return &SummaryResponse{
		Total:                total,
		CompletedCount:       completed,
		PendingCount:         pendingCount,
		PendingPriorityCount: pendingTasks,
		CompletionRate:       completionRate,
		CompletedToday:       completedToday,
		CompletedThisWeek:    completedThisWeek,
	}, nil
}

func (r *repository) GetTodosPerDay(
	ctx context.Context,
) ([]TodosPerDayResponse, error) {
	collection := r.db.Collection("todos")
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
