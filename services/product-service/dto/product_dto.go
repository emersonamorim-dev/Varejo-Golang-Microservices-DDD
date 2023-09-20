package dto

import (
	"Varejo-Golang-Microservices/services/product-service/domain/model"
	"time"
)

type ProductDTO struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Price       float64             `json:"price"`
	Category    CategoryDTO         `json:"category"`
	Stock       int                 `json:"stock"`
	AddedDate   time.Time           `json:"addedDate"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Status      model.ProductStatus `json:"status"`
}

type CategoryDTO struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
