package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inanmasov/music-service/internal/db"
	"github.com/inanmasov/music-service/internal/models"
)

// UpdateSong обновляет данные о песне по её ID
// @Summary Update song details
// @Description Updates the song information by its ID. Only provided fields will be updated.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Updated song data"
// @Success 200 {object} models.Song "Song updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid JSON data"
// @Failure 404 {object} models.ErrorResponse "Song not found"
// @Failure 500 {object} models.ErrorResponse "Failed to connect to database or update song"
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	log.Println("info: Starting UpdateSong handler")

	// Получаем ID из параметров URL
	id := c.Param("id")

	log.Printf("info: Request to update song with ID: %s", id)

	var newSong models.Song

	// Парсим данные из тела запроса в структуру
	if err := c.ShouldBindJSON(&newSong); err != nil {
		log.Printf("error: Invalid JSON data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	log.Printf("debug: Parsed new song data: %+v", newSong)

	db, err := db.Initialize()
	if err != nil {
		log.Printf("error: Failed to connect to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to database",
		})
		return
	}
	defer db.Close()

	log.Println("info: Successfully connected to the database")

	// Обновляем только те поля, которые были переданы
	query := "UPDATE songs SET "
	params := []interface{}{}
	index := 1

	if newSong.GroupName != "" {
		query += "group_name = $" + fmt.Sprint(index) + ", "
		params = append(params, newSong.GroupName)
		index++
	}

	if newSong.SongName != "" {
		query += "song = $" + fmt.Sprint(index) + ", "
		params = append(params, newSong.SongName)
		index++
	}

	if newSong.ReleaseDate != "" {
		query += "release_date = $" + fmt.Sprint(index) + ", "
		params = append(params, newSong.ReleaseDate)
		index++
	}

	if newSong.Text != "" {
		query += "text = $" + fmt.Sprint(index) + ", "
		params = append(params, newSong.Text)
		index++
	}

	if newSong.Link != "" {
		query += "link = $" + fmt.Sprint(index) + ", "
		params = append(params, newSong.Link)
		index++
	}

	// Удаляем последнюю запятую и добавляем условие WHERE
	query = query[:len(query)-2] + " WHERE id = $" + fmt.Sprint(index)
	params = append(params, id)

	log.Printf("debug: Prepared update query: %s with parameters: %+v", query, params)

	// Выполняем запрос
	result, err := db.Exec(query, params...)
	if err != nil {
		log.Printf("error: Failed to update song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	// Проверяем, было ли обновление
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Printf("warning: No song found with ID: %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	log.Printf("info: Song with ID %s updated successfully", id)

	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}
