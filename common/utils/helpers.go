package utils

import (
	"Varejo-Golang-Microservices/services/promotion-service/domain/model"
	"time"
)

// IsPromotionActive verifica se uma promoção está atualmente ativa
func IsPromotionActive(promotion *model.Promotion) bool {
	currentTime := time.Now()
	return promotion.StartDate.Before(currentTime) && promotion.EndDate.After(currentTime)
}
