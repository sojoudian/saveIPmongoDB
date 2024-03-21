package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	http.HandleFunc("/saveIP", saveIPHandler)
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	log.Println("Server started on port 3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func saveIPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTION" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error parsing JSON request body", http.StatusBadRequest)
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	data["time"] = currentTime

	collection := mongoClient.Database("IPdatabase").Collection("IPs")
	_, err = collection.InsertOne(context.TODO(), data)
	if err != nil {
		http.Error(w, "Error saving data to MongoDB", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Data saved successfully")
}
