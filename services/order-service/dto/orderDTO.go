package dto

import (
	"Varejo-Golang-Microservices/services/order-service/domain/model"
	"time"
)

type OrderDTO struct {
	ID              string            `json:"id"`
	CustomerID      string            `json:"customerId"`
	Products        []OrderItemDTO    `json:"products"`
	TotalPrice      float64           `json:"totalPrice"`
	ShippingAddress Address           `json:"shippingAddress" bson:"shippingAddress"`
	Status          model.OrderStatus `json:"status"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}

type OrderItemDTO struct {
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"Price"`
	Total       float64 `json:"total"`
}

type Address struct {
	Street     string `json:"street" bson:"street"`
	City       string `json:"city" bson:"city"`
	State      string `json:"state" bson:"state"`
	PostalCode string `json:"postalCode" bson:"postalCode"`
	Country    string `json:"country" bson:"country"`
}
