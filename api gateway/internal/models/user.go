package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Address      string    `json:"address"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"-"`
	RoleID       int       `json:"role_id"`
	RoleName     string    `json:"role_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRole struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required"`
	RoleID   int    `json:"role_id"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func (u *User) GetByEmail(db *sql.DB, email string) error {
	query := `
		SELECT u.id, u.full_name, u.email, u.address, u.phone, u.password_hash, 
		       u.role_id, ur.name, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN user_roles ur ON u.role_id = ur.id
		WHERE u.email = $1`

	return db.QueryRow(query, email).Scan(
		&u.ID, &u.FullName, &u.Email, &u.Address, &u.Phone,
		&u.PasswordHash, &u.RoleID, &u.RoleName, &u.CreatedAt, &u.UpdatedAt,
	)
}

func (u *User) Create(db *sql.DB) error {
	query := `
		INSERT INTO users (full_name, email, address, phone, password_hash, role_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return db.QueryRow(query,
		u.FullName, u.Email, u.Address, u.Phone, u.PasswordHash, u.RoleID,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (u *User) GetByID(db *sql.DB, id int) error {
	query := `
		SELECT u.id, u.full_name, u.email, u.address, u.phone, u.password_hash, 
		       u.role_id, ur.name, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN user_roles ur ON u.role_id = ur.id
		WHERE u.id = $1`

	return db.QueryRow(query, id).Scan(
		&u.ID, &u.FullName, &u.Email, &u.Address, &u.Phone,
		&u.PasswordHash, &u.RoleID, &u.RoleName, &u.CreatedAt, &u.UpdatedAt,
	)
}
