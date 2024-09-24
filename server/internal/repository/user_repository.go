package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/pgvector/pgvector-go"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, email, role, taste, created_at) 
              VALUES ($1, $2, $3, $4, $5)`

	user.CreatedAt = time.Now()
	user.Taste = pgvector.NewVector([]float32{})

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Role, user.Taste, user.CreatedAt,
	)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, role, taste, created_at 
              FROM users 
              WHERE email = $1`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Role, &user.Taste, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.Email = email
	return &user, nil
}
