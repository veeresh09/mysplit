// models/group.go

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	ID      primitive.ObjectID   `bson:"_id,omitempty"`
	Name    string               `bson:"name,omitempty"`
	Members []primitive.ObjectID `bson:"members,omitempty"` // References to User IDs
}
