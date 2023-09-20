package model

import (
	"time"
)

type Location struct {
	ID string `json:"id" bson:"_id"`

	Latitude    float64        `json:"latitude" bson:"latitude"`
	Longitude   float64        `json:"longitude" bson:"longitude"`
	Description string         `json:"description" bson:"description"`
	Address     string         `json:"address" bson:"address"`
	Data        string         `json:"data" bson:"data"`
	CreatedDate time.Time      `json:"createdDate" bson:"createdDate"`
	Status      LocationStatus `json:"status" bson:"status"`
}

type LocationStatus string

const (
	Active   LocationStatus = "ACTIVE"
	Inactive LocationStatus = "INACTIVE"
	Pending  LocationStatus = "PENDING"
)
