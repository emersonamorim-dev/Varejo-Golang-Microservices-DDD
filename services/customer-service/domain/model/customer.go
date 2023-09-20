package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Customer representa um cliente no sistema.
type Customer struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Email   string             `json:"email" bson:"email"`
	Cell    string             `json:"cell" bson:"cell"`
	Phone   string             `json:"phone" bson:"phone"`
	Address string             `json:"address" bson:"address"`
	ZipCode string             `json:"zipCode" bson:"zipCode"`
	City    string             `json:"city" bson:"city"`
}
