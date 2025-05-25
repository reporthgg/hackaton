package service

import (
	"database/sql"
	"fmt"
	"math"
	"sync"
	"time"

	"drones-api/internal/models"
)

type DroneService struct {
	db           *sql.DB
	movingDrones map[int]chan bool
	mu           sync.RWMutex
}

func NewDroneService(db *sql.DB) *DroneService {
	return &DroneService{
		db:           db,
		movingDrones: make(map[int]chan bool),
	}
}

func (ds *DroneService) CreateDrone(req models.CreateDroneRequest) (*models.Drone, error) {
	query := `
		INSERT INTO drones (name, owner_id, max_speed, current_status, battery_level, created_at, updated_at)
		VALUES ($1, $2, $3, 'offline', 100, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id`

	var droneID int
	err := ds.db.QueryRow(query, req.Name, req.UserID, req.MaxSpeed).Scan(&droneID)
	if err != nil {
		fmt.Printf("Ошибка вставки дрона: %v\n", err)
		return nil, err
	}

	selectQuery := `
		SELECT id, name, owner_id, current_lat, current_lng, current_altitude,
		       current_status, battery_level, max_speed, created_at, updated_at
		FROM drones WHERE id = $1`

	var drone models.Drone
	err = ds.db.QueryRow(selectQuery, droneID).Scan(
		&drone.ID, &drone.Name, &drone.OwnerID, &drone.CurrentLat,
		&drone.CurrentLng, &drone.CurrentAltitude, &drone.CurrentStatus,
		&drone.BatteryLevel, &drone.MaxSpeed, &drone.CreatedAt, &drone.UpdatedAt)

	if err != nil {
		fmt.Printf("Ошибка получения созданного дрона: %v\n", err)
		return nil, err
	}

	return &drone, nil
}

func (ds *DroneService) ActivateDrone(req models.ActivateDroneRequest) error {
	var ownerID int
	err := ds.db.QueryRow("SELECT owner_id FROM drones WHERE id = $1", req.DroneID).Scan(&ownerID)
	if err != nil {
		return err
	}

	if ownerID != req.UserID {
		return ErrAccessDenied
	}

	query := `
		UPDATE drones 
		SET current_lat = $1, current_lng = $2, current_altitude = $3, 
		    current_status = 'active', battery_level = 100, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4`

	_, err = ds.db.Exec(query, req.Lat, req.Lng, req.Altitude, req.DroneID)
	return err
}

func (ds *DroneService) MoveDrone(req models.MoveDroneRequest) error {
	var drone models.Drone
	query := `SELECT id, owner_id, current_lat, current_lng, current_altitude, current_status, battery_level 
	          FROM drones WHERE id = $1`

	err := ds.db.QueryRow(query, req.DroneID).Scan(
		&drone.ID, &drone.OwnerID, &drone.CurrentLat, &drone.CurrentLng,
		&drone.CurrentAltitude, &drone.CurrentStatus, &drone.BatteryLevel)

	if err != nil {
		return err
	}

	if drone.OwnerID != req.UserID {
		return ErrAccessDenied
	}

	if drone.CurrentLat == nil || drone.CurrentLng == nil || drone.CurrentAltitude == nil {
		return ErrDroneNotActivated
	}

	ds.stopDroneMovement(req.DroneID)

	go ds.simulateMovement(req.DroneID, *drone.CurrentLat, *drone.CurrentLng, *drone.CurrentAltitude,
		req.TargetLat, req.TargetLng, req.TargetAltitude, req.BatteryLevel, req.Speed)

	return nil
}

func (ds *DroneService) GetActiveDrones() ([]models.Drone, error) {
	query := `SELECT id, name, owner_id, current_lat, current_lng, current_altitude, 
	          current_status, battery_level, max_speed, created_at, updated_at 
	          FROM drones WHERE current_status != 'offline'`

	rows, err := ds.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drones []models.Drone
	for rows.Next() {
		var drone models.Drone
		err := rows.Scan(&drone.ID, &drone.Name, &drone.OwnerID, &drone.CurrentLat,
			&drone.CurrentLng, &drone.CurrentAltitude, &drone.CurrentStatus,
			&drone.BatteryLevel, &drone.MaxSpeed, &drone.CreatedAt, &drone.UpdatedAt)
		if err != nil {
			continue
		}
		drones = append(drones, drone)
	}

	return drones, nil
}

