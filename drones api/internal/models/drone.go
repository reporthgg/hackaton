package models

import "time"

type Drone struct {
	ID              int        `json:"id" db:"id"`
	Name            string     `json:"name" db:"name"`
	OwnerID         int        `json:"owner_id" db:"owner_id"`
	CurrentLat      *float64   `json:"current_lat" db:"current_lat"`
	CurrentLng      *float64   `json:"current_lng" db:"current_lng"`
	CurrentAltitude *float64   `json:"current_altitude" db:"current_altitude"`
	CurrentStatus   string     `json:"current_status" db:"current_status"`
	BatteryLevel    int        `json:"battery_level" db:"battery_level"`
	MaxSpeed        float64    `json:"max_speed" db:"max_speed"`
	CreatedAt       *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at" db:"updated_at"`
}

type CreateDroneRequest struct {
	UserID   int     `json:"user_id"`
	Name     string  `json:"name"`
	MaxSpeed float64 `json:"max_speed"`
}

type ActivateDroneRequest struct {
	UserID   int     `json:"user_id"`
	DroneID  int     `json:"drone_id"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Altitude float64 `json:"altitude"`
}

type MoveDroneRequest struct {
	UserID         int     `json:"user_id"`
	DroneID        int     `json:"drone_id"`
	TargetLat      float64 `json:"target_lat"`
	TargetLng      float64 `json:"target_lng"`
	TargetAltitude float64 `json:"target_altitude"`
	BatteryLevel   int     `json:"battery_level"`
	Speed          float64 `json:"speed"`
}

type DroneInfoRequest struct {
	UserID  int `json:"user_id"`
	DroneID int `json:"drone_id"`
}

type StopDroneRequest struct {
	UserID   int `json:"user_id"`
	DroneID  int `json:"drone_id"`
	UserRole int `json:"user_role"`
}

type GetUserDronesRequest struct {
	UserID int `json:"user_id"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
