package main

import (
	"fmt"
	"log"
	"pvz_system/internal/app"

	"github.com/joho/godotenv"
)

func LoadEnv() error {

	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("ошибка чтения конфигурации: %w", err)
	}
	return nil
}

func main() {
	err := LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("ошибка инициализации приложения: %v", err)
	}

	application.Run()

}
