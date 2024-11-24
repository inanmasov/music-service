package models

import "time"

type Song struct {
	ID          int       `json:"id"`
	GroupName   string    `json:"group"`
	SongName    string    `json:"song"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

// ErrorResponse представляет структуру для ошибок
type ErrorResponse struct {
	Error string `json:"error"`
}
