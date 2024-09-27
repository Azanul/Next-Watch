package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRatingRepository struct {
	mock.Mock
}

func (m *MockRatingRepository) GetByID(ctx context.Context, ratingID uuid.UUID) (*models.Rating, error) {
	args := m.Called(ctx, ratingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Rating), args.Error(1)
}

func (m *MockRatingRepository) GetByUserAndMovie(ctx context.Context, userID, movieID uuid.UUID) (*models.Rating, error) {
	args := m.Called(ctx, userID, movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Rating), args.Error(1)
}

func (m *MockRatingRepository) Create(ctx context.Context, rating *models.Rating) error {
	args := m.Called(ctx, rating)
	return args.Error(0)
}

func (m *MockRatingRepository) Update(ctx context.Context, rating *models.Rating) error {
	args := m.Called(ctx, rating)
	return args.Error(0)
}

func (m *MockRatingRepository) Delete(ctx context.Context, ratingID uuid.UUID) (*models.Rating, error) {
	args := m.Called(ctx, ratingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Rating), args.Error(1)
}

func TestRatingService_RateMovie(t *testing.T) {
	mockRatingRepo := new(MockRatingRepository)
	mockMovieRepo := new(MockMovieRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewRatingService(mockRatingRepo, mockMovieRepo, mockUserRepo)

	ctx := context.Background()
	user := &models.User{ID: uuid.New(), Taste: pgvector.NewVector(make([]float32, 512))}
	movieID := uuid.New()
	score := float32(4.5)

	tests := []struct {
		name      string
		mockSetup func()
		want      *models.Rating
		wantErr   bool
	}{
		{
			name: "Success - New Rating",
			mockSetup: func() {
				mockMovieRepo.On("GetByID", ctx, movieID).Return(&models.Movie{ID: movieID, Embedding: pgvector.NewVector(make([]float32, 512))}, nil)
				mockRatingRepo.On("GetByUserAndMovie", ctx, user.ID, movieID).Return(nil, nil)
				mockRatingRepo.On("Create", ctx, mock.AnythingOfType("*models.Rating")).Return(nil)
				mockUserRepo.On("Update", ctx, mock.AnythingOfType("*models.User")).Return(nil)
			},
			want:    &models.Rating{UserID: user.ID, MovieID: movieID, Score: score},
			wantErr: false,
		},
		{
			name: "Success - Update Existing Rating",
			mockSetup: func() {
				mockMovieRepo.On("GetByID", ctx, movieID).Return(&models.Movie{ID: movieID, Embedding: pgvector.NewVector(make([]float32, 512))}, nil)
				mockRatingRepo.On("GetByUserAndMovie", ctx, user.ID, movieID).Return(&models.Rating{ID: uuid.New(), UserID: user.ID, MovieID: movieID, Score: 3.0}, nil)
				mockRatingRepo.On("Update", ctx, mock.AnythingOfType("*models.Rating")).Return(nil)
				mockUserRepo.On("Update", ctx, mock.AnythingOfType("*models.User")).Return(nil)
			},
			want:    &models.Rating{UserID: user.ID, MovieID: movieID, Score: score},
			wantErr: false,
		},
		{
			name: "Error - Movie Not Found",
			mockSetup: func() {
				mockMovieRepo.On("GetByID", ctx, movieID).Return(nil, nil)
			},
			want:    nil,
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := service.RateMovie(ctx, user, movieID, score)
			if (err != nil) != tt.wantErr {
				t.Errorf("RatingService.RateMovie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want.UserID, got.UserID)
				assert.Equal(t, tt.want.MovieID, got.MovieID)
				assert.Equal(t, tt.want.Score, got.Score)
			}
		})
	}
}

func TestRatingService_GetRatingByID(t *testing.T) {
	mockRatingRepo := new(MockRatingRepository)
	service := NewRatingService(mockRatingRepo, nil, nil)

	ctx := context.Background()
	ratingID := uuid.New()

	tests := []struct {
		name      string
		mockSetup func()
		want      *models.Rating
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func() {
				mockRatingRepo.On("GetByID", ctx, ratingID).Return(&models.Rating{
					ID:      ratingID,
					UserID:  uuid.New(),
					MovieID: uuid.New(),
					Score:   4.5,
				}, nil)
			},
			want:    &models.Rating{ID: ratingID, Score: 4.5},
			wantErr: false,
		},
		{
			name: "Not Found",
			mockSetup: func() {
				mockRatingRepo.On("GetByID", ctx, ratingID).Return(nil, nil)
			},
			want:    nil,
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := service.GetRatingByID(ctx, ratingID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RatingService.GetRatingByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Score, got.Score)
			}
		})
	}
}

func TestRatingService_GetRatingByUserAndMovie(t *testing.T) {
	mockRatingRepo := new(MockRatingRepository)
	service := NewRatingService(mockRatingRepo, nil, nil)

	ctx := context.Background()
	userID := uuid.New()
	movieID := uuid.New()

	tests := []struct {
		name      string
		mockSetup func()
		want      *models.Rating
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func() {
				mockRatingRepo.On("GetByUserAndMovie", ctx, userID, movieID).Return(&models.Rating{
					ID:      uuid.New(),
					UserID:  userID,
					MovieID: movieID,
					Score:   4.5,
				}, nil)
			},
			want:    &models.Rating{UserID: userID, MovieID: movieID, Score: 4.5},
			wantErr: false,
		},
		{
			name: "Not Found",
			mockSetup: func() {
				mockRatingRepo.On("GetByUserAndMovie", ctx, userID, movieID).Return(nil, nil)
			},
			want:    nil,
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := service.GetRatingByUserAndMovie(ctx, userID, movieID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RatingService.GetRatingByUserAndMovie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				assert.Equal(t, tt.want.UserID, got.UserID)
				assert.Equal(t, tt.want.MovieID, got.MovieID)
				assert.Equal(t, tt.want.Score, got.Score)
			}
		})
	}
}

func TestRatingService_DeleteRating(t *testing.T) {
	mockRatingRepo := new(MockRatingRepository)
	service := NewRatingService(mockRatingRepo, nil, nil)

	ctx := context.Background()
	ratingID := uuid.New()

	tests := []struct {
		name      string
		mockSetup func()
		want      bool
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func() {
				mockRatingRepo.On("Delete", ctx, ratingID).Return(&models.Rating{}, nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Error",
			mockSetup: func() {
				mockRatingRepo.On("Delete", ctx, ratingID).Return(nil, errors.New("database error"))
			},
			want:    false,
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := service.DeleteRating(ctx, ratingID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RatingService.DeleteRating() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
