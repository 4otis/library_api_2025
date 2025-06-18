package main

import (
	"log"

	"github.com/4otis/library_api_2025/internal/handlers"
	"github.com/4otis/library_api_2025/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Library API
// @version 1.0
// @description test msg
func main() {
	// TODO: инициализация echo
	e := echo.New()

	// TODO: подключение к БД
	dsn := "host=localhost user=postgres password=password dbname=library port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error. Failed to connect to db.")
	}

	// TODO: создание миграции
	err = db.AutoMigrate(&models.Book{}, &models.Author{})
	if err != nil {
		log.Fatal("Error. Failed to migrate db.")
	}

	// TODO: запуск сервера
	handlers.SetupRoutes(e, db)

	e.Logger.Fatal(e.Start(":1323"))
}
