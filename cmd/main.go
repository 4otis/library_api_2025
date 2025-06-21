package main

import (
	"log"

	"github.com/4otis/library_api_2025/internal/handlers"
	"github.com/4otis/library_api_2025/internal/migrations"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// @title Library API
// @version 1.0
// @description test msg
func main() {
	// TODO: инициализация echo
	e := echo.New()

	// TODO: подключение к БД
	dsn := "host=localhost user=postgres password=password dbname=library port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Error. Failed to connect to db.")
	}

	// TODO: создание миграции
	// err = db.AutoMigrate(&models.Book{}, &models.Author{})
	// if err != nil {
	// 	log.Fatal("Error. Failed to migrate db.")
	// }
	err = migrations.RunInitMigrations(db)
	if err != nil {
		log.Fatal("Error. Failed to migrated db.")
	}

	// TODO: запуск сервера
	handlers.SetupRoutes(e, db)

	e.Logger.Fatal(e.Start(":1323"))
}

func SetupTestServer(db *gorm.DB) *echo.Echo {
	e := echo.New()
	handlers.SetupRoutes(e, db)
	return e
}
