package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/inanmasov/music-service/internal/db"
	"github.com/inanmasov/music-service/internal/logger"
	"github.com/inanmasov/music-service/internal/models"
)

// GetSongs возвращает список песен с фильтрацией и пагинацией
// @Summary Get songs list with filtering and pagination
// @Description Retrieves a paginated list of songs with optional filtering based on group, song name, release date, text, and link
// @Tags songs
// @Param groupName query string false "Group name for filtering"
// @Param song query string false "Song name for filtering"
// @Param releaseDate query string false "Release date for filtering"
// @Param text query string false "Text for filtering"
// @Param link query string false "Link for filtering"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of songs per page" default(10)
// @Success 200 {object} map[string]string "Songs retrieved successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid page or limit number"
// @Failure 500 {object} models.ErrorResponse "Failed to connect to database or retrieve songs"
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	log := logger.GetLogger()
	log.Info("Starting GetSongs handler")

	// Получение параметров фильтрации
	group := c.Query("groupName")
	songName := c.Query("song")
	releaseDate := c.Query("releaseDate")
	text := c.Query("text")
	link := c.Query("link")

	// Получение параметров пагинации
	pageParam := c.DefaultQuery("page", "1")
	limitParam := c.DefaultQuery("limit", "10") // По умолчанию 10 песен на страницу

	log.Debugf("Request to get songs with filters: group=%s, song=%s, releaseDate=%s, text=%s, link=%s", group, songName, releaseDate, text, link)

	// Преобразуем параметры в числа
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		log.Errorf("Invalid page number: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		log.Errorf("Invalid limit number: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}

	log.Debugf("Parsed pagination params: page=%d, limit=%d", page, limit)

	// Подключаемся к базе данных
	db, err := db.Initialize()
	if err != nil {
		log.Printf("error: Failed to connect to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	log.Info("Successfully connected to the database")

	query := `
		SELECT 
    		songs.id, 
    		groups.name AS group_name, 
    		songs.song AS song, 
    		songs.release_date, 
    		songs.text, 
    		songs.link
		FROM songs
		JOIN groups ON songs.group_id = groups.id
		WHERE 1=1`

	var args []interface{}

	// Добавляем фильтрацию по имени группы
	if group != "" {
		query += " AND groups.name ILIKE $" + strconv.Itoa(len(args)+1)
		args = append(args, "%"+group+"%")
		log.Debugf("Adding filter for group_name: %s", group)
	}

	// Добавляем фильтрацию по названию песни
	if songName != "" {
		query += " AND songs.song ILIKE $" + strconv.Itoa(len(args)+1)
		args = append(args, "%"+songName+"%")
		log.Debugf("Adding filter for song: %s", songName)
	}

	// Добавляем фильтрацию по дате выхода
	if releaseDate != "" {
		query += " AND songs.release_date = $" + strconv.Itoa(len(args)+1)
		args = append(args, releaseDate)
		log.Debugf("Adding filter for release_date: %s", releaseDate)
	}

	// Добавляем фильтрацию по тексту песни
	if text != "" {
		query += " AND songs.text ILIKE $" + strconv.Itoa(len(args)+1)
		args = append(args, "%"+text+"%")
		log.Debugf("Adding filter for text: %s", text)
	}

	// Добавляем фильтрацию по ссылке
	if link != "" {
		query += " AND songs.link ILIKE $" + strconv.Itoa(len(args)+1)
		args = append(args, "%"+link+"%")
		log.Debugf("Adding filter for link: %s", link)
	}

	// Добавляем пагинацию
	offset := (page - 1) * limit
	argLim := strconv.Itoa(len(args) + 1)
	argOff := strconv.Itoa(len(args) + 2)
	query += " LIMIT $" + argLim + " OFFSET $" + argOff
	args = append(args, limit, offset)

	log.Debugf("Adding pagination: LIMIT=%d, OFFSET=%d", limit, offset)

	// Выполняем SQL-запрос
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Errorf("Failed to retrieve songs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve songs"})
		return
	}
	defer rows.Close()

	// Считываем результаты
	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			log.Errorf("Failed to scan song: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan song"})
			return
		}
		songs = append(songs, song)
	}

	log.Infof("Retrieved %d songs successfully", len(songs))

	// Возвращаем песни в ответе
	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"songs": songs,
	})
}
