package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Dorrrke/rent-analytcs-service/internal/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AanalyticsService struct {
	repo   domain.AnalyticsRepository
	logger *zap.Logger
}

func NewAnalyticsService(repo domain.AnalyticsRepository, logger *zap.Logger) *AanalyticsService {
	return &AanalyticsService{
		repo:   repo,
		logger: logger,
	}
}

// --- User events ----------------

func (s *AanalyticsService) HandleUserRegistered(ctx context.Context, e domain.UserRegisteredEvent) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

	if err := s.repo.SaveEvent(ctx, &domain.AnalyticsEvent{
		ID:           uuid.NewString(),
		Subject:      domain.SubjectUserRegistered,
		EntityID:     e.UserID,
		Payload:      payload,
		OccurratedAt: e.CreatedAt,
		CreatedAt:    time.Now().UTC(),
	}); err != nil {
		return err
	}

	return s.repo.UpsertDailyMetric(ctx, &domain.DailyMetric{
		Date:     truncateToDay(e.CreatedAt),
		NewUsers: 1,
	})
}

func (s *AanalyticsService) HandleUserLoggedIn(ctx context.Context, e domain.UserLoggedInEvent) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

	if err := s.repo.SaveEvent(ctx, &domain.AnalyticsEvent{
		ID:           uuid.NewString(),
		Subject:      domain.SubjectUserLoggedIn,
		EntityID:     e.UserID,
		Payload:      payload,
		OccurratedAt: e.At,
		CreatedAt:    time.Now().UTC(),
	}); err != nil {
		return err
	}

	stats := &domain.UserStats{
		UserID:         e.UserID,
		TotalLogins:    1,
		LastActivityAt: e.At,
		UpdatedAt:      time.Now().UTC(),
	}
	if err := s.repo.UpsertUserStats(ctx, stats); err != nil {
		return err
	}

	return s.repo.UpsertDailyMetric(ctx, &domain.DailyMetric{
		Date:   truncateToDay(e.At),
		Logins: 1,
	})
}

func (s *AanalyticsService) HandleTokenRefreshed(ctx context.Context, e domain.TokenRefreshedEvent) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return s.repo.SaveEvent(ctx, &domain.AnalyticsEvent{
		ID:           uuid.NewString(),
		Subject:      domain.SubjectTokenRefreshed,
		EntityID:     e.UserID,
		Payload:      payload,
		OccurratedAt: e.At,
		CreatedAt:    time.Now().UTC(),
	})
}

// --- Car events ----------------

func (s *AanalyticsService) HandleCarAdded(ctx context.Context, e domain.CarAddedEvent) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return s.repo.SaveEvent(ctx, &domain.AnalyticsEvent{
		ID:           uuid.NewString(),
		Subject:      domain.SubjectCarAdded,
		EntityID:     e.CarID,
		Payload:      payload,
		OccurratedAt: e.AddeddAt,
		CreatedAt:    time.Now().UTC(),
	})
}

func (s *AanalyticsService) HandleRentStarted(ctx context.Context, e domain.RentStartedEvent) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

	if err := s.repo.SaveEvent(ctx, &domain.AnalyticsEvent{
		ID:           uuid.NewString(),
		Subject:      domain.SubjectRentStarted,
		EntityID:     e.RentID,
		Payload:      payload,
		OccurratedAt: e.StartedAt,
		CreatedAt:    time.Now().UTC(),
	}); err != nil {
		return err
	}

	return s.repo.UpsertDailyMetric(ctx, &domain.DailyMetric{
		Date:         truncateToDay(e.StartedAt),
		RentsStarted: 1,
	})
}

func (s *AanalyticsService) HandleRentEnded(ctx context.Context, e domain.RentEndedEvent) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

	if err := s.repo.SaveEvent(ctx, &domain.AnalyticsEvent{
		ID:           uuid.NewString(),
		Subject:      domain.SubjectRentEnded,
		EntityID:     e.RentID,
		Payload:      payload,
		OccurratedAt: e.EndedAt,
		CreatedAt:    time.Now().UTC(),
	}); err != nil {
		return err
	}

	if err := s.repo.UpsertUserStats(ctx, &domain.UserStats{
		UserID:         e.UserID,
		TotalRents:     1,
		TotalSpent:     e.TotalPrice,
		LastActivityAt: e.EndedAt,
		UpdatedAt:      time.Now().UTC(),
	}); err != nil {
		return err
	}

	if err := s.repo.UpsertCarStats(ctx, &domain.CarStats{
		CarID:        e.CarID,
		TotalRents:   1,
		TotalRevenue: e.TotalPrice,
		TotalMinutes: int64(e.DurationMinutes),
		UpdatedAt:    time.Now().UTC(),
	}); err != nil {
		return err
	}

	return s.repo.UpsertDailyMetric(ctx, &domain.DailyMetric{
		Date:           truncateToDay(e.EndedAt),
		RentsComplited: 1,
		Revenue:        e.TotalPrice,
	})
}

// --- helpers ----------------

func truncateToDay(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
