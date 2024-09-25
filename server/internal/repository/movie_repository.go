// internal/repository/rating_repository.go

package repository

import (
	"context"
	"database/sql"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

type MoviePage struct {
	Movies          []*models.Movie
	TotalCount      int
	HasNextPage     bool
	HasPreviousPage bool
}

func (r *MovieRepository) GetMovies(ctx context.Context, page, pageSize int) (*MoviePage, error) {
	offset := (page - 1) * pageSize

	// Query to get the movies for the current page
	query := `
    SELECT id, title, genre, year, wiki, plot, director, "cast"
    FROM movies
    ORDER BY id
    LIMIT $1 OFFSET $2
    `

	rows, err := r.db.QueryContext(ctx, query, pageSize+1, offset) // We fetch one extra to determine if there's a next page
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*models.Movie
	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Year, &movie.Wiki, &movie.Plot, &movie.Director, &movie.Cast)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Query to get the total count of movies
	var totalCount int
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM movies").Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	hasNextPage := len(movies) > pageSize
	if hasNextPage {
		movies = movies[:pageSize]
	}

	return &MoviePage{
		Movies:          movies,
		TotalCount:      totalCount,
		HasNextPage:     hasNextPage,
		HasPreviousPage: page > 1,
	}, nil
}

func (r *MovieRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Movie, error) {
	query := `SELECT id, title, genre, year, wiki, plot, "cast" 
              FROM movies 
              WHERE id = $1`

	var movie models.Movie
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		id.String(), &movie.ID, &movie.Title, &movie.Genre, &movie.Year, &movie.Wiki, &movie.Plot, &movie.Cast,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepository) GetByTitle(ctx context.Context, title string) (*models.Movie, error) {
	query := `SELECT id, genre, year, wiki, plot, director, "cast" 
              FROM movies 
              WHERE title = $1`

	var movie models.Movie
	err := r.db.QueryRowContext(ctx, query, title).Scan(
		&movie.ID, &movie.Genre, &movie.Year, &movie.Wiki, &movie.Plot, &movie.Director, &movie.Cast,
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

func (r *MovieRepository) GetSimilarMovies(ctx context.Context, embedding pgvector.Vector, page, pageSize int) (*MoviePage, error) {
	offset := (page - 1) * pageSize

	// Query to get the movies similar to taste
	query := `SELECT id, title, genre, year, wiki, plot, director, "cast" FROM movies 
	ORDER BY embedding <-> $1 
	LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, embedding, pageSize+1, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*models.Movie
	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Year, &movie.Wiki, &movie.Plot, &movie.Director, &movie.Cast)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Query to get the total count of movies
	var totalCount int
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM movies").Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	hasNextPage := len(movies) > pageSize
	if hasNextPage {
		movies = movies[:pageSize]
	}

	return &MoviePage{
		Movies:          movies,
		TotalCount:      totalCount,
		HasNextPage:     hasNextPage,
		HasPreviousPage: page > 1,
	}, nil
}

func (r *MovieRepository) Create(ctx context.Context, movie *models.Movie) error {
	query := `INSERT INTO movies (id, title, genre, year, wiki, plot, director, "cast", embedding) 
              VALUES ($1, $2, $3, $4, $5, $6)`

	movie.ID = uuid.New()

	_, err := r.db.ExecContext(ctx, query,
		movie.ID, movie.Genre, movie.Year, movie.Wiki, movie.Plot, movie.Director, movie.Cast, movie.Embedding,
	)
	return err
}

func (r *MovieRepository) Update(ctx context.Context, movie *models.Movie) error {
	query := `UPDATE movies 
              SET title = $1, genre = $2, year = $3, wiki = $4, plot = $5, director = $6, "cast" = $7, embedding = $8
              WHERE id = $9`

	_, err := r.db.ExecContext(ctx, query, movie.Title, movie.Genre, movie.Year, movie.Wiki, movie.Plot, movie.Cast, movie.Embedding, movie.ID)
	return err
}

func (r *MovieRepository) Delete(ctx context.Context, rating *models.Rating) error {
	query := `DELETE movies 
              WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, rating.ID)
	return err
}
