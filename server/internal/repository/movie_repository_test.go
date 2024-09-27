package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
)

func TestMovieRepository_GetMovies(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewMovieRepository(db)

	tests := []struct {
		name       string
		searchTerm string
		page       int
		pageSize   int
		mockSetup  func()
		want       *MoviePage
		wantErr    bool
	}{
		{
			name:       "Success - No search term",
			searchTerm: "",
			page:       1,
			pageSize:   10,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "genre", "year", "wiki", "plot", "director", "cast"}).
					AddRow(uuid.New(), "Movie 1", "Action", 2021, "wiki1", "plot1", "director1", "cast1").
					AddRow(uuid.New(), "Movie 2", "Comedy", 2022, "wiki2", "plot2", "director2", "cast2")
				mock.ExpectQuery("^SELECT (.+) FROM movies").WillReturnRows(rows)
				mock.ExpectQuery("^SELECT COUNT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
			},
			want: &MoviePage{
				Movies:          []*models.Movie{},
				TotalCount:      2,
				HasNextPage:     false,
				HasPreviousPage: false,
			},
			wantErr: false,
		},
		{
			name:       "Success - With search term",
			searchTerm: "Action",
			page:       1,
			pageSize:   10,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "genre", "year", "wiki", "plot", "director", "cast"}).
					AddRow(uuid.New(), "Action Movie", "Action", 2021, "wiki1", "plot1", "director1", "cast1")
				mock.ExpectQuery("^SELECT (.+) FROM movies WHERE").WillReturnRows(rows)
				mock.ExpectQuery("^SELECT COUNT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			want: &MoviePage{
				Movies:          []*models.Movie{},
				TotalCount:      1,
				HasNextPage:     false,
				HasPreviousPage: false,
			},
			wantErr: false,
		},
		{
			name:       "Error - Database query fails",
			searchTerm: "",
			page:       1,
			pageSize:   10,
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM movies").WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repo.GetMovies(context.Background(), tt.searchTerm, tt.page, tt.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieRepository.GetMovies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want.TotalCount, got.TotalCount)
				assert.Equal(t, tt.want.HasNextPage, got.HasNextPage)
				assert.Equal(t, tt.want.HasPreviousPage, got.HasPreviousPage)
			}
		})
	}
}

func TestMovieRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewMovieRepository(db)

	tests := []struct {
		name      string
		id        uuid.UUID
		mockSetup func()
		want      *models.Movie
		wantErr   bool
	}{
		{
			name: "Success",
			id:   uuid.New(),
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "genre", "year", "wiki", "plot", "cast", "embedding"}).
					AddRow(uuid.New(), "Movie 1", "Action", 2021, "wiki1", "plot1", "cast1", pgvector.NewVector([]float32{1, 2, 3}))
				mock.ExpectQuery("^SELECT (.+) FROM movies WHERE").WillReturnRows(rows)
			},
			want:    &models.Movie{},
			wantErr: false,
		},
		{
			name: "Not Found",
			id:   uuid.New(),
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM movies WHERE").WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Error",
			id:   uuid.New(),
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM movies WHERE").WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repo.GetByID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				assert.IsType(t, &models.Movie{}, got)
			}
		})
	}
}

func TestMovieRepository_GetByTitle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewMovieRepository(db)

	tests := []struct {
		name      string
		title     string
		mockSetup func()
		want      *models.Movie
		wantErr   bool
	}{
		{
			name:  "Success",
			title: "Test Movie",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "genre", "year", "wiki", "plot", "director", "cast"}).
					AddRow(uuid.New(), "Action", 2021, "wiki", "plot", "director", "cast")
				mock.ExpectQuery("^SELECT (.+) FROM movies WHERE").WillReturnRows(rows)
			},
			want:    &models.Movie{Title: "Test Movie"},
			wantErr: false,
		},
		{
			name:  "Not Found",
			title: "Non-existent Movie",
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM movies WHERE").WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:  "Error",
			title: "Error Movie",
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM movies WHERE").WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repo.GetByTitle(context.Background(), tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieRepository.GetByTitle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				assert.Equal(t, tt.title, got.Title)
			}
		})
	}
}

