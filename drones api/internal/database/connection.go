package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewConnection(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения с БД: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка проверки соединения с БД: %w", err)
	}

	fmt.Println("Подключение к базе данных успешно установлено")
	return db, nil
}
