package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/google/uuid"
)

type RatingRepository struct {
	db *sql.DB
}

// Checking if RatingRepository implements RatingRepositoryInterface during compile time
var _ RatingRepositoryInterface = (*RatingRepository)(nil)

func NewRatingRepository(db *sql.DB) *RatingRepository {
	return &RatingRepository{db: db}
}

func (r *RatingRepository) GetByID(ctx context.Context, ratingID uuid.UUID) (*models.Rating, error) {
	query := `SELECT user_id, movie_id, score, created_at, updated_at 
              FROM ratings 
              WHERE id = $1`

	var rating models.Rating
	err := r.db.QueryRowContext(ctx, query, ratingID).Scan(
		&rating.UserID, &rating.MovieID, &rating.Score, &rating.CreatedAt, &rating.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rating.ID = ratingID
	return &rating, nil
}

func (r *RatingRepository) GetByUserAndMovie(ctx context.Context, userID, movieID uuid.UUID) (*models.Rating, error) {
	query := `SELECT id, user_id, movie_id, score, created_at, updated_at 
              FROM ratings 
              WHERE user_id = $1 AND movie_id = $2`

	var rating models.Rating
	err := r.db.QueryRowContext(ctx, query, userID, movieID).Scan(
		&rating.ID, &rating.UserID, &rating.MovieID, &rating.Score, &rating.CreatedAt, &rating.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *RatingRepository) Create(ctx context.Context, rating *models.Rating) error {
	query := `INSERT INTO ratings (id, user_id, movie_id, score, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`

	rating.ID = uuid.New()
	rating.CreatedAt = time.Now()
	rating.UpdatedAt = rating.CreatedAt

	_, err := r.db.ExecContext(ctx, query,
		rating.ID, rating.UserID, rating.MovieID, rating.Score, rating.CreatedAt, rating.UpdatedAt,
	)
	return err
}

func (r *RatingRepository) Update(ctx context.Context, rating *models.Rating) error {
	query := `UPDATE ratings 
              SET score = $1, updated_at = $2 
              WHERE id = $3`

	rating.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query, rating.Score, rating.UpdatedAt, rating.ID)
	return err
}

func (r *RatingRepository) Delete(ctx context.Context, ratingID uuid.UUID) (*models.Rating, error) {
	query := `DELETE FROM ratings 
              WHERE id = $1
              RETURNING id, user_id, movie_id, score, created_at, updated_at`

	var deletedRating models.Rating
	err := r.db.QueryRowContext(ctx, query, ratingID).Scan(
		&deletedRating.ID,
		&deletedRating.UserID,
		&deletedRating.MovieID,
		&deletedRating.Score,
		&deletedRating.CreatedAt,
		&deletedRating.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &deletedRating, nil
}