func TestMovieRepository_GetSimilarMovies(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewMovieRepository(db)

	tests := []struct {
		name      string
		embedding pgvector.Vector
		page      int
		pageSize  int
		mockSetup func()
		want      *MoviePage
		wantErr   bool
	}{
		{
			name:      "Success",
			embedding: pgvector.NewVector([]float32{1, 2, 3}),
			page:      1,
			pageSize:  10,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "genre", "year", "wiki", "plot", "director", "cast"}).
					AddRow(uuid.New(), "Similar Movie 1", "Action", 2021, "wiki1", "plot1", "director1", "cast1").
					AddRow(uuid.New(), "Similar Movie 2", "Comedy", 2022, "wiki2", "plot2", "director2", "cast2")
				mock.ExpectQuery("^SELECT (.+) FROM movies ORDER BY").WillReturnRows(rows)
				mock.ExpectQuery("^SELECT COUNT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
			},
			want: &MoviePage{
				Movies:          []*models.Movie{},
				TotalCount:      2,
				HasNextPage:     false,
				HasPreviousPage: false,
			},
			wantErr: false,
		},
		{
			name:      "Error - Database query fails",
			embedding: pgvector.NewVector([]float32{1, 2, 3}),
			page:      1,
			pageSize:  10,
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM movies ORDER BY").WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repo.GetSimilarMovies(context.Background(), tt.embedding, tt.page, tt.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieRepository.GetSimilarMovies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want.TotalCount, got.TotalCount)
				assert.Equal(t, tt.want.HasNextPage, got.HasNextPage)
				assert.Equal(t, tt.want.HasPreviousPage, got.HasPreviousPage)
			}
		})
	}
}

func TestMovieRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewMovieRepository(db)

	tests := []struct {
		name      string
		movie     *models.Movie
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			movie: &models.Movie{
				Title:     "New Movie",
				Genre:     "Action",
				Year:      2023,
				Wiki:      "wiki",
				Plot:      "plot",
				Director:  "director",
				Cast:      "cast",
				Embedding: pgvector.NewVector([]float32{1, 2, 3}),
			},
			mockSetup: func() {
				mock.ExpectExec("^INSERT INTO movies").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Error",
			movie: &models.Movie{
				Title: "Error Movie",
			},
			mockSetup: func() {
				mock.ExpectExec("^INSERT INTO movies").WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Create(context.Background(), tt.movie)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMovieRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewMovieRepository(db)

	tests := []struct {
		name      string
		movie     *models.Movie
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			movie: &models.Movie{
				ID:        uuid.New(),
				Title:     "Updated Movie",
				Genre:     "Action",
				Year:      2023,
				Wiki:      "updated wiki",
				Plot:      "updated plot",
				Director:  "updated director",
				Cast:      "updated cast",
				Embedding: pgvector.NewVector([]float32{1, 2, 3}),
			},
			mockSetup: func() {
				mock.ExpectExec("^UPDATE movies").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Error",
			movie: &models.Movie{
				ID:    uuid.New(),
				Title: "Error Movie",
			},
			mockSetup: func() {
				mock.ExpectExec("^UPDATE movies").WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Update(context.Background(), tt.movie)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMovieRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewMovieRepository(db)

	tests := []struct {
		name      string
		rating    *models.Rating
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			rating: &models.Rating{
				ID: uuid.New(),
			},
			mockSetup: func() {
				mock.ExpectExec("^DELETE movies").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Error",
			rating: &models.Rating{
				ID: uuid.New(),
			},
			mockSetup: func() {
				mock.ExpectExec("^DELETE movies").WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Delete(context.Background(), tt.rating)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
