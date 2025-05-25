package main

import (
	"log"

	"drones-api/internal/app"
)

func main() {
	application := app.NewApp()
	if err := application.Run(); err != nil {
		log.Fatal("Ошибка запуска приложения:", err)
	}
}
