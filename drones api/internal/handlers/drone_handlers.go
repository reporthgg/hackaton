package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"drones-api/internal/models"
	"drones-api/internal/service"
)

type DroneHandlers struct {
	droneService *service.DroneService
}

func NewDroneHandlers(droneService *service.DroneService) *DroneHandlers {
	return &DroneHandlers{
		droneService: droneService,
	}
}

func (h *DroneHandlers) CreateDrone(w http.ResponseWriter, r *http.Request) {
	var req models.CreateDroneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendResponse(w, false, "Неверный формат JSON", nil, http.StatusBadRequest)
		return
	}

	drone, err := h.droneService.CreateDrone(req)
	if err != nil {
		fmt.Printf("Ошибка создания дрона: %v\n", err)
		h.sendResponse(w, false, "Ошибка создания дрона", nil, http.StatusInternalServerError)
		return
	}

	h.sendResponse(w, true, "Дрон успешно создан", drone, http.StatusCreated)
}

func (h *DroneHandlers) ActivateDrone(w http.ResponseWriter, r *http.Request) {
	var req models.ActivateDroneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendResponse(w, false, "Неверный формат JSON", nil, http.StatusBadRequest)
		return
	}

	err := h.droneService.ActivateDrone(req)
	if err != nil {
		if err == service.ErrAccessDenied {
			h.sendResponse(w, false, "Доступ запрещен", nil, http.StatusForbidden)
			return
		}
		h.sendResponse(w, false, "Ошибка активации дрона", nil, http.StatusInternalServerError)
		return
	}

	h.sendResponse(w, true, "Дрон успешно активирован", nil, http.StatusOK)
}

func (h *DroneHandlers) MoveDrone(w http.ResponseWriter, r *http.Request) {
	var req models.MoveDroneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendResponse(w, false, "Неверный формат JSON", nil, http.StatusBadRequest)
		return
	}

	err := h.droneService.MoveDrone(req)
	if err != nil {
		switch err {
		case service.ErrAccessDenied:
			h.sendResponse(w, false, "Доступ запрещен", nil, http.StatusForbidden)
		case service.ErrDroneNotActivated:
			h.sendResponse(w, false, "Дрон не активирован", nil, http.StatusBadRequest)
		default:
			h.sendResponse(w, false, "Ошибка запуска движения дрона", nil, http.StatusInternalServerError)
		}
		return
	}

	h.sendResponse(w, true, "Дрон начал движение", nil, http.StatusOK)
}

func (h *DroneHandlers) GetActiveDrones(w http.ResponseWriter, r *http.Request) {
	drones, err := h.droneService.GetActiveDrones()
	if err != nil {
		h.sendResponse(w, false, "Ошибка получения дронов", nil, http.StatusInternalServerError)
		return
	}

	h.sendResponse(w, true, "Активные дроны получены", drones, http.StatusOK)
}

func (h *DroneHandlers) GetUserDrones(w http.ResponseWriter, r *http.Request) {
	var req models.GetUserDronesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendResponse(w, false, "Неверный формат JSON", nil, http.StatusBadRequest)
		return
	}

	drones, err := h.droneService.GetUserDrones(req)
	if err != nil {
		h.sendResponse(w, false, "Ошибка получения дронов пользователя", nil, http.StatusInternalServerError)
		return
	}

	h.sendResponse(w, true, "Дроны пользователя получены", drones, http.StatusOK)
}

func (h *DroneHandlers) GetDroneInfo(w http.ResponseWriter, r *http.Request) {
	var req models.DroneInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendResponse(w, false, "Неверный формат JSON", nil, http.StatusBadRequest)
		return
	}

	drone, err := h.droneService.GetDroneInfo(req)
	if err != nil {
		if err == service.ErrAccessDenied {
			h.sendResponse(w, false, "Доступ запрещен", nil, http.StatusForbidden)
			return
		}
		h.sendResponse(w, false, "Дрон не найден", nil, http.StatusNotFound)
		return
	}

	h.sendResponse(w, true, "Информация о дроне получена", drone, http.StatusOK)
}

func (h *DroneHandlers) StopDrone(w http.ResponseWriter, r *http.Request) {
	var req models.StopDroneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendResponse(w, false, "Неверный формат JSON", nil, http.StatusBadRequest)
		return
	}

	err := h.droneService.StopDrone(req)
	if err != nil {
		if err == service.ErrAccessDenied {
			h.sendResponse(w, false, "Доступ запрещен", nil, http.StatusForbidden)
			return
		}
		h.sendResponse(w, false, "Дрон не найден", nil, http.StatusNotFound)
		return
	}

	h.sendResponse(w, true, "Дрон остановлен", nil, http.StatusOK)
}

func (h *DroneHandlers) sendResponse(w http.ResponseWriter, success bool, message string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.APIResponse{
		Success: success,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}
