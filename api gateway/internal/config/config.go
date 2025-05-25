package config

import "os"

type Config struct {
	DatabaseURL      string
	JWTSecret        string
	DroneServiceURL  string
	PoliceServiceURL string
	MapServiceURL    string
}

func New() *Config {
	return &Config{
		DatabaseURL: getDatabaseURL(),
		JWTSecret:   getEnv("JWT_SECRET", "dfpokgpofdjgiodfjgifdgj"),
		// DroneServiceURL:  getEnv("DRONE_SERVICE_URL", "http://35.192.62.136:8084"),
		DroneServiceURL:  getEnv("DRONE_SERVICE_URL", "http://localhost:8083"),
		PoliceServiceURL: getEnv("POLICE_SERVICE_URL", "http://localhost:8081"),
		MapServiceURL:    getEnv("MAP_SERVICE_URL", "http://localhost:8082"),
	}
}

func getDatabaseURL() string {
	return "postgres://postgres.akodbsqofninasbpqxbx:HgO3lFR752WPNGN@aws-0-us-east-2.pooler.supabase.com:5432/postgres"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
