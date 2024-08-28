package repositories

import (
	"backend/internal/db"
	"backend/internal/models"
)

func GetAllStreams() ([]models.Stream, error) {
	var streams []models.Stream
	result := db.DB.Find(&streams)
	return streams, result.Error
}

func CreateStream(stream *models.Stream) error {
	result := db.DB.Create(stream)
	return result.Error
}
