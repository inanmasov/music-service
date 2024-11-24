package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/inanmasov/music-service/internal/handlers"
	"github.com/inanmasov/music-service/internal/logger"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/zhashkevych/todo-app/docs"
)

// @title Music Service API
// @version 1.0
// @description This is a service to manage songs in a library.
// @host localhost:8080
// @BasePath /
func main() {
	log := logger.GetLogger()

	err := godotenv.Load()
	if err != nil {
		log.Error("Loading .env file")
	}

	m, err := migrate.New(
		"file://migrations",
		"postgres://"+os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@"+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+"/"+os.Getenv("DB_NAME")+"?sslmode=disable",
	)
	if err != nil {
		log.Errorf("Initializing migrations: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Errorf("Applying migrations: %v", err)
	}

	log.Info("Migrations applied successfully!")

	// Инициализация роутера
	r := gin.Default()
	log.Info("Gin router initialized")

	// Добавляем статические файлы для swagger
	r.Static("/docs", "./docs")

	// Маршруты для работы с песнями
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("/docs/swagger.json"))) // swagger
	r.GET("/songs", handlers.GetSongs)             // Получение списка песен с фильтрацией и пагинацией
	r.GET("/songs/:id/text", handlers.GetSongText) // Получение текста песни с пагинацией по куплетам
	r.POST("/songs", handlers.AddSong)             // Добавление новой песни
	r.PUT("/songs/:id", handlers.UpdateSong)       // Изменение данных песни
	r.DELETE("/songs/:id", handlers.DeleteSong)    // Удаление песни

	// Запуск сервера
	port := os.Getenv("SERVER_PORT")
	log.Info("Starting server on :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Errorf("Failed to start server: %v", err)
	}
	log.Info("Server started successfully")
}
