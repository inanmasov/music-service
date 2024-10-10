package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/inanmasov/music-service/internal/handlers"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Инициализация роутера
	r := gin.Default()
	log.Println("info: Gin router initialized")

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
	log.Println("info: Starting server on :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("error: Failed to start server: %v", err)
	}
	log.Println("info: Server started successfully")
}
