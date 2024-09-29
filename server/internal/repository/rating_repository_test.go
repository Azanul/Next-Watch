package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRatingRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	tests := []struct {
		name      string
		ratingID  uuid.UUID
		mockSetup func()
		want      *models.Rating
		wantErr   bool
	}{
		{
			name:     "Success",
			ratingID: uuid.New(),
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"user_id", "movie_id", "score", "created_at", "updated_at"}).
					AddRow(uuid.New(), uuid.New(), 5, time.Now(), time.Now())
				mock.ExpectQuery("^SELECT (.+) FROM ratings WHERE").WillReturnRows(rows)
			},
			want:    &models.Rating{},
			wantErr: false,
		},
		{
			name:     "Not Found",
			ratingID: uuid.New(),
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM ratings WHERE").WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:     "Error",
			ratingID: uuid.New(),
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM ratings WHERE").WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repo.GetByID(context.Background(), tt.ratingID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RatingRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				assert.IsType(t, &models.Rating{}, got)
			}
		})
	}
}

func TestRatingRepository_GetByUserAndMovie(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	tests := []struct {
		name      string
		userID    uuid.UUID
		movieID   uuid.UUID
		mockSetup func()
		want      *models.Rating
		wantErr   bool
	}{
		{
			name:    "Success",
			userID:  uuid.New(),
			movieID: uuid.New(),
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "movie_id", "score", "created_at", "updated_at"}).
					AddRow(uuid.New(), uuid.New(), uuid.New(), 5, time.Now(), time.Now())
				mock.ExpectQuery("^SELECT (.+) FROM ratings WHERE").WillReturnRows(rows)
			},
			want:    &models.Rating{},
			wantErr: false,
		},
		{
			name:    "Not Found",
			userID:  uuid.New(),
			movieID: uuid.New(),
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM ratings WHERE").WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Error",
			userID:  uuid.New(),
			movieID: uuid.New(),
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM ratings WHERE").WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repo.GetByUserAndMovie(context.Background(), tt.userID, tt.movieID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RatingRepository.GetByUserAndMovie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				assert.IsType(t, &models.Rating{}, got)
			}
		})
	}
}

func TestRatingRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	tests := []struct {
		name      string
		rating    *models.Rating
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			rating: &models.Rating{
				UserID:  uuid.New(),
				MovieID: uuid.New(),
				Score:   5,
			},
			mockSetup: func() {
				mock.ExpectExec("^INSERT INTO ratings").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Error",
			rating: &models.Rating{
				UserID:  uuid.New(),
				MovieID: uuid.New(),
				Score:   5,
			},
			mockSetup: func() {
				mock.ExpectExec("^INSERT INTO ratings").WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Create(context.Background(), tt.rating)
			if (err != nil) != tt.wantErr {
				t.Errorf("RatingRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRatingRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	tests := []struct {
		name      string
		rating    *models.Rating
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			rating: &models.Rating{
				ID:    uuid.New(),
				Score: 4,
			},
			mockSetup: func() {
				mock.ExpectExec("^UPDATE ratings").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Error",
			rating: &models.Rating{
				ID:    uuid.New(),
				Score: 4,
			},
			mockSetup: func() {
				mock.ExpectExec("^UPDATE ratings").WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Update(context.Background(), tt.rating)
			if (err != nil) != tt.wantErr {
				t.Errorf("RatingRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRatingRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	tests := []struct {
		name      string
		ratingID  uuid.UUID
		mockSetup func()
		want      *models.Rating
		wantErr   bool
	}{
		{
			name:     "Success",
			ratingID: uuid.New(),
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "movie_id", "score", "created_at", "updated_at"}).
					AddRow(uuid.New(), uuid.New(), uuid.New(), 5, time.Now(), time.Now())
				mock.ExpectQuery("^DELETE FROM ratings WHERE").WillReturnRows(rows)
			},
			want:    &models.Rating{},
			wantErr: false,
		},
		{
			name:     "Not Found",
			ratingID: uuid.New(),
			mockSetup: func() {
				mock.ExpectQuery("^DELETE FROM ratings WHERE").WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:     "Error",
			ratingID: uuid.New(),
			mockSetup: func() {
				mock.ExpectQuery("^DELETE FROM ratings WHERE").WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repo.Delete(context.Background(), tt.ratingID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RatingRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				assert.IsType(t, &models.Rating{}, got)
			}
		})
	}
}
