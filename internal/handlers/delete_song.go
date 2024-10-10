package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inanmasov/music-service/internal/db"
	_ "github.com/inanmasov/music-service/internal/models"
)

// DeleteSong удаляет песню из базы данных
// @Summary Delete a song by ID
// @Description Deletes a song from the database by its ID
// @Tags songs
// @Param id path int true "Song ID"
// @Success 200 {object} map[string]string "Song deleted successfully"
// @Failure 404 {object} models.ErrorResponse "Song not found"
// @Failure 500 {object} models.ErrorResponse "Failed to delete song or connect to database"
// @Router /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	log.Println("info: Starting DeleteSong handler")

	id := c.Param("id")

	log.Printf("info: Request to delete song with ID: %s", id)

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

	// Выполняем запрос
	result, err := db.Exec("DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		log.Printf("error: Failed to delete song with ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete song from database",
		})
		return
	}

	log.Println("debug: SQL query executed to delete song")

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("error: Failed to retrieve affected rows for song with ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve result",
		})
		return
	}

	log.Printf("info: Number of affected rows: %d", rowsAffected)

	// Если не найдено, возвращаем ошибку 404
	if rowsAffected == 0 {
		log.Printf("info: Song with ID %s not found", id)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Song not found",
		})
		return
	}

	log.Printf("info: Song with ID %s deleted successfully", id)

	// Успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"message": "Song deleted successfully",
	})
}
