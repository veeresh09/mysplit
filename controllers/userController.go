package controllers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mysplit/models"
	"net/http"
	"time"
)

// UsersCollection Assuming you have a global variable for MongoDB collection
var UsersCollection *mongo.Collection

// CreateUser creates a new user if it doesn't exist
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var existingUser models.User
	err := UsersCollection.FindOne(ctx, bson.M{"email": newUser.Email}).Decode(&existingUser)
	if err == nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	// Create new user
	newUser.ID = primitive.NewObjectID()
	_, err = UsersCollection.InsertOne(ctx, newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
