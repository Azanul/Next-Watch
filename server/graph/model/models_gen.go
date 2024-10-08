// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Movie struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Genre string `json:"genre"`
	Year  int    `json:"year"`
	Wiki  string `json:"wiki"`
	Plot  string `json:"plot"`
	Cast  string `json:"cast"`
}

type MovieConnection struct {
	Edges      []*MovieEdge `json:"edges"`
	PageInfo   *PageInfo    `json:"pageInfo"`
	TotalCount int          `json:"totalCount"`
}

type MovieEdge struct {
	Node *Movie `json:"node"`
}

type MovieInput struct {
	Title string `json:"title"`
	Genre string `json:"genre"`
	Year  int    `json:"year"`
	Wiki  string `json:"wiki"`
	Plot  string `json:"plot"`
	Cast  string `json:"cast"`
}

type Mutation struct {
}

type PageInfo struct {
	HasNextPage     bool `json:"hasNextPage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
}

type Query struct {
}

type Rating struct {
	ID    string  `json:"id"`
	User  *User   `json:"user"`
	Movie *Movie  `json:"movie"`
	Score float64 `json:"score"`
}

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	Role         string `json:"role"`
}
