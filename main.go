package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mysplit/controllers"
	"net/http"
	"time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello, world!")
}

func main() {
	http.HandleFunc("/", helloHandler)
	clientOptions := options.Client().ApplyURI(
		"mongodb://localhost:27017")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	controllers.UsersCollection = client.Database("mySplit").Collection("users")
	controllers.GroupsCollection = client.Database("mySplit").Collection("groups")

	http.HandleFunc("/create_user", controllers.CreateUser)
	http.HandleFunc("/create_group", controllers.CreateGroup)
	http.HandleFunc("/add_user_to_group", controllers.AddUserToGroup)
	http.HandleFunc("/list_group_members", controllers.ListGroupMembers)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
}
