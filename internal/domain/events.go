package domain

import "time"

const (
	SubjectUserRegistered = "auth.user.registered"
	SubjectUserLoggedIn   = "auth.user.login"
	SubjectTokenRefreshed = "auth.user.refreshed"

	SubjectCarAdded    = "cars.car.added"
	SubjectRentStarted = "cars.rent.started"
	SubjectRentEnded   = "cars.rent.ended"
)

// Auth events

type UserRegisteredEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserLoggedInEvent struct {
	UserID    string    `json:"user_id"`
	IP        string    `json:"ip,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	At        time.Time `json:"at"`
}

type TokenRefreshedEvent struct {
	UserID string    `json:"user_id"`
	At     time.Time `json:"at"`
}

// Cars events

type CarAddedEvent struct {
	CarID    string    `json:"car_id"`
	OwnerID  string    `json:"owner_id"`
	Brand    string    `json:"brand"`
	Model    string    `json:"model"`
	AddeddAt time.Time `json:"added_at"`
}

type RentStartedEvent struct {
	RentID    string    `json:"rent_id"`
	CarID     string    `json:"car_id"`
	UserID    string    `json:"user_id"`
	StartedAt time.Time `json:"started_at"`
}

type RentEndedEvent struct {
	RentID          string    `json:"rent_id"`
	CarID           string    `json:"car_id"`
	UserID          string    `json:"user_id"`
	EndedAt         time.Time `json:"ended_at"`
	DurationMinutes int64     `json:"duration_minutes"`
	TotalPrice      float64   `json:"total_price"`
}
