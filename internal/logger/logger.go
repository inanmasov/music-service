package logger

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	instance *logrus.Logger
	once     sync.Once
)

// GetLogger возвращает общий логгер
func GetLogger() *logrus.Logger {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("error: Loading .env file")
		}

		level := os.Getenv("LOG_LEVEL")

		parsedLevel, err := logrus.ParseLevel(level)
		if err != nil {
			log.Fatalf("Invalid log level: %v", err)
		}

		instance = logrus.New()
		instance.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
		instance.SetLevel(parsedLevel)
	})
	return instance
}
