package postgres

import (
	"context"
	"fmt"

	"github.com/Dorrrke/rent-analytcs-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, logger *zap.Logger) *Repository {
	return &Repository{pool: pool, logger: logger}
}

func (r *Repository) SaveEvent(ctx context.Context, e *domain.AnalyticsEvent) error {
	_, err := r.pool.Exec(ctx, `
	INSERT INTO analytics_events (id, subject, entity_id, payload, occurrated_at, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`, e.ID, e.Subject, e.EntityID, e.Payload, e.OccurratedAt, e.CreatedAt)
	if err != nil {
		return fmt.Errorf("SaveEvent: %w", err)
	}
	return nil
}

func (r *Repository) UpsertUserStats(ctx context.Context, e *domain.UserStats) error {
	_, err := r.pool.Exec(ctx, `
	INSERT INTO user_stats (user_id, total_logins, total_rents, total_spent, last_activity_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (user_id) DO UPDATE SET
		total_logins = user_stats.total_logins + EXCLUDED.total_logins,
		total_rents = user_stats.total_rents + EXCLUDED.total_rents,
		total_spent = user_stats.total_spent + EXCLUDED.total_spent,
		last_activity_at = GREATEST(user_stats.last_activity_at, EXCLUDED.last_activity_at),
		updated_at = EXCLUDED.updated_at
	`)
	if err != nil {
		return fmt.Errorf("UpsertUserStats: %w", err)
	}
	return nil
}

func (r *Repository) GetUserStats(ctx context.Context, userID string) (*domain.UserStats, error) {
	row := r.pool.QueryRow(ctx, `
	SELECT user_id, total_logins, total_rents, total_spent, last_activity_at, updated_at
	FROM user_stats WHERE user_id = $1
	`, userID)

	var s domain.UserStats
	if err := row.Scan(&s.UserID, &s.TotalLogins, &s.TotalRents, &s.TotalSpent, &s.LastActivityAt, &s.UpdatedAt); err != nil {
		return nil, fmt.Errorf("GetUserStats: %w", err)
	}
	return &s, nil
}

func (r *Repository) UpsertCarStats(ctx context.Context, s *domain.CarStats) error {
	_, err := r.pool.Exec(ctx, `
	INSERT INTO car_stats (car_id, total_rents, total_revenue, total_minutes, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (car_id) DO UPDATE SET
		total_rents = car_stats.total_rents + EXCLUDED.total_rents,
		total_revenue = car_stats.total_revenue + EXCLUDED.total_revenue,
		total_minutes = car_stats.total_minutes + EXCLUDED.total_minutes,
		updated_at = EXCLUDED.updated_at
	`, s.CarID, s.TotalRents, s.TotalRevenue, s.TotalMinutes, s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("UpsertCarStats: %w", err)
	}
	return nil
}

func (r *Repository) GetCarStats(ctx context.Context, carID string) (*domain.CarStats, error) {
	row := r.pool.QueryRow(ctx, `
	SELECT car_id, total_rents, total_revenue, total_minutes, updated_at
	FROM car_stats WHERE car_id = $1
	`, carID)

	var s domain.CarStats
	if err := row.Scan(&s.CarID, &s.TotalRents, &s.TotalRevenue, &s.TotalMinutes, &s.UpdatedAt); err != nil {
		return nil, fmt.Errorf("GetCarStats: %w", err)
	}
	return &s, nil
}

func (r *Repository) UpsertDailyMetric(ctx context.Context, m *domain.DailyMetric) error {
	_, err := r.pool.Exec(ctx, `
	INSERT INTO daily_metrics (date, new_users, logins, rents_started, rents_completed, revenue)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (date) DO UPDATE SET
		new_users = daily_metrics.new_users + EXCLUDED.new_users,
		logins = daily_metrics.logins + EXCLUDED.logins,
		rents_started = daily_metrics.rents_started + EXCLUDED.rents_started,
		rents_completed = daily_metrics.rents_completed + EXCLUDED.rents_completed,
		revenue = daily_metrics.revenue + EXCLUDED.revenue
	`, m.Date, m.NewUsers, m.Logins, m.RentsStarted, m.RentsComplited, m.Revenue)
	if err != nil {
		return fmt.Errorf("UpsertDailyMetric: %w", err)
	}
	return nil
}

func (r *Repository) GetDailyMetrics(ctx context.Context, f domain.DailyMetricFilter) ([]*domain.DailyMetric, error) {
	query := `SELECT date, new_users, logins, rents_started, rents_completed, revenue
	FROM daily_metrics WHERE 1=1`

	args := []any{}
	argN := 1

	if f.From != nil {
		query += fmt.Sprintf(" AND date >= $%d", argN)
		args = append(args, f.From)
		argN++
	}

	if f.To != nil {
		query += fmt.Sprintf(" AND date <= $%d", argN)
		args = append(args, f.To)
		argN++
	}

	query += " ORDER BY date DESC"
	if f.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argN)
		args = append(args, f.Limit)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetDailyMetric: %w", err)
	}
	defer rows.Close()

	var metrics []*domain.DailyMetric
	for rows.Next() {
		var m domain.DailyMetric
		if err := rows.Scan(&m.Date, &m.NewUsers, &m.Logins, &m.RentsStarted, &m.RentsComplited, &m.Revenue); err != nil {
			return nil, fmt.Errorf("GetDailyMetric scan: %w", err)
		}
		metrics = append(metrics, &m)
	}

	return metrics, rows.Err()
}
