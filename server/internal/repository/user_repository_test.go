package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	tests := []struct {
		name      string
		user      *models.User
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			user: &models.User{
				ID:    uuid.New(),
				Email: "test@example.com",
				Name:  "Test User",
				Role:  "user",
			},
			mockSetup: func() {
				mock.ExpectExec("^INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Error",
			user: &models.User{
				ID:    uuid.New(),
				Email: "error@example.com",
				Name:  "Error User",
				Role:  "user",
			},
			mockSetup: func() {
				mock.ExpectExec("^INSERT INTO users").WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Create(context.Background(), tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, 512, len(tt.user.Taste.Slice()))
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	tests := []struct {
		name      string
		email     string
		mockSetup func()
		want      *models.User
		wantErr   bool
	}{
		{
			name:  "Success",
			email: "test@example.com",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "role", "taste", "created_at"}).
					AddRow(uuid.New(), "user", pgvector.NewVector(make([]float32, 512)), time.Now())
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WillReturnRows(rows)
			},
			want:    &models.User{},
			wantErr: false,
		},
		{
			name:  "Not Found",
			email: "notfound@example.com",
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:  "Error",
			email: "error@example.com",
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repo.GetByEmail(context.Background(), tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				assert.IsType(t, &models.User{}, got)
				assert.Equal(t, tt.email, got.Email)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	tests := []struct {
		name      string
		user      *models.User
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			user: &models.User{
				ID:    uuid.New(),
				Email: "update@example.com",
				Role:  "admin",
				Taste: pgvector.NewVector(make([]float32, 512)),
			},
			mockSetup: func() {
				mock.ExpectExec("^UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Error",
			user: &models.User{
				ID:    uuid.New(),
				Email: "error@example.com",
				Role:  "user",
				Taste: pgvector.NewVector(make([]float32, 512)),
			},
			mockSetup: func() {
				mock.ExpectExec("^UPDATE users").WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Update(context.Background(), tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
