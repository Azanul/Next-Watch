package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type Movie struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Genre     string    `json:"genre"`
	Year      int       `json:"year"`
	Wiki      string    `json:"wiki"`
	Plot      string    `json:"plot"`
	Cast      string    `json:"cast"`
	Embedding pgvector.Vector
}

type User struct {
	ID        uuid.UUID       `json:"id"`
	Email     string          `json:"email"`
	Role      string          `json:"role"`
	Taste     pgvector.Vector `json:"-"`
	CreatedAt time.Time       `json:"createdAt"`
}

type Rating struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	MovieID   uuid.UUID `json:"movieId"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
