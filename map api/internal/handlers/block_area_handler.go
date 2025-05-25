package handlers

import (
	"encoding/json"
	"fmt"
	"map-api/internal/models"
	"map-api/internal/repository"
	"net/http"
)

type BlockAreaHandler struct {
	repo *repository.BlockAreaRepository
}

func NewBlockAreaHandler(repo *repository.BlockAreaRepository) *BlockAreaHandler {
	return &BlockAreaHandler{repo: repo}
}

func (h *BlockAreaHandler) CreateBlockArea(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBlockAreaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка декодирования JSON: %v", err), http.StatusBadRequest)
		return
	}

	if req.UserID <= 0 {
		http.Error(w, "user_id обязателен и должен быть больше 0", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name обязательно", http.StatusBadRequest)
		return
	}
	if req.State == "" {
		req.State = "active"
	}

	area, err := h.repo.Create(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(area)
}

func (h *BlockAreaHandler) GetAllBlockAreas(w http.ResponseWriter, r *http.Request) {
	areas, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(areas)
}

func (h *BlockAreaHandler) UpdateBlockArea(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateBlockAreaRequest
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

	area, err := h.repo.Update(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var changes []string
	if req.Name != nil {
		changes = append(changes, fmt.Sprintf("название: %s", *req.Name))
	}
	if req.Radius != nil {
		changes = append(changes, fmt.Sprintf("радиус: %.2f", *req.Radius))
	}
	if req.Latitude != nil {
		changes = append(changes, fmt.Sprintf("широта: %.6f", *req.Latitude))
	}
	if req.Longitude != nil {
		changes = append(changes, fmt.Sprintf("долгота: %.6f", *req.Longitude))
	}
	if req.Altitude != nil {
		changes = append(changes, fmt.Sprintf("высота: %.2f", *req.Altitude))
	}
	if req.State != nil {
		changes = append(changes, fmt.Sprintf("состояние: %s", *req.State))
	}
	if req.ExpiresAt != nil {
		changes = append(changes, fmt.Sprintf("срок действия до: %s", req.ExpiresAt.Format("2006-01-02 15:04:05")))
	}

	logMessage := fmt.Sprintf("Пользователь ID %d изменил запись ID %d. Изменения: %v", req.UserID, req.ID, changes)

	if err := h.repo.CreateActivityLog(req.UserID, logMessage); err != nil {
		fmt.Printf("Ошибка создания лога активности: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(area)
}
