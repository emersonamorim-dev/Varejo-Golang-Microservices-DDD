package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Category    Category           `json:"category" bson:"category"`
	Stock       int                `json:"stock" bson:"stock"`
	AddedDate   time.Time          `json:"addedDate" bson:"addedDate"`
	Status      ProductStatus      `json:"status" bson:"status"`
}

type Category struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

type ProductStatus string

const (
	Available    ProductStatus = "AVAILABLE"
	OutOfStock   ProductStatus = "OUT_OF_STOCK"
	Discontinued ProductStatus = "DISCONTINUED"
)
