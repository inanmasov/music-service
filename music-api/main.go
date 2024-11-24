package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SongDetail структура ответа
type SongDetail struct {
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

func main() {
	r := gin.Default()

	// Обработчик для маршрута /info
	r.GET("/info", func(c *gin.Context) {
		group := c.Query("group")
		song := c.Query("song")

		// Проверяем, что обязательные параметры переданы
		if group == "" || song == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Group and song are required parameters"})
			return
		}

		log.Printf("info: Received request for group: %s, song: %s", group, song)

		dateStr := "16.07.2006"
		layout := "02.01.2006" // Формат даты
		parsedDate, err := time.Parse(layout, dateStr)
		if err != nil {
			log.Printf("error: Failed to parse release date: %v", err)
			return
		}
		fmt.Println(group, song)
		// Создаем фиктивный ответ
		response := SongDetail{
			ReleaseDate: parsedDate,
			Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
			Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		}

		// Возвращаем успешный ответ
		c.JSON(http.StatusOK, response)
	})

	// Запускаем сервер
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("error: Failed to start server: %v", err)
	}
}
