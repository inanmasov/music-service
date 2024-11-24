package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/inanmasov/music-service/internal/db"
	"github.com/inanmasov/music-service/internal/logger"
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
	log := logger.GetLogger()
	log.Info("Starting AddSong handler")

	// Структура для получения данных из тела запроса
	var input struct {
		Group string `json:"group" binding:"required"`
		Song  string `json:"song" binding:"required"`
	}

	// Привязываем данные из запроса к структуре input
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Debugf("Received request to add song - Group: %s, Song: %s", input.Group, input.Song)

	songDetail, err := GetSongInfoFromAPI(input.Group, input.Song)
	if err != nil {
		log.Errorf("Failed to get song info from external API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call external API"})
		return
	}

	log.Debugf("Retrieved song details from external API: %+v", songDetail)

	// Подключаемся к базе данных
	db, err := db.Initialize()
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
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

	var groupID int
	// Проверяем, существует ли группа с указанным именем
	err = db.QueryRow("SELECT id FROM groups WHERE name = $1", input.Group).Scan(&groupID)
	if err == sql.ErrNoRows {
		// Группа не найдена, добавляем новую
		log.Debugf("Group not found, adding new group: %s", input.Group)
		err = db.QueryRow("INSERT INTO groups (name) VALUES ($1) RETURNING id", input.Group).Scan(&groupID)
		if err != nil {
			log.Errorf("Failed to insert group into database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert group into database"})
			return
		}
		log.Debugf("New group added with ID: %d", groupID)
	} else if err != nil {
		// Обработка других ошибок
		log.Errorf("Failed to check if group exists: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if group exists"})
		return
	}

	log.Debugf("Group successfully added with ID: %d", groupID)

	// Подготавливаем SQL-запрос для добавления песни
	query := `
		INSERT INTO songs (group_id, song, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`

	// Выполняем запрос
	var songID int
	err = db.QueryRow(query, groupID, input.Song, songDetail.ReleaseDate, songDetail.Text, songDetail.Link).Scan(&songID)
	if err != nil {
		log.Errorf("Failed to insert song into database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert song into database"})
		return
	}

	log.Debugf("Song successfully added to database with ID: %d", songID)

	// Возвращаем ответ с добавленной песней
	c.JSON(http.StatusCreated, gin.H{
		"id":           songID,
		"group_name":   input.Group,
		"song":         input.Song,
		"release_date": songDetail.ReleaseDate,
		"text":         songDetail.Text,
		"link":         songDetail.Link,
	})

	log.Info("Successfully completed AddSong handler")
}

func GetSongInfoFromAPI(group, song string) (models.Song, error) {
	log := logger.GetLogger()

	// Формируем полный URL с параметрами запроса
	url := fmt.Sprintf("http://music-api:8080/info?group=%s&song=%s", url.QueryEscape(group), url.QueryEscape(song))

	// Отправляем GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Get API request: %v", err)
		return models.Song{}, err
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Reading response body: %v", err)
		return models.Song{}, err
	}

	var songDetail models.Song
	err = json.Unmarshal(body, &songDetail)
	if err != nil {
		log.Errorf("unmarshalling JSON response: %v", err)
		return models.Song{}, err
	}

	// Теперь, используя данные из SongDetail, создаем объект Song
	return songDetail, nil
}
