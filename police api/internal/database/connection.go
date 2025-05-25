package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewConnection() (*DB, error) {
	dbHost := "aws-0-us-east-2.pooler.supabase.com"
	dbPort := 5432
	dbUser := "postgres.akodbsqofninasbpqxbx"
	dbPassword := "HgO3lFR752WPNGN"
	dbName := "postgres"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping базы данных: %v", err)
	}

	log.Println("Успешное подключение к базе данных PostgreSQL")
	return &DB{db}, nil
}

func (db *DB) CreateTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS flightrequest (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		drone_id INTEGER NOT NULL,
		departure_time TIMESTAMP NOT NULL,
		altitude DOUBLE PRECISION NOT NULL,
		start_lat DOUBLE PRECISION NOT NULL,
		start_lng DOUBLE PRECISION NOT NULL,
		end_lat DOUBLE PRECISION NOT NULL,
		end_lng DOUBLE PRECISION NOT NULL,
		state VARCHAR(20) NOT NULL DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы flightrequest: %v", err)
	}

	log.Println("Таблица flightrequest создана успешно")

	activityLogsQuery := `
	CREATE TABLE IF NOT EXISTS activitylogs (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		logs TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(activityLogsQuery)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы activitylogs: %v", err)
	}
	log.Println("Таблица activitylogs создана успешно")

	return nil
}
