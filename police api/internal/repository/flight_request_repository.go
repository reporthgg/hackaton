package repository

import (
	"fmt"
	"police-api/internal/database"
	"police-api/internal/models"
)

type FlightRequestRepository struct {
	db *database.DB
}

func NewFlightRequestRepository(db *database.DB) *FlightRequestRepository {
	return &FlightRequestRepository{db: db}
}

func (r *FlightRequestRepository) Create(req *models.CreateFlightRequestRequest) (*models.FlightRequest, error) {
	query := `
        INSERT INTO flightrequest (user_id, drone_id, departure_time, altitude, start_lat, start_lng, end_lat, end_lng, state)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending')
        RETURNING id, user_id, drone_id, departure_time, altitude, start_lat, start_lng, end_lat, end_lng, state, created_at`

	var flightRequest models.FlightRequest
	err := r.db.QueryRow(query,
		req.UserID,
		req.DroneID,
		req.DepartureTime,
		req.Altitude,
		req.StartLat,
		req.StartLng,
		req.EndLat,
		req.EndLng,
	).Scan(
		&flightRequest.ID,
		&flightRequest.UserID,
		&flightRequest.DroneID,
		&flightRequest.DepartureTime,
		&flightRequest.Altitude,
		&flightRequest.StartLat,
		&flightRequest.StartLng,
		&flightRequest.EndLat,
		&flightRequest.EndLng,
		&flightRequest.State,
		&flightRequest.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания заявки: %v", err)
	}

	usernameQuery := `SELECT full_name FROM users WHERE id = $1`
	err = r.db.QueryRow(usernameQuery, flightRequest.UserID).Scan(&flightRequest.Username)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения имени пользователя: %v", err)
	}

	return &flightRequest, nil
}

func (r *FlightRequestRepository) GetPendingRequests() ([]models.FlightRequest, error) {
	query := `
        SELECT fr.id, fr.user_id, u.full_name, fr.drone_id, fr.departure_time, fr.altitude, 
               fr.start_lat, fr.start_lng, fr.end_lat, fr.end_lng, fr.state, fr.created_at
        FROM flightrequest fr
        JOIN users u ON fr.user_id = u.id
        WHERE fr.state = 'pending'
        ORDER BY fr.created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заявок: %v", err)
	}
	defer rows.Close()

	var requests []models.FlightRequest
	for rows.Next() {
		var req models.FlightRequest
		err := rows.Scan(
			&req.ID,
			&req.UserID,
			&req.Username,
			&req.DroneID,
			&req.DepartureTime,
			&req.Altitude,
			&req.StartLat,
			&req.StartLng,
			&req.EndLat,
			&req.EndLng,
			&req.State,
			&req.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования заявки: %v", err)
		}
		requests = append(requests, req)
	}

	return requests, nil
}

func (r *FlightRequestRepository) GetUserRequests(userID int) ([]models.FlightRequest, error) {
	query := `
        SELECT id, user_id, drone_id, departure_time, altitude, start_lat, start_lng, end_lat, end_lng, state, created_at
        FROM flightrequest 
        WHERE user_id = $1
        ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заявок пользователя: %v", err)
	}
	defer rows.Close()

	var requests []models.FlightRequest
	for rows.Next() {
		var req models.FlightRequest
		err := rows.Scan(
			&req.ID,
			&req.UserID,
			&req.DroneID,
			&req.DepartureTime,
			&req.Altitude,
			&req.StartLat,
			&req.StartLng,
			&req.EndLat,
			&req.EndLng,
			&req.State,
			&req.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования заявки пользователя: %v", err)
		}
		requests = append(requests, req)
	}

	return requests, nil
}

func (r *FlightRequestRepository) UpdateState(id int, state string) error {
	if state != "approved" && state != "denied" {
		return fmt.Errorf("неверный статус: %s. Разрешены только 'approved' или 'denied'", state)
	}

	query := `UPDATE flightrequest SET state = $1 WHERE id = $2 AND state = 'pending'`
	result, err := r.db.Exec(query, state, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления статуса заявки: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("заявка с ID %d не найдена или уже обработана", id)
	}

	return nil
}

func (r *FlightRequestRepository) CreateActivityLog(userID int, logMessage string) error {
	query := `INSERT INTO activitylogs (user_id, logs) VALUES ($1, $2)`

	_, err := r.db.Exec(query, userID, logMessage)
	if err != nil {
		return fmt.Errorf("ошибка создания лога активности: %v", err)
	}

	return nil
}
