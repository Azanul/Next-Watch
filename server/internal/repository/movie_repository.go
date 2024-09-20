// internal/repository/rating_repository.go

package repository

import (
	"context"
	"database/sql"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/google/uuid"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) GetByTitle(ctx context.Context, title string) (*models.Movie, error) {
	query := `SELECT id, genre, year, wiki, plot, cast 
              FROM movies 
              WHERE title = $1`

	var movie models.Movie
	err := r.db.QueryRowContext(ctx, query, title).Scan(
		&movie.ID, &movie.Genre, &movie.Year, &movie.Wiki, &movie.Plot, &movie.Cast,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	movie.Title = title
	return &movie, nil
}

func (r *MovieRepository) Create(ctx context.Context, movie *models.Movie) error {
	query := `INSERT INTO movies (id, title, genre, year, wiki, plot, cast) 
              VALUES ($1, $2, $3, $4, $5, $6)`

	movie.ID = uuid.New()

	_, err := r.db.ExecContext(ctx, query,
		movie.ID, movie.Genre, movie.Year, movie.Wiki, movie.Plot, movie.Cast,
	)
	return err
}

func (r *MovieRepository) Update(ctx context.Context, movie *models.Movie) error {
	query := `UPDATE movies 
              SET title = $1, genre = $2, year = $3, wiki = $4, plot = $5, cast = $6,
              WHERE id = $7`

	_, err := r.db.ExecContext(ctx, query, movie.Title, movie.Genre, movie.Year, movie.Wiki, movie.Plot, movie.Cast, movie.ID)
	return err
}

func (r *MovieRepository) Delete(ctx context.Context, rating *models.Rating) error {
	query := `DELETE movies 
              WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, rating.ID)
	return err
}
