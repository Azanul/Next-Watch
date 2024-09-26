package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.54

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azanul/Next-Watch/graph/model"
	"github.com/Azanul/Next-Watch/internal/auth"
	"github.com/google/uuid"
)

// RateMovie is the resolver for the rateMovie field.
func (r *mutationResolver) RateMovie(ctx context.Context, movieID string, score float64) (*model.Rating, error) {
	currentUser, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Validate inputs
	var movieUUID uuid.UUID
	if movieUUID, err = uuid.Parse(movieID); err != nil {
		return nil, errors.New("invalid movie id")
	}
	if score < 0 || score > 5 {
		return nil, errors.New("rating score must be between 0 and 5")
	}

	// Call service to rate movie
	rating, err := r.RatingService.RateMovie(ctx, currentUser, movieUUID, float32(score))
	if err != nil {
		return nil, err
	}

	// Convert internal model to GraphQL model
	return &model.Rating{
		ID:    rating.ID.String(),
		User:  &model.User{ID: rating.UserID.String()},
		Movie: &model.Movie{ID: rating.MovieID.String()},
		Score: float64(rating.Score),
	}, nil
}

// DeleteRating is the resolver for the deleteRating field.
func (r *mutationResolver) DeleteRating(ctx context.Context, id string) (bool, error) {
	// Get current user from context
	currentUser, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return false, err
	}

	ratingID, err := uuid.Parse(id)
	if err != nil {
		return false, errors.New("invalid rating ID")
	}

	// Fetch the rating
	rating, err := r.RatingService.GetRatingByID(ctx, ratingID)
	if err != nil {
		return false, err
	}
	if rating == nil {
		return false, errors.New("rating not found")
	}

	// Check if the user is authorized to delete this rating
	isAdmin := currentUser.Role == "ADMIN"
	isOwner := rating.UserID == currentUser.ID

	if !isAdmin && !isOwner {
		return false, errors.New("not authorized to delete this rating")
	}

	// Delete the rating
	return r.RatingService.DeleteRating(ctx, ratingID)
}

// Movie is the resolver for the movie field.
func (r *queryResolver) Movie(ctx context.Context, id string) (*model.Movie, error) {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid movie ID")
	}

	// Fetch the movie
	movie, err := r.MovieService.GetMovieByID(ctx, movieID)
	if err != nil {
		return nil, err
	}
	if movie == nil {
		return nil, errors.New("movie not found")
	}
	return &model.Movie{
		ID:    movieID.String(),
		Title: movie.Title,
		Genre: movie.Genre,
		Year:  movie.Year,
		Wiki:  movie.Wiki,
		Plot:  movie.Plot,
		Cast:  movie.Cast,
	}, nil
}

// MovieByTitle is the resolver for the movieByTitle field.
func (r *queryResolver) MovieByTitle(ctx context.Context, title string) (*model.Movie, error) {
	// Fetch the movie
	movie, err := r.MovieService.GetMovieByTitle(ctx, title)
	if err != nil {
		return nil, err
	}
	if movie == nil {
		return nil, errors.New("movie not found")
	}
	return &model.Movie{
		ID:    movie.ID.String(),
		Title: movie.Title,
		Genre: movie.Genre,
		Year:  movie.Year,
		Wiki:  movie.Wiki,
		Plot:  movie.Plot,
		Cast:  movie.Cast,
	}, nil
}

// Movies is the resolver for the movies field.
func (r *queryResolver) Movies(ctx context.Context, page int, pageSize int) (*model.MovieConnection, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Default page size
	}

	moviePage, err := r.MovieService.GetMovies(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.MovieEdge, len(moviePage.Movies))
	for i, movie := range moviePage.Movies {
		edges[i] = &model.MovieEdge{
			Node: &model.Movie{
				ID:    movie.ID.String(),
				Title: movie.Title,
				Genre: movie.Genre,
				Year:  movie.Year,
				Wiki:  movie.Wiki,
				Plot:  movie.Plot,
				Cast:  movie.Cast,
			},
		}
	}

	return &model.MovieConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     moviePage.HasNextPage,
			HasPreviousPage: moviePage.HasPreviousPage,
		},
		TotalCount: moviePage.TotalCount,
	}, nil
}

// SearchMovies is the resolver for the searchMovies field.
func (r *queryResolver) SearchMovies(ctx context.Context, query string, page int, pageSize int) (*model.MovieConnection, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Default page size
	}

	moviePage, err := r.MovieService.SearchMovies(ctx, query, page, pageSize)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.MovieEdge, len(moviePage.Movies))
	for i, movie := range moviePage.Movies {
		edges[i] = &model.MovieEdge{
			Node: &model.Movie{
				ID:    movie.ID.String(),
				Title: movie.Title,
				Genre: movie.Genre,
				Year:  movie.Year,
				Wiki:  movie.Wiki,
				Plot:  movie.Plot,
				Cast:  movie.Cast,
			},
		}
	}

	return &model.MovieConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     moviePage.HasNextPage,
			HasPreviousPage: moviePage.HasPreviousPage,
		},
		TotalCount: moviePage.TotalCount,
	}, nil
}

// Recommendations is the resolver for the recommendations field.
func (r *queryResolver) Recommendations(ctx context.Context, page int, pageSize int) (*model.MovieConnection, error) {
	currentUser, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Default page size
	}

	moviePage, err := r.RecommendationService.GetSimilarMovies(ctx, currentUser.Taste, page, pageSize)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.MovieEdge, len(moviePage.Movies))
	for i, movie := range moviePage.Movies {
		edges[i] = &model.MovieEdge{
			Node: &model.Movie{
				ID:    movie.ID.String(),
				Title: movie.Title,
				Genre: movie.Genre,
				Year:  movie.Year,
				Wiki:  movie.Wiki,
				Plot:  movie.Plot,
				Cast:  movie.Cast,
			},
		}
	}

	return &model.MovieConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     moviePage.HasNextPage,
			HasPreviousPage: moviePage.HasPreviousPage,
		},
		TotalCount: moviePage.TotalCount,
	}, nil
}

// Ratings is the resolver for the ratings field.
func (r *queryResolver) Ratings(ctx context.Context, userID string) ([]*model.Rating, error) {
	panic(fmt.Errorf("not implemented: Ratings - ratings"))
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	panic(fmt.Errorf("not implemented: User - user"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
