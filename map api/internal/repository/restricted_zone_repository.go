package repository

import (
	"database/sql"
	"fmt"
	"map-api/internal/models"
	"strings"
	"time"
)

type RestrictedZoneRepository struct {
	db *sql.DB
}

func NewRestrictedZoneRepository(db *sql.DB) *RestrictedZoneRepository {
	return &RestrictedZoneRepository{db: db}
}

func (r *RestrictedZoneRepository) Create(zone *models.CreateRestrictedZoneRequest) (*models.RestrictedZone, error) {
	query := `
        INSERT INTO restricted_zones (user_id, height, radius, duration_hours, latitude, longitude, altitude, state)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at`

	var id int
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(query, zone.UserID, zone.Height, zone.Radius, zone.DurationHours,
		zone.Latitude, zone.Longitude, zone.Altitude, zone.State).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания запретной зоны: %v", err)
	}

	return &models.RestrictedZone{
		ID:            id,
		UserID:        zone.UserID,
		Height:        zone.Height,
		Radius:        zone.Radius,
		DurationHours: zone.DurationHours,
		Latitude:      zone.Latitude,
		Longitude:     zone.Longitude,
		Altitude:      zone.Altitude,
		State:         zone.State,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

func (r *RestrictedZoneRepository) GetAll() ([]models.RestrictedZone, error) {
	query := `
        SELECT id, user_id, height, radius, duration_hours, latitude, longitude, altitude, state, created_at, updated_at
        FROM restricted_zones
        ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения запретных зон: %v", err)
	}
	defer rows.Close()

	var zones []models.RestrictedZone
	for rows.Next() {
		var zone models.RestrictedZone
		err := rows.Scan(&zone.ID, &zone.UserID, &zone.Height, &zone.Radius, &zone.DurationHours,
			&zone.Latitude, &zone.Longitude, &zone.Altitude, &zone.State, &zone.CreatedAt, &zone.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования запретной зоны: %v", err)
		}
		zones = append(zones, zone)
	}

	return zones, nil
}

func (r *RestrictedZoneRepository) Update(req *models.UpdateRestrictedZoneRequest) (*models.RestrictedZone, error) {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Height != nil {
		setParts = append(setParts, fmt.Sprintf("height = $%d", argIndex))
		args = append(args, *req.Height)
		argIndex++
	}
	if req.Radius != nil {
		setParts = append(setParts, fmt.Sprintf("radius = $%d", argIndex))
		args = append(args, *req.Radius)
		argIndex++
	}
	if req.DurationHours != nil {
		setParts = append(setParts, fmt.Sprintf("duration_hours = $%d", argIndex))
		args = append(args, *req.DurationHours)
		argIndex++
	}
	if req.Latitude != nil {
		setParts = append(setParts, fmt.Sprintf("latitude = $%d", argIndex))
		args = append(args, *req.Latitude)
		argIndex++
	}
	if req.Longitude != nil {
		setParts = append(setParts, fmt.Sprintf("longitude = $%d", argIndex))
		args = append(args, *req.Longitude)
		argIndex++
	}
	if req.Altitude != nil {
		setParts = append(setParts, fmt.Sprintf("altitude = $%d", argIndex))
		args = append(args, *req.Altitude)
		argIndex++
	}
	if req.State != nil {
		setParts = append(setParts, fmt.Sprintf("state = $%d", argIndex))
		args = append(args, *req.State)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("нет полей для обновления")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, req.ID)

	query := fmt.Sprintf(`
        UPDATE restricted_zones 
        SET %s 
        WHERE id = $%d
        RETURNING id, user_id, height, radius, duration_hours, latitude, longitude, altitude, state, created_at, updated_at`,
		strings.Join(setParts, ", "), argIndex)

	var zone models.RestrictedZone
	err := r.db.QueryRow(query, args...).Scan(
		&zone.ID, &zone.UserID, &zone.Height, &zone.Radius, &zone.DurationHours,
		&zone.Latitude, &zone.Longitude, &zone.Altitude, &zone.State, &zone.CreatedAt, &zone.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("запретная зона с ID %d не найдена", req.ID)
		}
		return nil, fmt.Errorf("ошибка обновления запретной зоны: %v", err)
	}

	return &zone, nil
}

func (r *RestrictedZoneRepository) CreateActivityLog(userID int, logMessage string) error {
	query := `INSERT INTO activitylogs (user_id, logs) VALUES ($1, $2)`

	_, err := r.db.Exec(query, userID, logMessage)
	if err != nil {
		return fmt.Errorf("ошибка создания лога активности: %v", err)
	}

	return nil
}
