package main

import (
	"context"
	"fmt"
	"time"

	// "github.com/andyVB2012/referral-service/proto"

	pb "github.com/andyVB2012/referral-service/block"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func main() {
	// create a database connection

	// // Use the SetServerAPIOptions() method to set the Stable API version to 1
	// serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// // password := "77J8Knd6vrhvVGjL"
	// opts := options.Client().ApplyURI(
	// 	"mongodb+srv://user insert here",
	// ).SetServerAPIOptions(serverAPI)

	// // Create a new client and connect to the server
	// client, err := mongo.Connect(context.TODO(), opts)
	// if err != nil {
	// 	panic(err)
	// }

	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// // Send a ping to confirm a successful connection
	// if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	// // create a repository
	// repository := repository.NewRepository(client.Database("stfx-referral"))

	// // create an http server
	// server := http.NewServer(repository)

	// // create a gin router
	// router := gin.Default()
	// {
	// 	router.GET("/follows/:user1/:user2", server.)
	// 	router.GET("/follows/:user1", server.GetFollowings)
	// 	router.GET("/followers/:user2", server.GetFollowers)
	// 	router.GET("/follows", server.GetAll)

	// 	router.POST("/follows/one", server.CreateFollow)
	// 	router.POST("/follows/batch", server.CreateFollowBatch)
	// 	router.DELETE("/follows/:user1/:user2", server.DeleteFollow)
	// }

	// // start the router
	// router.Run(":9080")

	// type Consumer struct {
	// 	reader *kafka.Reader
	// }

	// make a new reader that consumes from topic-A, partition 0, at offset 42
	// and only reads new messages
	// bb,_ := kafka.NewConsumerGroup(
	// 	kafka.ConsumerGroupConfig{
	// 		ID:          "my-consumer",
	// 		Brokers:     []string{"localhost:9092"},
	// 		Topics:      []string{"stfx.stream.vault"},
	// 		StartOffset: kafka.LastOffset,
	// 	},
	// )
	// bb.Next(context.Background())

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "stfx.stream.block.arb.stfx.vault.v2", 0)
	if err != nil {
		fmt.Println(err)
	}
	conn.SetReadDeadline(time.Now().Add(100 * time.Second))

	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
	defer batch.Close()
	count := 0

	for {
		b := make([]byte, 10e3) // 10KB max per message
		n, err := batch.Read(b)
		if err != nil {
			break // end of batch
		}
		b = b[:n] // truncate the buffer to the actual message size

		var blockEvent pb.BlockEvent
		if err := proto.Unmarshal(b, &blockEvent); err != nil {
			fmt.Printf("error unmarshalling message: %v", err)
			continue
		}

		fmt.Printf("Received message: %+v\n", blockEvent)
		count++
	}
	fmt.Println("count", count)
}
