package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/parsaakbari1209/go-mongo-crud-rest-api/http"
	"github.com/parsaakbari1209/go-mongo-crud-rest-api/repository"
)

func main() {
	// create a database connection

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// password := "77J8Knd6vrhvVGjL"
	opts := options.Client().ApplyURI(
		"mongodb+srv://user insert here",
	).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	// create a repository
	repository := repository.NewRepository(client.Database("stfx"))

	// create an http server
	server := http.NewServer(repository)

	// create a gin router
	router := gin.Default()
	{
		router.GET("/follows/:user1/:user2", server.GetFollow)
		router.GET("/follows/:user1", server.GetFollowings)
		router.GET("/followers/:user2", server.GetFollowers)
		router.GET("/follows", server.GetAll)

		router.POST("/follows/one", server.CreateFollow)
		router.POST("/follows/batch", server.CreateFollowBatch)
		router.DELETE("/follows/:user1/:user2", server.DeleteFollow)
	}

	// start the router
	router.Run(":9080")
}
