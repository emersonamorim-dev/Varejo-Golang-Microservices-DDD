package service

import (
	"Varejo-Golang-Microservices/services/promotion-service/domain/model"
	"Varejo-Golang-Microservices/services/promotion-service/domain/repository"
)

type PromotionService interface {
	ListAllPromotions() ([]*model.Promotion, error)
	GetPromotionByID(id string) (*model.Promotion, error)
	SavePromotion(promotion *model.Promotion) error
	UpdatePromotion(promotion *model.Promotion) error
	DeletePromotion(id string) error
}

type PromotionServiceImpl struct {
	promotionRepo *repository.MongoPromotionRepository
}

func NewPromotionService(promotionRepo *repository.MongoPromotionRepository) PromotionService {
	return &PromotionServiceImpl{
		promotionRepo: promotionRepo,
	}
}

func (s *PromotionServiceImpl) ListAllPromotions() ([]*model.Promotion, error) {
	return s.promotionRepo.ListAll()
}

func (s *PromotionServiceImpl) GetPromotionByID(id string) (*model.Promotion, error) {
	return s.promotionRepo.FindByID(id)
}

func (s *PromotionServiceImpl) SavePromotion(promotion *model.Promotion) error {
	return s.promotionRepo.Save(promotion)
}

func (s *PromotionServiceImpl) UpdatePromotion(promotion *model.Promotion) error {
	return s.promotionRepo.Update(promotion)
}

func (s *PromotionServiceImpl) DeletePromotion(id string) error {
	return s.promotionRepo.Delete(id)
}
