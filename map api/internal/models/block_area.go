package models

import "time"

type BlockArea struct {
	ID        int        `json:"id" db:"id"`
	UserID    int        `json:"user_id" db:"user_id"`
	Name      string     `json:"name" db:"name"`
	Radius    float64    `json:"radius" db:"radius"`
	Latitude  float64    `json:"latitude" db:"latitude"`
	Longitude float64    `json:"longitude" db:"longitude"`
	Altitude  float64    `json:"altitude" db:"altitude"`
	State     string     `json:"state" db:"state"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateBlockAreaRequest struct {
	UserID    int        `json:"user_id"`
	Name      string     `json:"name"`
	Radius    float64    `json:"radius"`
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
	Altitude  float64    `json:"altitude"`
	State     string     `json:"state"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type UpdateBlockAreaRequest struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Name      *string    `json:"name,omitempty"`
	Radius    *float64   `json:"radius,omitempty"`
	Latitude  *float64   `json:"latitude,omitempty"`
	Longitude *float64   `json:"longitude,omitempty"`
	Altitude  *float64   `json:"altitude,omitempty"`
	State     *string    `json:"state,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}
