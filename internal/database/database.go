package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var UsersColl *mongo.Collection
var ProductsColl *mongo.Collection
var InteractionsColl *mongo.Collection

func ConnectDB(uri string) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	UsersColl = client.Database("databaseproject").Collection("users")
	ProductsColl = client.Database("databaseproject").Collection("products")
	InteractionsColl = client.Database("databaseproject").Collection("interactions")
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}
