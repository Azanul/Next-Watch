// internal/repository/rating_repository.go

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

func NewRatingRepository(db *sql.DB) *RatingRepository {
	return &RatingRepository{db: db}
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

func (r *RatingRepository) Delete(ctx context.Context, rating *models.Rating) error {
	query := `DELETE ratings 
              WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, rating.ID)
	return err
}
