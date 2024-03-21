package main

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func main() {
	var err error
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	log.Println("Server started on port 3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
