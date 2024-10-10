package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/inanmasov/music-service/internal/db"
	_ "github.com/inanmasov/music-service/internal/models"
)

// GetSongText возвращает текст песни с пагинацией по куплетам
// @Summary Get song text by verses with pagination
// @Description Retrieves the song's text, paginated by verses, based on the song's ID
// @Tags songs
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of verses per page" default(2)
// @Success 200 {object} map[string]string "Song text retrieved successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid page or limit number"
// @Failure 404 {object} models.ErrorResponse "Song not found or no verses on this page"
// @Failure 500 {object} models.ErrorResponse "Failed to connect to database or retrieve song text"
// @Router /songs/{id}/text [get]
func GetSongText(c *gin.Context) {
	log.Println("info: Starting GetSongText handler")

	// Получаем ID песни из URL параметров
	id := c.Param("id")

	log.Printf("info: Request to get song text with ID: %s", id)

	// Получаем параметры пагинации из URL
	pageParam := c.DefaultQuery("page", "1")
	limitParam := c.DefaultQuery("limit", "2")

	log.Printf("info: Pagination params received: page = %s, limit = %s", pageParam, limitParam)

	// Преобразуем параметры в числа
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		log.Printf("error: Invalid page number: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		log.Printf("error: Invalid limit number: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}

	log.Printf("debug: Parsed pagination params: page = %d, limit = %d", page, limit)

	// Подключаемся к базе данных
	db, err := db.Initialize()
	if err != nil {
		log.Printf("error: Failed to connect to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	log.Println("info: Successfully connected to the database")

	// Выполняем SQL-запрос для получения текста песни
	var text string
	err = db.QueryRow("SELECT text FROM songs WHERE id = $1", id).Scan(&text)
	if err == sql.ErrNoRows {
		log.Printf("info: Song with ID %s not found", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	} else if err != nil {
		log.Printf("error: Failed to retrieve song text for ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve song text"})
		return
	}

	log.Println("debug: Successfully retrieved song text from database")

	// Разбиваем текст песни на куплеты по разделителю (например, "\n\n" между куплетами)
	verses := strings.Split(text, "\n\n")

	// Вычисляем срез куплетов для текущей страницы
	start := (page - 1) * limit
	if start >= len(verses) {
		log.Printf("info: No verses on page %d", page)
		c.JSON(http.StatusNotFound, gin.H{"error": "No verses on this page"})
		return
	}

	end := start + limit
	if end > len(verses) {
		end = len(verses)
	}

	paginatedVerses := verses[start:end]

	log.Printf("info: Text song with ID %s get successfully", id)

	// Возвращаем куплеты в ответе
	c.JSON(http.StatusOK, gin.H{
		"page":    page,
		"limit":   limit,
		"verses":  paginatedVerses,
		"total":   len(verses),
		"message": "Song text retrieved successfully",
	})
}
