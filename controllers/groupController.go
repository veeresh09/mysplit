package controllers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"mysplit/models"
)

var GroupsCollection *mongo.Collection

// CreateGroup creates a new group
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var newGroup models.Group
	if err := json.NewDecoder(r.Body).Decode(&newGroup); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newGroup.ID = primitive.NewObjectID()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := GroupsCollection.InsertOne(ctx, newGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newGroup)
}

type groupAddRequest struct {
	groupID string
	userID  string
}

// AddUserToGroup adds a user to an existing group
func AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	var newReq groupAddRequest

	if err := json.NewDecoder(r.Body).Decode(&newReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := GroupsCollection.UpdateOne(
		ctx,
		bson.M{"_id": newReq.groupID}, // groupID should be the ObjectID of the group
		bson.M{"$push": bson.M{"members": newReq.userID}}, // userID should be the ObjectID of the user to add
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User added to group"))
}

// ListGroupMembers lists all members of a group
func ListGroupMembers(w http.ResponseWriter, r *http.Request) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	groupID := r.URL.Query().Get("groupID")
	if groupID == "" {
		http.Error(w, "empty groupID", http.StatusBadRequest)
		return
	}
	// Find the group by ID
	var group models.Group
	err := GroupsCollection.FindOne(ctx, bson.M{"_id": groupID}).Decode(&group) // groupID should be the ObjectID of the group
	if err != nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	// Fetch user details for each member
	var members []models.User
	for _, memberID := range group.Members {
		var user models.User
		err := UsersCollection.FindOne(ctx, bson.M{"_id": memberID}).Decode(&user)
		if err == nil {
			members = append(members, user)
		}
		// You might want to handle the error in case a user is not found
	}

	// Send the list of members as a response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(members)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
