package main

import (
	"context"
	"fmt"

	"github.com/andyVB2012/referral-service/http"
	"github.com/andyVB2012/referral-service/kafka"
	"github.com/andyVB2012/referral-service/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// create a database connection

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(
		"mongodb+srv://andrewvb2012:U3aQdGnoYzznB1d2@cluster0.gmjmvad.mongodb.net/",
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
	repository := repository.NewRepository(client.Database("stfx-referral-system"))

	// create an http server
	server := http.NewServer(repository)

	// create a kafka consumer
	go kafka.RunnConsumers(repository)
	// // create a gin router
	router := gin.Default()
	{
		router.GET("/referral-kafka-health", server.KafkaHealthCheck)
		router.GET("/referral-health", server.HealthCheck)
		router.GET("/referral-stats/:address", server.GetAttributionStats)
		router.GET("/referral-attributors/:address", server.GetAllAttributions)
		router.GET("/referral-code/:address", server.GetReferralCode)
		router.POST("/referral-generatecode/:address", server.CreateReferralCode)
		router.POST("/referral-addattributor/:code", server.AddAttributor)
	}

	// start the router
	router.Run(":9090")

}
