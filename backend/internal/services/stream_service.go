package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
)

func GetStreams() ([]models.Stream, error) {
	return repositories.GetAllStreams()
}

func CreateStream(stream *models.Stream) error {
	return repositories.CreateStream(stream)
}
