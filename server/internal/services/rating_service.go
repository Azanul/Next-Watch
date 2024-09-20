package services

import (
	"context"
	"errors"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/Azanul/Next-Watch/internal/repository"
	"github.com/google/uuid"
)

type RatingService struct {
	ratingRepo *repository.RatingRepository
	movieRepo  *repository.MovieRepository
}

func NewRatingService(ratingRepo *repository.RatingRepository, movieRepo *repository.MovieRepository) *RatingService {
	return &RatingService{
		ratingRepo: ratingRepo,
		movieRepo:  movieRepo,
	}
}

func (s *RatingService) RateMovie(ctx context.Context, userID uuid.UUID, movieID uuid.UUID, score int) (*models.Rating, error) {
	// Validate movie exists
	movie, err := s.movieRepo.GetByID(ctx, movieID)
	if err != nil {
		return nil, err
	}
	if movie == nil {
		return nil, errors.New("movie not found")
	}

	// Check for existing rating
	existingRating, err := s.ratingRepo.GetByUserAndMovie(ctx, userID, movieID)
	if err != nil {
		return nil, err
	}

	if existingRating != nil {
		// Update existing rating
		existingRating.Score = score
		err = s.ratingRepo.Update(ctx, existingRating)
		if err != nil {
			return nil, err
		}
		return existingRating, nil
	}

	// Create new rating
	newRating := &models.Rating{
		UserID:  userID,
		MovieID: movieID,
		Score:   score,
	}
	err = s.ratingRepo.Create(ctx, newRating)
	if err != nil {
		return nil, err
	}

	return newRating, nil
}

func (s *RatingService) GetRatingByID(ctx context.Context, ratingID uuid.UUID) (*models.Rating, error) {
	rating, err := s.ratingRepo.GetByID(ctx, ratingID)
	if err != nil {
		return nil, err
	}

	return &models.Rating{
		ID:      ratingID,
		UserID:  rating.UserID,
		MovieID: rating.MovieID,
		Score:   rating.Score,
	}, nil
}

func (s *RatingService) GetRatingByUserAndMovie(ctx context.Context, userID, movieID uuid.UUID) (*models.Rating, error) {
	rating, err := s.ratingRepo.GetByUserAndMovie(ctx, userID, movieID)
	if err != nil {
		return nil, err
	}

	return &models.Rating{
		ID:      rating.ID,
		UserID:  rating.UserID,
		MovieID: rating.MovieID,
		Score:   rating.Score,
	}, nil
}

func (s *RatingService) DeleteRating(ctx context.Context, ratingID uuid.UUID) (bool, error) {
	_, err := s.ratingRepo.Delete(ctx, ratingID)
	if err != nil {
		return false, nil
	}
	return true, nil
}