func (ds *DroneService) GetDroneInfo(req models.DroneInfoRequest) (*models.Drone, error) {
	query := `SELECT id, name, owner_id, current_lat, current_lng, current_altitude, 
	          current_status, battery_level, max_speed, created_at, updated_at 
	          FROM drones WHERE id = $1`

	var drone models.Drone
	err := ds.db.QueryRow(query, req.DroneID).Scan(
		&drone.ID, &drone.Name, &drone.OwnerID, &drone.CurrentLat,
		&drone.CurrentLng, &drone.CurrentAltitude, &drone.CurrentStatus,
		&drone.BatteryLevel, &drone.MaxSpeed, &drone.CreatedAt, &drone.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if drone.OwnerID != req.UserID {
		return nil, ErrAccessDenied
	}

	return &drone, nil
}

func (ds *DroneService) StopDrone(req models.StopDroneRequest) error {
	var ownerID int
	err := ds.db.QueryRow("SELECT owner_id FROM drones WHERE id = $1", req.DroneID).Scan(&ownerID)
	if err != nil {
		return err
	}

	if ownerID != req.UserID && req.UserRole != 2 {
		return ErrAccessDenied
	}

	ds.stopDroneMovement(req.DroneID)
	return nil
}

func (ds *DroneService) simulateMovement(droneID int, startLat, startLng, startAlt, targetLat, targetLng, targetAlt float64, batteryLevel int, speed float64) {
	ds.mu.Lock()
	stopChan := make(chan bool)
	ds.movingDrones[droneID] = stopChan
	ds.mu.Unlock()

	ds.db.Exec("UPDATE drones SET current_status = 'flying', updated_at = CURRENT_TIMESTAMP WHERE id = $1", droneID)

	currentLat, currentLng, currentAlt := startLat, startLng, startAlt
	currentBattery := batteryLevel

	totalDistance := ds.calculateDistance(startLat, startLng, targetLat, targetLng)
	altitudeDiff := targetAlt - startAlt

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			ds.db.Exec(`UPDATE drones SET current_status = 'stopped', battery_level = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
				currentBattery, droneID)
			ds.mu.Lock()
			delete(ds.movingDrones, droneID)
			ds.mu.Unlock()
			return

		case <-ticker.C:
			currentBattery = int(float64(currentBattery) - 0.02)
			if currentBattery <= 0 {
				currentBattery = 0
				ds.db.Exec(`UPDATE drones SET current_status = 'nullbattery', battery_level = 0, updated_at = CURRENT_TIMESTAMP WHERE id = $1`, droneID)
				ds.mu.Lock()
				delete(ds.movingDrones, droneID)
				ds.mu.Unlock()
				return
			}

			distancePerSecond := speed / 3.6

			distanceToTarget := ds.calculateDistance(currentLat, currentLng, targetLat, targetLng)

			if distanceToTarget <= distancePerSecond {
				// Достигли цели
				currentLat, currentLng, currentAlt = targetLat, targetLng, targetAlt
				ds.db.Exec(`UPDATE drones SET current_lat = $1, current_lng = $2, current_altitude = $3, 
				            current_status = 'active', battery_level = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`,
					currentLat, currentLng, currentAlt, currentBattery, droneID)
				ds.mu.Lock()
				delete(ds.movingDrones, droneID)
				ds.mu.Unlock()
				return
			}

			ratio := distancePerSecond / distanceToTarget
			latDiff := (targetLat - currentLat) * ratio
			lngDiff := (targetLng - currentLng) * ratio
			altDiff := altitudeDiff * (distancePerSecond / totalDistance)

			currentLat += latDiff
			currentLng += lngDiff
			currentAlt += altDiff

			ds.db.Exec(`UPDATE drones SET current_lat = $1, current_lng = $2, current_altitude = $3, 
			            battery_level = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`,
				currentLat, currentLng, currentAlt, currentBattery, droneID)
		}
	}
}

func (ds *DroneService) calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371000

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func (ds *DroneService) stopDroneMovement(droneID int) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if stopChan, exists := ds.movingDrones[droneID]; exists {
		close(stopChan)
		delete(ds.movingDrones, droneID)
	}
}
