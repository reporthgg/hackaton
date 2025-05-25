package main

import (
	"log"
	"net/http"

	"police-api/internal/database"
	"police-api/internal/handlers"
	"police-api/internal/repository"

	"github.com/gorilla/mux"
)

func main() {
	// Подключение к базе данных
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer db.Close()

	// Создание таблиц
	if err := db.CreateTables(); err != nil {
		log.Fatal("Ошибка создания таблиц:", err)
	}

	// Инициализация репозитория и обработчиков
	flightRequestRepo := repository.NewFlightRequestRepository(db)
	flightRequestHandler := handlers.NewFlightRequestHandler(flightRequestRepo)

	// Настройка маршрутов
	router := mux.NewRouter()

	// API маршруты
	api := router.PathPrefix("/api").Subrouter()
	api.Use(jsonMiddleware)

	// Маршруты для заявок
	api.HandleFunc("/requests/create", flightRequestHandler.CreateRequest).Methods("POST")
	api.HandleFunc("/requests", flightRequestHandler.GetPendingRequests).Methods("GET")
	api.HandleFunc("/requests/user", flightRequestHandler.GetUserRequests).Methods("POST")
	api.HandleFunc("/requests/update", flightRequestHandler.UpdateRequestState).Methods("POST")

	// Запуск сервера
	log.Println("Сервер запущен на порту :8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
