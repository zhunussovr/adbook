package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Cache struct {
	conn   *mongo.Client
	config *CacheConfig
}

type CacheConfig struct {
	Username       string
	Password       string
	DbName         string
	CollectionName string
	Port           int
}

func NewCache(config *CacheConfig) (*Cache, error) {
	conn, err := ConnectDB(config)
	if err != nil {
		return nil, err
	}

	return &Cache{conn, config}, nil

}

// ConnectDB function to verify database connectivity
func ConnectDB(config *CacheConfig) (*mongo.Client, error) {
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
func (cache *Cache) GetCollections() (*mongo.Collection, error) {

	collection := cache.conn.Database(cache.config.DbName).Collection(cache.config.CollectionName)

	return collection, nil
}

func (cache *Cache) GetEmployees(username string) ([]Employee, error) {
	var employees []Employee

	collection, err := cache.GetCollections()
	if err != nil {
		return nil, err
	}

	var filter bson.M = bson.M{}
	objID, _ := primitive.ObjectIDFromHex(username)
	filter = bson.M{"_id": objID}

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	cur.All(context.Background(), &results)

	// TODO: unmarshaling to employees from bson

	return employees, nil
}
