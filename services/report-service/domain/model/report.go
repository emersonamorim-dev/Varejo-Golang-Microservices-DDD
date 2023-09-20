package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	CreatedDate time.Time          `json:"createdDate" bson:"createdDate"`
	Data        string             `json:"data" bson:"data"`
	Status      ReportStatus       `json:"status" bson:"status"`
	Category    Category           `json:"category" bson:"category"`
}

type Category struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

type ReportStatus string

const (
	Published ReportStatus = "PUBLISHED"
	Draft     ReportStatus = "DRAFT"
	Archived  ReportStatus = "ARCHIVED"
)
