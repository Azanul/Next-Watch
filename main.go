package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://Azan:<password>@moviecluster.t9pn0.mongodb.net/movieDB?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			fmt.Println("Error disconnecting:" + err.Error())
		}
	}(client, ctx)

	movieDB := client.Database("movieDB")
	fmt.Println(movieDB)
	//var movieType = graphql.NewObject(
	//	graphql.ObjectConfig{
	//		Name: "Movie",
	//		Fields: graphql.Fields{
	//			"id": &graphql.Field{
	//				Type: graphql.String,
	//			},
	//			"name": &graphql.Field{
	//				Type: graphql.String,
	//			},
	//			"picture": &graphql.Field{
	//				Type: graphql.String,
	//			},
	//			"director": &graphql.Field{
	//				Type: graphql.Float,
	//			},
	//		},
	//	},
	//)
}
