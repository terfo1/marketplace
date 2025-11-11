package database

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	clientInstance      *mongo.Client
	clientInstanceError error
	once                sync.Once
)

var (
	UsersColl        *mongo.Collection
	ProductsColl     *mongo.Collection
	InteractionsColl *mongo.Collection
)

func ConnectDB(uri string) {
	once.Do(func() {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

		client, err := mongo.Connect(opts)
		if err != nil {
			clientInstanceError = fmt.Errorf("failed to connect to MongoDB: %v", err)
			return
		}

		if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
			clientInstanceError = fmt.Errorf("failed to ping MongoDB: %v", err)
			return
		}

		clientInstance = client
		log.Println("Successfully connected to MongoDB!")

		// Initialize collections
		UsersColl = client.Database("databaseproject").Collection("users")
		ProductsColl = client.Database("databaseproject").Collection("products")
		InteractionsColl = client.Database("databaseproject").Collection("interactions")
	})

	if clientInstanceError != nil {
		log.Fatal(clientInstanceError)
	}
}

// GetClient returns the MongoDB client (optional, if needed elsewhere).
func GetClient() *mongo.Client {
	return clientInstance
}

// CloseDB disconnects the MongoDB client (call on program shutdown).
func CloseDB() {
	if clientInstance != nil {
		if err := clientInstance.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting MongoDB: %v\n", err)
		}
		log.Println("Disconnected from MongoDB")
	}
}
