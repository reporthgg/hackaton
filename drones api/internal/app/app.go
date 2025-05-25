package app

import (
	"fmt"
	"log"
	"net/http"

	"drones-api/internal/config"
	"drones-api/internal/database"
	"drones-api/internal/handlers"
	"drones-api/internal/service"

	"github.com/gorilla/mux"
)

type App struct {
	config       *config.Config
	droneService *service.DroneService
	handlers     *handlers.DroneHandlers
}

func NewApp() *App {
	cfg := config.New()

	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	droneService := service.NewDroneService(db)
	droneHandlers := handlers.NewDroneHandlers(droneService)

	return &App{
		config:       cfg,
		droneService: droneService,
		handlers:     droneHandlers,
	}
}

func (a *App) Run() error {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/drone/create", a.handlers.CreateDrone).Methods("POST")
	api.HandleFunc("/drone/activate", a.handlers.ActivateDrone).Methods("POST")
	api.HandleFunc("/drone/move", a.handlers.MoveDrone).Methods("POST")
	api.HandleFunc("/drones", a.handlers.GetActiveDrones).Methods("GET")
	api.HandleFunc("/drone/getlist", a.handlers.GetUserDrones).Methods("POST")
	api.HandleFunc("/drone/info", a.handlers.GetDroneInfo).Methods("POST")
	api.HandleFunc("/drone/stop", a.handlers.StopDrone).Methods("POST")

	fmt.Printf("Сервер запущен на порту %s\n", a.config.Port)
	return http.ListenAndServe(":"+a.config.Port, r)
}
