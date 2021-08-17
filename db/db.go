package db

import (
	"context"
	"github.com/Azanul/Next-Watch/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type DB interface {
	GetMovies(attr string) ([]*model.Movie, error)
}

type MongoDB struct {
	collection *mongo.Collection
}

func New(client *mongo.Client) *MongoDB {
	movies := client.Database("movieDB").Collection("movies")
	return &MongoDB{
		collection: movies,
	}
}

func (db MongoDB) GetMovies(attr string) ([]*model.Movie, error) {
	res, err := db.collection.Find(context.TODO(), db.filter(attr))
	if err != nil {
		log.Printf("Error while fetching movies: %s", err.Error())
		return nil, err
	}
	var p []*model.Movie
	err = res.All(context.TODO(), &p)
	if err != nil {
		log.Printf("Error while decoding movies: %s", err.Error())
		return nil, err
	}
	return p, nil
}

func (db MongoDB) filter(attr string) bson.D {
	return bson.D{{
		"attribute.name",
		bson.D{{
			"$regex",
			"^" + attr + ".*$",
		}, {
			"$options",
			"i",
		}},
	}}
}
