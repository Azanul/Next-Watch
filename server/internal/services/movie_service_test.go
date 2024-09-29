package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/Azanul/Next-Watch/internal/repository"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMovieRepository is a mock of MovieRepository
type MockMovieRepository struct {
	mock.Mock
}

func (m *MockMovieRepository) GetMovies(ctx context.Context, searchTerm string, page, pageSize int) (*repository.MoviePage, error) {
	args := m.Called(ctx, searchTerm, page, pageSize)
	return args.Get(0).(*repository.MoviePage), args.Error(1)
}

func (m *MockMovieRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Movie, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Movie), args.Error(1)
}

func (m *MockMovieRepository) GetByTitle(ctx context.Context, title string) (*models.Movie, error) {
	args := m.Called(ctx, title)
	return args.Get(0).(*models.Movie), args.Error(1)
}

func (m *MockMovieRepository) GetSimilarMovies(ctx context.Context, embedding pgvector.Vector, page, pageSize int) (*repository.MoviePage, error) {
	args := m.Called(ctx, embedding, page, pageSize)
	return args.Get(0).(*repository.MoviePage), args.Error(1)
}

func (m *MockMovieRepository) Create(ctx context.Context, movie *models.Movie) error {
	args := m.Called(ctx, movie)
	return args.Error(0)
}

func (m *MockMovieRepository) Update(ctx context.Context, movie *models.Movie) error {
	args := m.Called(ctx, movie)
	return args.Error(0)
}

func (m *MockMovieRepository) Delete(ctx context.Context, rating *models.Rating) error {
	args := m.Called(ctx, rating)
	return args.Error(0)
}

func TestMovieService_GetMovies(t *testing.T) {
	mockRepo := new(MockMovieRepository)
	service := NewMovieService(mockRepo)

	movieID := uuid.New()
	tests := []struct {
		name      string
		page      int
		pageSize  int
		mockSetup func()
		want      *repository.MoviePage
		wantErr   bool
	}{
		{
			name:     "Success",
			page:     1,
			pageSize: 10,
			mockSetup: func() {
				mockRepo.On("GetMovies", mock.Anything, "", 1, 10).Return(&repository.MoviePage{
					Movies:     []*models.Movie{{ID: movieID, Title: "Test Movie"}},
					TotalCount: 1,
				}, nil)
			},
			want: &repository.MoviePage{
				Movies:     []*models.Movie{{ID: movieID, Title: "Test Movie"}},
				TotalCount: 1,
			},
			wantErr: false,
		},
		{
			name:     "Error",
			page:     1,
			pageSize: 2,
			mockSetup: func() {
				mockRepo.On("GetMovies", mock.Anything, "", 1, 2).Return((*repository.MoviePage)(nil), errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := service.GetMovies(context.Background(), tt.page, tt.pageSize)
			t.Log(tt.page, tt.name, tt.want, tt.wantErr, got, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieService.GetMovieByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil && tt.want != nil {
				assert.Equal(t, tt.want.TotalCount, got.TotalCount)
			}
		})
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	}
}

func TestMovieService_GetMovieByID(t *testing.T) {
	mockRepo := new(MockMovieRepository)
	service := NewMovieService(mockRepo)

	movieID := uuid.New()
	tests := []struct {
		name      string
		movieID   uuid.UUID
		mockSetup func()
		want      *models.Movie
		wantErr   bool
	}{
		{
			name:    "Success",
			movieID: movieID,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, movieID).Return(&models.Movie{
					ID:    movieID,
					Title: "Test Movie",
				}, nil)
			},
			want:    &models.Movie{Title: "Test Movie"},
			wantErr: false,
		},
		{
			name:    "Not Found",
			movieID: uuid.New(),
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return((*models.Movie)(nil), nil)
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Error",
			movieID: uuid.New(),
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return((*models.Movie)(nil), errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := service.GetMovieByID(context.Background(), tt.movieID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieService.GetMovieByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil && tt.want != nil {
				assert.Equal(t, tt.want.Title, got.Title)
			}
		})
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	}
}
