package repository

import (
	"context"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type MovieRepositoryInterface interface {
	GetMovies(ctx context.Context, searchTerm string, page, pageSize int) (*MoviePage, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Movie, error)
	GetByTitle(ctx context.Context, title string) (*models.Movie, error)
	GetSimilarMovies(ctx context.Context, embedding pgvector.Vector, page, pageSize int) (*MoviePage, error)
	Create(ctx context.Context, movie *models.Movie) error
	Update(ctx context.Context, movie *models.Movie) error
	Delete(ctx context.Context, rating *models.Rating) error
}

type RatingRepositoryInterface interface {
	GetByID(ctx context.Context, ratingID uuid.UUID) (*models.Rating, error)
	GetByUserAndMovie(ctx context.Context, userID, movieID uuid.UUID) (*models.Rating, error)
	Create(ctx context.Context, rating *models.Rating) error
	Update(ctx context.Context, rating *models.Rating) error
	Delete(ctx context.Context, ratingID uuid.UUID) (*models.Rating, error)
}

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}
