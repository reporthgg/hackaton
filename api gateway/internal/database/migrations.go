package database

import "database/sql"

func Migrate(db *sql.DB) error {
	createUserRolesTable := `
	CREATE TABLE IF NOT EXISTS user_roles (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) UNIQUE NOT NULL,
		description TEXT
	);`

	if _, err := db.Exec(createUserRolesTable); err != nil {
		return err
	}

	insertRoles := `
	INSERT INTO user_roles (name, description) 
	VALUES 
		('user', 'Обычный пользователь'),
		('police', 'Сотрудник полиции')
	ON CONFLICT (name) DO NOTHING;`

	if _, err := db.Exec(insertRoles); err != nil {
		return err
	}

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		full_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		address TEXT,
		phone VARCHAR(20),
		password_hash VARCHAR(255) NOT NULL,
		role_id INTEGER REFERENCES user_roles(id) DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createUsersTable); err != nil {
		return err
	}

	return nil
}
