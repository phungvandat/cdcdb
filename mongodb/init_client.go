package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func initClient() (*mongo.Database, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	opts := options.Client()
	opts = opts.ApplyURI("mongodb://localhost:27017/?connect=direct")
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	log.Println("db connected")
	db := client.Database("cdc")
	return db, func() {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Printf("close db connection failed: %v\n", err)
		} else {
			log.Println("db disconnected")
		}
	}
}
