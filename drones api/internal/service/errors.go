package service

import "errors"

var (
	ErrAccessDenied      = errors.New("доступ запрещен")
	ErrDroneNotFound     = errors.New("дрон не найден")
	ErrDroneNotActivated = errors.New("дрон не активирован")
)
