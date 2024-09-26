package services

import (
	"context"
	"errors"
	"math"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/Azanul/Next-Watch/internal/repository"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type RatingService struct {
	ratingRepo *repository.RatingRepository
	movieRepo  *repository.MovieRepository
	userRepo   *repository.UserRepository
}

func NewRatingService(ratingRepo *repository.RatingRepository, movieRepo *repository.MovieRepository, userRepo *repository.UserRepository) *RatingService {
	return &RatingService{
		ratingRepo: ratingRepo,
		movieRepo:  movieRepo,
		userRepo:   userRepo,
	}
}

func (s *RatingService) RateMovie(ctx context.Context, user *models.User, movieID uuid.UUID, score float32) (*models.Rating, error) {
	// Validate movie exists
	movie, err := s.movieRepo.GetByID(ctx, movieID)
	if err != nil {
		return nil, err
	}
	if movie == nil {
		return nil, errors.New("movie not found")
	}

	// Check for existing rating
	existingRating, err := s.ratingRepo.GetByUserAndMovie(ctx, user.ID, movieID)
	if err != nil {
		return nil, err
	}

	var rating *models.Rating

	if existingRating != nil {
		// Update existing rating
		existingRating.Score = score
		err = s.ratingRepo.Update(ctx, existingRating)
		if err != nil {
			return nil, err
		}
		rating = existingRating
	} else {
		// Create new rating
		newRating := &models.Rating{
			UserID:  user.ID,
			MovieID: movieID,
			Score:   score,
		}
		err = s.ratingRepo.Create(ctx, newRating)
		if err != nil {
			return nil, err
		}
		rating = newRating
	}

	// Update user's taste
	err = s.updateUserTaste(ctx, user, movie, score)
	if err != nil {
		return nil, err
	}

	return rating, nil
}

func (s *RatingService) updateUserTaste(ctx context.Context, user *models.User, movie *models.Movie, score float32) error {
	tasteVector, movieEmbeddingVector := user.Taste.Slice(), movie.Embedding.Slice()
	if len(tasteVector) != len(movieEmbeddingVector) {
		return errors.New("taste and embedding dimensions do not match")
	}

	// Calculate the weight based on the score
	weight := (score - 2.5) / 2.5 // Normalize score to [-1, 1]

	for i := range tasteVector {
		tasteVector[i] += weight * movieEmbeddingVector[i]
	}

	// Normalize the taste vector
	magnitude := float32(0)
	for _, v := range tasteVector {
		magnitude += v * v
	}
	magnitude = float32(math.Sqrt(float64(magnitude)))

	for i := range tasteVector {
		tasteVector[i] /= magnitude
	}
	user.Taste = pgvector.NewVector(tasteVector)

	return s.userRepo.Update(ctx, user)
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
