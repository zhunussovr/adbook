package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB function to verify database connectivity
func ConnectDB() (*mongo.Client, error) {
	// For using SCRAM-SHA-1 auth.mechanism
	dbcreds := options.Credential{
		Username: "adbookadm",
		Password: "adbookadm",
	}
	clientOptions := options.Client().ApplyURI("mongodb://localhost/adbook").SetAuth(dbcreds)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	// Ferify connections
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Connection to database - FAILED:")
		log.Fatal(err)
	} else {
		fmt.Println("-------------------------------")
		fmt.Println("Connection to database - SUCCESS")
		fmt.Println("-------------------------------")
	}
	return client, nil
}

// GetCollections - to retrieve collections
func GetCollections(DbName string, CollectionName string) (*mongo.Collection, error) {
	client, err := ConnectDB()

	if err != nil {
		return nil, err
	}

	collection := client.Database(DbName).Collection(CollectionName)

	return collection, nil
}
