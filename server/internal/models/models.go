package models

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Genre string    `json:"genre"`
	Year  int       `json:"year"`
	Wiki  string    `json:"wiki"`
	Plot  string    `json:"plot"`
	Cast  string    `json:"cast"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type Rating struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	MovieID   uuid.UUID `json:"movieId"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
