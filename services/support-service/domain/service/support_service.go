package service

import (
	"Varejo-Golang-Microservices/services/support-service/domain/model"
	"Varejo-Golang-Microservices/services/support-service/domain/repository"
)

type SupportService interface {
	GetSupportByID(id string) (*model.Support, error)
	SaveSupport(support *model.Support) error
	UpdateSupport(support *model.Support) error
	DeleteSupport(id string) error
	ListAllSupports() ([]*model.Support, error)
}

type SupportServiceImpl struct {
	supportRepo *repository.MongoSupportRepository
}

func NewSupportService(supportRepo *repository.MongoSupportRepository) SupportService {
	return &SupportServiceImpl{
		supportRepo: supportRepo,
	}
}

func (s *SupportServiceImpl) GetSupportByID(id string) (*model.Support, error) {
	return s.supportRepo.FindByID(id)
}

func (s *SupportServiceImpl) SaveSupport(support *model.Support) error {
	return s.supportRepo.Save(support)
}

func (s *SupportServiceImpl) UpdateSupport(support *model.Support) error {
    return s.supportRepo.Update(support)
}

func (s *SupportServiceImpl) DeleteSupport(id string) error {
	return s.supportRepo.Delete(id)
}

func (s *SupportServiceImpl) ListAllSupports() ([]*model.Support, error) {
	return s.supportRepo.ListAll()
}
