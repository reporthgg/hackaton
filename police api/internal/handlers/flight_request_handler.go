package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"police-api/internal/models"
	"police-api/internal/repository"
)

type FlightRequestHandler struct {
	repo *repository.FlightRequestRepository
}

func NewFlightRequestHandler(repo *repository.FlightRequestRepository) *FlightRequestHandler {
	return &FlightRequestHandler{repo: repo}
}

func (h *FlightRequestHandler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	var req models.CreateFlightRequestRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Неверный формат JSON")
		return
	}

	// Валидация
	if req.UserID <= 0 {
		h.sendErrorResponse(w, http.StatusBadRequest, "user_id не найден")
		return
	}

	if req.DroneID <= 0 {
		h.sendErrorResponse(w, http.StatusBadRequest, "drone_id не найден")
		return
	}

	if req.Altitude <= 0 {
		h.sendErrorResponse(w, http.StatusBadRequest, "altitude не найден")
		return
	}

	flightRequest, err := h.repo.Create(&req)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusCreated, flightRequest)
}

func (h *FlightRequestHandler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.repo.GetPendingRequests()
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, requests)
}

func (h *FlightRequestHandler) GetUserRequests(w http.ResponseWriter, r *http.Request) {
	var req models.UserRequestsRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Неверный формат JSON")
		return
	}

	if req.UserID <= 0 {
		h.sendErrorResponse(w, http.StatusBadRequest, "user_id не найден")
		return
	}

	requests, err := h.repo.GetUserRequests(req.UserID)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, requests)
}

func (h *FlightRequestHandler) UpdateRequestState(w http.ResponseWriter, r *http.Request) {
	var updateReq models.UpdateFlightRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Неверный формат JSON")
		return
	}

	if updateReq.ID <= 0 {
		h.sendErrorResponse(w, http.StatusBadRequest, "id не найден")
		return
	}

	if updateReq.UserID <= 0 {
		h.sendErrorResponse(w, http.StatusBadRequest, "user_id не найден")
		return
	}

	if updateReq.State == "" {
		h.sendErrorResponse(w, http.StatusBadRequest, "state не может быть пустым")
		return
	}

	if err := h.repo.UpdateState(updateReq.ID, updateReq.State); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	logMessage := fmt.Sprintf("Пользователь ID %d изменил заявку ID %d на статус '%s'", updateReq.UserID, updateReq.ID, updateReq.State)
	if err := h.repo.CreateActivityLog(updateReq.UserID, logMessage); err != nil {
		fmt.Printf("Ошибка создания лога активности: %v\n", err)
	}

	h.sendSuccessResponse(w, http.StatusOK, map[string]interface{}{
		"id":      updateReq.ID,
		"state":   updateReq.State,
		"message": "Статус заявки успешно обновлен",
	})
}

func (h *FlightRequestHandler) sendSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.APIResponse{
		Success: true,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *FlightRequestHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.APIResponse{
		Success: false,
		Error:   errorMessage,
	}

	json.NewEncoder(w).Encode(response)
}
