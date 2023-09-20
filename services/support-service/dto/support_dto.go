package dto

import (
	"Varejo-Golang-Microservices/services/support-service/domain/model"
	"time"
)

type SupportDTO struct {
	ID          string              `json:"id"`
	Subject     string              `json:"subject"`
	Message     string              `json:"message"`
	CreatedDate time.Time           `json:"createdDate"`
	Data        string              `json:"data"`
	Response    string              `json:"response"`
	Status      model.SupportStatus `json:"status"`
	Category    SupportCategoryDTO  `json:"category"`
}

type SupportCategoryDTO struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
