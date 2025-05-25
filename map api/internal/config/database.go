package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func NewDatabaseConnection() (*sql.DB, error) {
	config := DatabaseConfig{
		Host:     "aws-0-us-east-2.pooler.supabase.com",
		Port:     "5432",
		Database: "postgres",
		User:     "postgres.akodbsqofninasbpqxbx",
		Password: "HgO3lFR752WPNGN",
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		config.Host, config.Port, config.User, config.Password, config.Database)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка проверки подключения к базе данных: %v", err)
	}

	return db, nil
}
