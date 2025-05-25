package repository

import (
	"database/sql"
	"fmt"
	"map-api/internal/models"
	"strings"
	"time"
)

type BlockAreaRepository struct {
	db *sql.DB
}

func NewBlockAreaRepository(db *sql.DB) *BlockAreaRepository {
	return &BlockAreaRepository{db: db}
}

func (r *BlockAreaRepository) Create(area *models.CreateBlockAreaRequest) (*models.BlockArea, error) {
	query := `
        INSERT INTO block_areas (user_id, name, radius, latitude, longitude, altitude, state, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at`

	var id int
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(query, area.UserID, area.Name, area.Radius,
		area.Latitude, area.Longitude, area.Altitude, area.State, area.ExpiresAt).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания запретной зоны: %v", err)
	}

	return &models.BlockArea{
		ID:        id,
		UserID:    area.UserID,
		Name:      area.Name,
		Radius:    area.Radius,
		Latitude:  area.Latitude,
		Longitude: area.Longitude,
		Altitude:  area.Altitude,
		State:     area.State,
		ExpiresAt: area.ExpiresAt,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (r *BlockAreaRepository) GetAll() ([]models.BlockArea, error) {
	query := `
        SELECT id, user_id, name, radius, latitude, longitude, altitude, state, expires_at, created_at, updated_at
        FROM block_areas
        ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения запретных зон: %v", err)
	}
	defer rows.Close()

	var areas []models.BlockArea
	for rows.Next() {
		var area models.BlockArea
		err := rows.Scan(&area.ID, &area.UserID, &area.Name, &area.Radius,
			&area.Latitude, &area.Longitude, &area.Altitude, &area.State, &area.ExpiresAt, &area.CreatedAt, &area.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования запретной зоны: %v", err)
		}
		areas = append(areas, area)
	}

	return areas, nil
}

func (r *BlockAreaRepository) Update(req *models.UpdateBlockAreaRequest) (*models.BlockArea, error) {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Radius != nil {
		setParts = append(setParts, fmt.Sprintf("radius = $%d", argIndex))
		args = append(args, *req.Radius)
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
	if req.ExpiresAt != nil {
		setParts = append(setParts, fmt.Sprintf("expires_at = $%d", argIndex))
		args = append(args, req.ExpiresAt)
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
        UPDATE block_areas 
        SET %s 
        WHERE id = $%d
        RETURNING id, user_id, name, radius, latitude, longitude, altitude, state, expires_at, created_at, updated_at`,
		strings.Join(setParts, ", "), argIndex)

	var area models.BlockArea
	err := r.db.QueryRow(query, args...).Scan(
		&area.ID, &area.UserID, &area.Name, &area.Radius,
		&area.Latitude, &area.Longitude, &area.Altitude, &area.State, &area.ExpiresAt, &area.CreatedAt, &area.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("запретная зона с ID %d не найдена", req.ID)
		}
		return nil, fmt.Errorf("ошибка обновления запретной зоны: %v", err)
	}

	return &area, nil
}

func (r *BlockAreaRepository) CreateActivityLog(userID int, logMessage string) error {
	query := `INSERT INTO activitylogs (user_id, logs) VALUES ($1, $2)`

	_, err := r.db.Exec(query, userID, logMessage)
	if err != nil {
		return fmt.Errorf("ошибка создания лога активности: %v", err)
	}

	return nil
}
