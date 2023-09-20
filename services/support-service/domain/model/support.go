package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Support struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Subject     string             `json:"subject" bson:"subject"`
	Message     string             `json:"message" bson:"message"`
	CreatedDate time.Time          `json:"createdDate" bson:"createdDate"`
	Data        string             `json:"data" bson:"data"`
	Response    string             `json:"response" bson:"response"`
	Status      SupportStatus      `json:"status" bson:"status"`
	Category    Category           `json:"category" bson:"category"`
}

type Category struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

type SupportStatus string

const (
	Open     SupportStatus = "OPEN"
	Resolved SupportStatus = "RESOLVED"
	Closed   SupportStatus = "CLOSED"
	Pending  SupportStatus = "PENDING"
)
