package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database

func InitClient(uri string, dbName string) {
	clientOptions := options.Client().ApplyURI(uri)

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("failed to create client: %s\n", err)
	}

	err = c.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("failed to connect to database: %s\n", err)
	}

	log.Println("connected to MongoDB!")
	database = c.Database(dbName)
}

func RegisterCollections() {
	InitOrderRepo(database)
}
