package db

import (
	"context"
	"fmt"
	"github.com/Azanul/Next-Watch/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var mongoHost = "Azan:<password>@moviecluster.t9pn0.mongodb.net/movieDB"

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://" + mongoHost + "?retryWrites=true&w=majority"))
	if err != nil {
		fmt.Println("Error getting client: " + err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	return &DB{
		client: client,
	}
}

func (*DB) Save(movie *model.NewMovie) *model.Movie {
	collection := Connect().client.Database("movieDB").Collection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, movie)
	if err != nil {
		fmt.Println("Error inserting movie: " + err.Error())
	}

	inputActors := make([]*model.Attr, len(movie.Actors))
	for i, ele := range movie.Actors {
		inputActors[i] = (*model.Attr)(ele)
	}

	return &model.Movie{
		ID:      res.InsertedID.(primitive.ObjectID).Hex(),
		Name:    movie.Name,
		Poster:  movie.Poster,
		Actors:  inputActors,
		Watched: movie.Watched,
	}
}
