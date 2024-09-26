package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/Azanul/Next-Watch/graph"
	"github.com/Azanul/Next-Watch/internal/auth"
	"github.com/Azanul/Next-Watch/internal/database"
	"github.com/Azanul/Next-Watch/internal/handlers"
	"github.com/Azanul/Next-Watch/internal/repository"
	"github.com/Azanul/Next-Watch/internal/services"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	logFile := initLogFile()
	defer logFile.Close()

	db := database.ConnectDB()

	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	ratingRepo := repository.NewRatingRepository(db)

	userService := services.NewUserService(userRepo, movieRepo)
	movieService := services.NewMovieService(movieRepo)
	ratingService := services.NewRatingService(ratingRepo, movieRepo, userRepo)
	recommendationService := services.NewRecommendationService(ratingRepo, movieRepo)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(
		graph.Config{
			Resolvers: &graph.Resolver{
				RatingService: *ratingService, MovieService: *movieService, RecommendationService: *recommendationService,
			},
			Directives: graph.DirectiveRoot{
				HasRole: hasRoleDirective,
			},
		},
	))
	restHandler := handlers.NewHandler(userService, auth.NewGoogleAuthClient())

	http.HandleFunc("/auth/signin/google", cors(restHandler.GoogleSignin))
	http.HandleFunc("/auth/callback/google", restHandler.GoogleCallback)

	http.HandleFunc("/query", cors(http.HandlerFunc(restHandler.AuthMiddleware(srv).ServeHTTP)))
	http.Handle("/", http.FileServer(getFrontendFileSystem()))

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func hasRoleDirective(ctx context.Context, obj interface{}, next graphql.Resolver, role string) (interface{}, error) {
	user, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check if the user has the required role
	if user.Role != role {
		return nil, errors.New("access denied")
	}

	// If the user has the role, continue to the next resolver
	return next(ctx)
}

func cors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:64139")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h(w, r)
	}
}

// Function to initialize the log file and set the log output to the file
func initLogFile() *os.File {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return logFile
}
