package main

import (
	"context"
	"fmt"

	"github.com/andyVB2012/referral-service/http"
	"github.com/andyVB2012/referral-service/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/andyVB2012/referral-service/proto"
	// pb "github.com/andyVB2012/referral-service/block"
	// "github.com/andyVB2012/referral-service"
)

func main() {
	// create a database connection

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// password := "77J8Knd6vrhvVGjL"
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

	// res, err := repository.GetAllAttributors(context.Background(), "0x1")
	// if err != nil {
	// 	fmt.Println("error: ", err)
	// }
	// fmt.Println("Stats: ", res)

	// return
	// // repository.CreateReferralCode(context.Background(), "0x87ad83dc2f12a14c85d20f178a918a65edfe1b42")
	// refCode, err := repository.GetCode(context.Background(), "0x87ad83dc2f12a14c85d20f178a918a65edfe1b42")
	// if err != nil {
	// 	fmt.Println("error: ", err)
	// }
	// fmt.Println("Refcode: ", refCode)
	// // repository.AddAttributor(context.Background(), refCode, "0x87ad83dc2f12a14c85d20f178a918a65edfe1b42")

	// isIt := repository.IsTraderAddrInDb(context.Background(), "0x87ad83dc2f12a14c85d20f178a918a65edfe1b42")
	// fmt.Println("Hello", isIt)
	// kafka.RunnConsumers(repository)
	// // create a gin router
	router := gin.Default()
	{
		router.GET("/referral-stats/:address", server.GetAttributionStats)
		router.GET("/referral-attributors/:address", server.GetAllAttributions)
		router.GET("/referral-code/:address", server.GetReferralCode)
		// router.GET("/follows/:user1/:user2", server.CreateFollow)
		// router.GET("/follows/:user1", server.GetFollowings)
		// router.GET("/followers/:user2", server.GetFollowers)
		// router.GET("/follows", server.GetAll)

		// router.POST("/follows/one", server.CreateFollow)
		// router.POST("/follows/batch", server.CreateFollowBatch)
		// router.DELETE("/follows/:user1/:user2", server.DeleteFollow)
	}

	// start the router
	router.Run(":9090")

}
