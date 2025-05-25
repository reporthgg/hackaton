package handlers

import (
	"api-gateway/internal/models"
	"api-gateway/internal/utils"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db        *sql.DB
	jwtSecret string
}

func NewAuthHandler(db *sql.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пароль должен содержать минимум 6 символов"})
		return
	}

	var existingUser models.User
	if err := existingUser.GetByEmail(h.db, req.Email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким email уже существует"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании пароля"})
		return
	}

	if req.RoleID == 0 {
		req.RoleID = 1
	}

	user := models.User{
		FullName:     req.FullName,
		Email:        req.Email,
		Address:      req.Address,
		Phone:        req.Phone,
		PasswordHash: string(hashedPassword),
		RoleID:       req.RoleID,
	}

	if err := user.Create(h.db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}

	if err := user.GetByID(h.db, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных пользователя"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.RoleID, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	user.PasswordHash = ""

	c.JSON(http.StatusCreated, models.LoginResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := user.GetByEmail(h.db, req.Email); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.RoleID, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	user.PasswordHash = ""

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  user,
	})
}
