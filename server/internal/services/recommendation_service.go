package services

import (
	"context"

	"github.com/Azanul/Next-Watch/internal/repository"
	"github.com/pgvector/pgvector-go"
)

type RecommendationService struct {
	ratingRepo *repository.RatingRepository
	movieRepo  *repository.MovieRepository
}

func NewRecommendationService(ratingRepo *repository.RatingRepository, movieRepo *repository.MovieRepository) *RecommendationService {
	return &RecommendationService{
		ratingRepo: ratingRepo,
		movieRepo:  movieRepo,
	}
}

func (s *RecommendationService) GetSimilarMovies(ctx context.Context, taste_embedding pgvector.Vector, page, pageSize int) (*repository.MoviePage, error) {
	return s.movieRepo.GetSimilarMovies(ctx, taste_embedding, page, pageSize)
}
