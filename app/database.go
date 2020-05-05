package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB function to verify database connectivity
func ConnectDB() *mongo.Client {
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
		fmt.Println("-------------------------------")
		log.Fatal("Connection to database - FAILED", err)
		fmt.Println("-------------------------------")
	} else {
		fmt.Println("-------------------------------")
		fmt.Println("Connection to database - SUCCESS")
		fmt.Println("-------------------------------")
	}
	return client
}
