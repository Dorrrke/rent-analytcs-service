package domain

import (
	"context"
)

type AnalyticsRepository interface {
	SaveEvent(ctx context.Context, event *AnalyticsEvent) error

	UpsertUserStats(ctx context.Context, stats *UserStats) error
	GetUserStats(ctx context.Context, userID string) (*UserStats, error)

	UpsertCarStats(ctx context.Context, stats *CarStats) error
	GetCarStats(ctx context.Context, carID string) (*CarStats, error)

	UpsertDailyMetric(ctx context.Context, metric *DailyMetric) error
	GetDailyMetrics(ctx context.Context, filter DailyMetricFilter) ([]*DailyMetric, error)
}

type DailyMetricFilter struct {
	From  *string // ISO date string YYYY-MM-DD
	To    *string
	Limit int
}
