package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

func New() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=aws-0-us-east-2.pooler.supabase.com port=5432 user=postgres.akodbsqofninasbpqxbx password=HgO3lFR752WPNGN dbname=postgres sslmode=require"
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
	}
}
