package dto

import (
	"Varejo-Golang-Microservices/services/report-service/domain/model"
	"time"
)

type ReportDTO struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	CreatedDate time.Time          `json:"createdDate"`
	Data        string             `json:"data"`
	Status      model.ReportStatus `json:"status"`
	Category    ReportCategoryDTO  `json:"category"`
}

type Category struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

type ReportCategoryDTO struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
