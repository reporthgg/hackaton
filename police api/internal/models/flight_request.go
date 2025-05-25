package models

import (
	"time"
)

type FlightRequest struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"user_id" db:"user_id"`
	Username      string    `json:"username" db:"username"`
	DroneID       int       `json:"drone_id" db:"drone_id"`
	DepartureTime time.Time `json:"departure_time" db:"departure_time"`
	Altitude      float64   `json:"altitude" db:"altitude"`
	StartLat      float64   `json:"start_lat" db:"start_lat"`
	StartLng      float64   `json:"start_lng" db:"start_lng"`
	EndLat        float64   `json:"end_lat" db:"end_lat"`
	EndLng        float64   `json:"end_lng" db:"end_lng"`
	State         string    `json:"state" db:"state"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type CreateFlightRequestRequest struct {
	UserID        int       `json:"user_id"`
	DroneID       int       `json:"drone_id"`
	DepartureTime time.Time `json:"departure_time"`
	Altitude      float64   `json:"altitude"`
	StartLat      float64   `json:"start_lat"`
	StartLng      float64   `json:"start_lng"`
	EndLat        float64   `json:"end_lat"`
	EndLng        float64   `json:"end_lng"`
}

type UpdateFlightRequestRequest struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	State  string `json:"state"`
}

type UserRequestsRequest struct {
	UserID int `json:"user_id"`
}

type ActivityLog struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Logs      string    `json:"logs" db:"logs"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
