package repositories

import "revel-dynamodb-v1/app/models"

type Repository interface {
	GetMovies() ([]*models.Movie, error)
}
