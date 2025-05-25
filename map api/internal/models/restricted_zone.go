package models

import "time"

type RestrictedZone struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"user_id" db:"user_id"`
	Height        float64   `json:"height" db:"height"`
	Radius        float64   `json:"radius" db:"radius"`
	DurationHours int       `json:"duration_hours" db:"duration_hours"`
	Latitude      float64   `json:"latitude" db:"latitude"`
	Longitude     float64   `json:"longitude" db:"longitude"`
	Altitude      float64   `json:"altitude" db:"altitude"`
	State         string    `json:"state" db:"state"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CreateRestrictedZoneRequest struct {
	UserID        int     `json:"user_id"`
	Height        float64 `json:"height"`
	Radius        float64 `json:"radius"`
	DurationHours int     `json:"duration_hours"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Altitude      float64 `json:"altitude"`
	State         string  `json:"state"`
}

type UpdateRestrictedZoneRequest struct {
	ID            int      `json:"id"`
	UserID        int      `json:"user_id"`
	Height        *float64 `json:"height,omitempty"`
	Radius        *float64 `json:"radius,omitempty"`
	DurationHours *int     `json:"duration_hours,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	Altitude      *float64 `json:"altitude,omitempty"`
	State         *string  `json:"state,omitempty"`
}
