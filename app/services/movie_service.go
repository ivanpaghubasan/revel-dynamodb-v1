package services

import (
	"revel-dynamodb-v1/app/models"
	"revel-dynamodb-v1/app/repositories"
)

type MovieService struct {
	repo repositories.Repository
}

func New(repository repositories.Repository) *MovieService {
	return &MovieService{
		repo: repository,
	}
}

func (s *MovieService) GetAllMovies() ([]*models.Movie, error) {
	return s.repo.GetMovies()
}
