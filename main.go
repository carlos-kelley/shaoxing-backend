package main

// Context handles concurrency stuff
// Need to import the driver stuff I got from go get in the module
import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set up MongoDB client options
	// Other options include auth, SSL, etc.
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB, and return a client whcih we can use to access db
	// This creates a root context. We can manage timeout across API boundaries
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Disconnect from MongoDB when finished
	// Go convention puts these defers right after the initialization
	defer func() {
		// Assign the return of the disconnect to err and check it
		// This is special Go syntax, initializing a var with an if
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Access database and collection
	collection := client.Database("shaoxing").Collection("words")

	// Define a filter to match all documents
	filter := bson.M{} // An empty filter matches all documents

	// Find documents that match the filter
	// Mongo returns a cursor, which is an iterator over the documents
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Iterate over the documents returned by the cursor
	for cursor.Next(context.Background()) {
		// Initialize a map with zero value nil
		// Why not empty??
		var word bson.M
		// Decode expects a pointer to a value, and only returns the error
		if err := cursor.Decode(&word); err != nil {
			log.Fatal(err)
		}
		fmt.Println(word)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
