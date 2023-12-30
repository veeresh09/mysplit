// models/userController.go

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty"`
	Name     string               `bson:"name,omitempty"`
	Email    string               `bson:"email,omitempty"`
	Password string               `bson:"password,omitempty"` // Store hashed passwords only
	Groups   []primitive.ObjectID `bson:"groups,omitempty"`   // References to Group IDs
}
