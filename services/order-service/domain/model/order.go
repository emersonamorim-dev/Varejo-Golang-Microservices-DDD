package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	CustomerID      string             `json:"customerId" bson:"customerId"`
	Products        []OrderProduct     `json:"products" bson:"products"`
	TotalPrice      float64            `json:"totalPrice" bson:"totalPrice"`
	ShippingAddress Address            `json:"shippingAddress" bson:"shippingAddress"`
	Status          OrderStatus        `json:"status" bson:"status"`
	OrderDate       time.Time          `json:"orderDate" bson:"orderDate"`
	DeliveryDate    time.Time          `json:"deliveryDate" bson:"deliveryDate"`
}

type OrderProduct struct {
	ProductID   string  `json:"productId" bson:"productId"`
	ProductName string  `json:"productName" bson:"productName"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	Price       float64 `json:"price" bson:"price"`
}

type Address struct {
	Street     string `json:"street" bson:"street"`
	City       string `json:"city" bson:"city"`
	State      string `json:"state" bson:"state"`
	PostalCode string `json:"postalCode" bson:"postalCode"`
	Country    string `json:"country" bson:"country"`
}

type OrderStatus string

const (
	Pending   OrderStatus = "PENDING"
	Shipped   OrderStatus = "SHIPPED"
	Delivered OrderStatus = "DELIVERED"
	Canceled  OrderStatus = "CANCELED"
)
