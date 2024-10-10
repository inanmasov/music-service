package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Подключение к базе данных
func Initialize() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error: Error loading .env file")
		return nil, err
	}

	log.Println("info: Loaded environment variables")

	connect_db := "user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := sql.Open("postgres", connect_db)
	if err != nil {
		log.Printf("error: Failed to open database connection: %v", err)
		return nil, err
	}

	log.Println("info: Database connection established")

	if err = db.Ping(); err != nil {
		log.Printf("error: Database connection failed: %v", err)
		return nil, err
	}

	log.Println("info: Database ping successful")

	return db, nil
}
