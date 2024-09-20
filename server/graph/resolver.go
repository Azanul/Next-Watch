package graph

import "github.com/Azanul/Next-Watch/internal/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	services.RatingService
	services.MovieService
}
