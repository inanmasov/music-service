package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inanmasov/music-service/internal/db"
	"github.com/inanmasov/music-service/internal/logger"
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
	log := logger.GetLogger()
	log.Info("Starting UpdateSong handler")

	// Получаем ID из параметров URL
	id := c.Param("id")

	log.Debugf("Request to update song with ID: %s", id)

	var rawData map[string]interface{}

	if err := c.ShouldBindJSON(&rawData); err != nil {
		log.Printf("error: Invalid JSON data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	log.Debugf("Raw data: %+v", rawData)

	db, err := db.Initialize()
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to database",
		})
		return
	}
	defer db.Close()

	log.Info("Successfully connected to the database")

	// Начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}

	// Откат транзакции в случае ошибки
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				log.Errorf("Failed to commit transaction: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			}
		}
	}()

	if value, exists := rawData["group"]; exists {
		// Обновляем имя группы, связанной с песней
		query := `
        	UPDATE groups
        	SET name = $1
        	WHERE id = (SELECT group_id FROM songs WHERE id = $2)
    	`
		_, err := tx.Exec(query, value, id)
		if err != nil {
			log.Errorf("Failed to update group name: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update group name"})
			return
		}
		log.Debugf("Group name successfully updated for song ID: %s", id)
	}

	// Обновляем только те поля, которые были переданы
	query := "UPDATE songs SET "
	params := []interface{}{}
	index := 1

	if value, exists := rawData["song"]; exists {
		query += "song = $" + fmt.Sprint(index) + ", "
		params = append(params, value)
		index++
	}

	if value, exists := rawData["releaseDate"]; exists {
		query += "release_date = $" + fmt.Sprint(index) + ", "
		params = append(params, value)
		index++
	}

	if value, exists := rawData["text"]; exists {
		query += "text = $" + fmt.Sprint(index) + ", "
		params = append(params, value)
		index++
	}

	if value, exists := rawData["link"]; exists {
		query += "link = $" + fmt.Sprint(index) + ", "
		params = append(params, value)
		index++
	}

	if index == 1 {
		c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
		return
	}

	// Удаляем последнюю запятую и добавляем условие WHERE
	query = query[:len(query)-2] + " WHERE id = $" + fmt.Sprint(index)
	params = append(params, id)

	log.Debugf("Prepared update query: %s with parameters: %+v", query, params)

	// Выполняем запрос
	result, err := db.Exec(query, params...)
	if err != nil {
		log.Errorf("Failed to update song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	// Проверяем, было ли обновление
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Warnf("No song found with ID: %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	log.Infof("Song with ID %s updated successfully", id)

	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}
