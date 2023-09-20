package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Promotion struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	Title         string             `json:"title" bson:"title"`
	Description   string             `json:"description" bson:"description"`
	StartDate     time.Time          `json:"startDate" bson:"startDate"`
	EndDate       time.Time          `json:"endDate" bson:"endDate"`
	Discount      float64            `json:"discount" bson:"discount"`
	DiscountValue float64            `json:"discountValue" bson:"discountValue"`
	Status        PromoStatus        `json:"status" bson:"status"`
}

type PromoStatus string

const (
	Active    PromoStatus = "ACTIVE"
	Expired   PromoStatus = "EXPIRED"
	Scheduled PromoStatus = "SCHEDULED"
)
