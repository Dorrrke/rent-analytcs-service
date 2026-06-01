package domain

import "time"

type AnalyticsEvent struct {
	ID           string    `db:"id"`
	Subject      string    `db:"subject"`
	EntityID     string    `db:"entity_id"`
	Payload      []byte    `db:"payload"`
	OccurratedAt time.Time `db:"occurred_at"`
	CreatedAt    time.Time `db:"created_at"`
}

type UserStats struct {
	UserID         string    `db:"user_id"`
	TotalLogins    int64     `db:"total_logins"`
	TotalRents     int64     `db:"total_rents"`
	TotalSpent     float64   `db:"total_spent"`
	LastActivityAt time.Time `db:"last_activity_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type CarStats struct {
	CarID        string    `db:"car_id"`
	TotalRents   int64     `db:"total_rents"`
	TotalRevenue float64   `db:"total_revenue"`
	TotalMinutes int64     `db:"total_minutes"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type DailyMetric struct {
	Date           time.Time `db:"date"`
	NewUsers       int64     `db:"new_users"`
	Logins         int64     `db:"logins"`
	RentsStarted   int64     `db:"rents_started"`
	RentsComplited int64     `db:"rents_completed"`
	Revenue        float64   `db:"revenue"`
}
