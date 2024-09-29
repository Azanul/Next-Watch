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

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

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
				Taste: pgvector.NewVector(make([]float32, 512)),
			},
			mockSetup: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
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
				Taste: pgvector.NewVector(make([]float32, 512)),
			},
			mockSetup: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := service.CreateUser(context.Background(), tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

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
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(&models.User{
					ID:    uuid.New(),
					Email: "test@example.com",
					Name:  "Test User",
					Role:  "user",
					Taste: pgvector.NewVector(make([]float32, 512)),
				}, nil)
			},
			want:    &models.User{Email: "test@example.com"},
			wantErr: false,
		},
		{
			name:  "Not Found",
			email: "notfound@example.com",
			mockSetup: func() {
				mockRepo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(nil, nil)
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:  "Error",
			email: "error@example.com",
			mockSetup: func() {
				mockRepo.On("GetByEmail", mock.Anything, "error@example.com").Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := service.GetUserByEmail(context.Background(), tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				assert.Equal(t, tt.want.Email, got.Email)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

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
				Name:  "Updated User",
				Role:  "admin",
				Taste: pgvector.NewVector(make([]float32, 512)),
			},
			mockSetup: func() {
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
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
				Taste: pgvector.NewVector(make([]float32, 512)),
			},
			mockSetup: func() {
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := service.UpdateUser(context.Background(), tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	}
}
