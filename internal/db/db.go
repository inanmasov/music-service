package db

import (
	"database/sql"
	"os"

	"github.com/inanmasov/music-service/internal/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Подключение к базе данных
func Initialize() (*sql.DB, error) {
	log := logger.GetLogger()

	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
		return nil, err
	}

	log.Info("Loaded environment variables")

	connect_db := "user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := sql.Open("postgres", connect_db)
	if err != nil {
		log.Errorf("Failed to open database connection: %v", err)
		return nil, err
	}

	log.Info("Database connection established")

	if err = db.Ping(); err != nil {
		log.Errorf("Database connection failed: %v", err)
		return nil, err
	}

	log.Info("Database ping successful")

	return db, nil
}
