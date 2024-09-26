package services

import (
	"context"
	"errors"

	"github.com/Azanul/Next-Watch/graph/model"
	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/Azanul/Next-Watch/internal/repository"
	"github.com/google/uuid"
)

type MovieService struct {
	movieRepo *repository.MovieRepository
}

func NewMovieService(movieRepo *repository.MovieRepository) *MovieService {
	return &MovieService{
		movieRepo: movieRepo,
	}
}

// func (s *MovieService) CreateMovie(ctx context.Context, movie *model.Movie) (*models.Movie, error) {
// 	movieID := uuid.New()

// 	return s.movieRepo.Create(ctx, &models.Movie{
// 		ID:    movieID,
// 		Title: movie.Title,
// 		Genre: movie.Genre,
// 		Year:  movie.Year,
// 		Wiki:  movie.Wiki,
// 		Plot:  movie.Plot,
// 		Cast:  movie.Cast,
// 	})
// }

func (s *MovieService) GetMovies(ctx context.Context, page, pageSize int) (*repository.MoviePage, error) {
	return s.movieRepo.GetMovies(ctx, "", page, pageSize)
}

func (s *MovieService) SearchMovies(ctx context.Context, searchTerm string, page, pageSize int) (*repository.MoviePage, error) {
	return s.movieRepo.GetMovies(ctx, searchTerm, page, pageSize)
}

func (s *MovieService) GetMovieByID(ctx context.Context, movieID uuid.UUID) (*models.Movie, error) {
	return s.movieRepo.GetByID(ctx, movieID)
}

func (s *MovieService) GetMovieByTitle(ctx context.Context, title string) (*models.Movie, error) {
	return s.movieRepo.GetByTitle(ctx, title)
}

func (s *MovieService) UpdateMovie(ctx context.Context, movie *model.Movie) (*models.Movie, error) {
	// Implementation to update movie in database
	return nil, errors.New("not implemented")
}

func (s *MovieService) DeleteMovie(ctx context.Context, movie *model.Movie) (*model.Movie, error) {
	// Implementation to create new movie in database
	return nil, errors.New("not implemented")
}
