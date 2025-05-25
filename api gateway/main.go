package main

import (
	"log"
	"net/http"
	"os"

	"api-gateway/internal/config"
	"api-gateway/internal/database"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"api-gateway/internal/proxy"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.New()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
	proxyHandler := proxy.NewProxyHandler()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	router.Use(middleware.Logger())

	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.POST("/drone/create", proxyHandler.ProxyToDroneService)
		api.POST("/drone/activate", proxyHandler.ProxyToDroneService)
		api.POST("/drone/move", proxyHandler.ProxyToDroneService)
		api.GET("/drones", proxyHandler.ProxyToDroneService)
		api.POST("/drone/getlist", proxyHandler.ProxyToDroneService)
		api.POST("/drone/info", proxyHandler.ProxyToDroneService)
		api.POST("/drone/stop", proxyHandler.ProxyToDroneService)

		api.POST("/requests/create", proxyHandler.ProxyToPoliceService)
		api.POST("/requests/user", proxyHandler.ProxyToPoliceService)

		policeOnly := api.Group("")
		policeOnly.Use(middleware.PoliceOnlyMiddleware())
		{
			policeOnly.GET("/requests", proxyHandler.ProxyToPoliceService)
			policeOnly.POST("/requests/update", proxyHandler.ProxyToPoliceService)
		}

		api.GET("/map", proxyHandler.ProxyToMapService)
		api.POST("/map/create", proxyHandler.ProxyToMapService)
		api.POST("/map/update", proxyHandler.ProxyToMapService)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
