package dto

import (
	"Varejo-Golang-Microservices/services/promotion-service/domain/model"
	"time"
)

type PromotionDTO struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"` 
	Discount      float64           `json:"discount"` 
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	StartDate     time.Time         `json:"startDate"`
	EndDate       time.Time         `json:"endDate"`
	DiscountValue float64           `json:"discountValue"`
	ProductID     string            `json:"productId"`
	Status        model.PromoStatus `json:"status"`
	UpdatedAt     time.Time         `json:"updatedAt"`

}

// A estrutura ProductReferenceDTO pode ser usada se você quiser incluir 
// detalhes limitados do produto na resposta da promoção.
type ProductReferenceDTO struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
