package handlers

import (
	"encoding/json"
	"fmt"
	"map-api/internal/models"
	"map-api/internal/repository"
	"net/http"
)

type RestrictedZoneHandler struct {
	repo *repository.RestrictedZoneRepository
}

func NewRestrictedZoneHandler(repo *repository.RestrictedZoneRepository) *RestrictedZoneHandler {
	return &RestrictedZoneHandler{repo: repo}
}

func (h *RestrictedZoneHandler) CreateRestrictedZone(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRestrictedZoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка декодирования JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Валидация
	if req.UserID <= 0 {
		http.Error(w, "user_id обязателен и должен быть больше 0", http.StatusBadRequest)
		return
	}
	if req.State == "" {
		req.State = "active"
	}

	zone, err := h.repo.Create(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zone)
}

func (h *RestrictedZoneHandler) GetAllRestrictedZones(w http.ResponseWriter, r *http.Request) {
	zones, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}

func (h *RestrictedZoneHandler) UpdateRestrictedZone(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateRestrictedZoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка декодирования JSON: %v", err), http.StatusBadRequest)
		return
	}

	if req.ID <= 0 {
		http.Error(w, "id обязателен и должен быть больше 0", http.StatusBadRequest)
		return
	}
	if req.UserID <= 0 {
		http.Error(w, "user_id обязателен и должен быть больше 0", http.StatusBadRequest)
		return
	}

	zone, err := h.repo.Update(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var changes []string
	if req.Height != nil {
		changes = append(changes, fmt.Sprintf("высота: %.2f", *req.Height))
	}
	if req.Radius != nil {
		changes = append(changes, fmt.Sprintf("радиус: %.2f", *req.Radius))
	}
	if req.DurationHours != nil {
		changes = append(changes, fmt.Sprintf("время действия: %d часов", *req.DurationHours))
	}
	if req.Latitude != nil {
		changes = append(changes, fmt.Sprintf("широта: %.6f", *req.Latitude))
	}
	if req.Longitude != nil {
		changes = append(changes, fmt.Sprintf("долгота: %.6f", *req.Longitude))
	}
	if req.Altitude != nil {
		changes = append(changes, fmt.Sprintf("высота над уровнем моря: %.2f", *req.Altitude))
	}
	if req.State != nil {
		changes = append(changes, fmt.Sprintf("состояние: %s", *req.State))
	}

	logMessage := fmt.Sprintf("Пользователь ID %d изменил запись ID %d. Изменения: %v", req.UserID, req.ID, changes)

	if err := h.repo.CreateActivityLog(req.UserID, logMessage); err != nil {
		fmt.Printf("Ошибка создания лога активности: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zone)
}
