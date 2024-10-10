package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inanmasov/music-service/internal/db"
	"github.com/inanmasov/music-service/internal/models"
)

// AddSong добавляет новую песню в библиотеку
// @Summary Add a new song to the library
// @Description Adds a new song with details like group, song name, release date, text, and link
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song details"
// @Success 201 {object} models.Song "Song created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input data"
// @Failure 500 {object} models.ErrorResponse "Failed to call external API or insert data into database"
// @Router /songs [post]
func AddSong(c *gin.Context) {
	log.Println("info: Starting AddSong handler")

	// Структура для получения данных из тела запроса
	var input struct {
		Group string `json:"group" binding:"required"`
		Song  string `json:"song" binding:"required"`
	}

	// Привязываем данные из запроса к структуре input
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("error: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("info: Received request to add song - Group: %s, Song: %s", input.Group, input.Song)

	songDetail, err := GetSongInfoFromAPI(input.Group, input.Song)
	if err != nil {
		log.Printf("error: Failed to get song info from external API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call external API"})
		return
	}

	log.Printf("info: Retrieved song details from external API: %+v", songDetail)

	// Подключаемся к базе данных
	db, err := db.Initialize()
	if err != nil {
		log.Printf("error: Failed to connect to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	log.Println("info: Successfully connected to the database")

	// Подготавливаем SQL-запрос для добавления песни
	query := `
        INSERT INTO songs (group_name, song, release_date, text, link)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	// Выполняем запрос
	var songID int
	err = db.QueryRow(query, input.Group, input.Song, songDetail.ReleaseDate, songDetail.Text, songDetail.Link).Scan(&songID)
	if err != nil {
		log.Printf("error: Failed to insert song into database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert song into database"})
		return
	}

	log.Printf("info: Song successfully added to database with ID: %d", songID)

	// Возвращаем ответ с добавленной песней
	c.JSON(http.StatusCreated, gin.H{
		"id":           songID,
		"group_name":   input.Group,
		"song":         input.Song,
		"release_date": songDetail.ReleaseDate,
		"text":         songDetail.Text,
		"link":         songDetail.Link,
	})

	log.Println("info: Successfully completed AddSong handler")
}

func GetSongInfoFromAPI(group, song string) (models.Song, error) {
	log.Printf("debug: Mocking external API call for group: %s, song: %s", group, song)

	// Имитация ответа от внешнего API
	songDetails := models.Song{
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	log.Printf("info: Mocked song details returned: %+v", songDetails)

	// Возвращаем данные без ошибок
	return songDetails, nil
}
