package dto

import (
	"Varejo-Golang-Microservices/services/location-service/domain/model"
	"time"
)

type LocationDTO struct {
	ID          string                 `json:"id,omitempty"`
	Address     string                 `json:"address,omitempty"`
	City        string                 `json:"city,omitempty"`
	State       string                 `json:"state,omitempty"`
	Country     string                 `json:"country,omitempty"`
	PostalCode  string                 `json:"postalCode,omitempty"`
	Description string                 `json:"description,omitempty"`
	Data        string                 `json:"data,omitempty"`
	CreatedDate time.Time              `json:"createdDate"`
	Coordinates LocationCoordinatesDTO `json:"coordinates"`
	Latitude    float64                `json:"latitude"`
	Longitude   float64                `json:"longitude"`
	Status      model.LocationStatus   `json:"status"`
}

type LocationCoordinatesDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
