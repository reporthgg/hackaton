package main

import (
	"log"
	"net/http"

	"map-api/internal/config"
	"map-api/internal/handlers"
	"map-api/internal/repository"

	"github.com/gorilla/mux"
)

func main() {
	db, err := config.NewDatabaseConnection()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer db.Close()

	blockAreaRepo := repository.NewBlockAreaRepository(db)
	blockAreaHandler := handlers.NewBlockAreaHandler(blockAreaRepo)

	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/map/create", blockAreaHandler.CreateBlockArea).Methods("POST")
	api.HandleFunc("/map/update", blockAreaHandler.UpdateBlockArea).Methods("POST")
	api.HandleFunc("/map", blockAreaHandler.GetAllBlockAreas).Methods("GET")

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	log.Println("Сервер запущен на порту 8082...")
	log.Fatal(http.ListenAndServe(":8082", router))
}
