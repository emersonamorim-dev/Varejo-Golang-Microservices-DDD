package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IntegrationData representa as informações de uma integração em particular.
type IntegrationData struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`        
	Endpoint string             `bson:"endpoint"`   
	APIKey   string             `bson:"api_key"`   
	Data        string    `json:"data" bson:"data"`
	Other    map[string]string  `bson:"other"`       
}



