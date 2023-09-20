package service

import (
	"Varejo-Golang-Microservices/services/integration-service/domain/model"
	"Varejo-Golang-Microservices/services/integration-service/domain/repository"
)

type IntegrationService interface {
	GetIntegrationDataByID(id string) (*model.IntegrationData, error)
	SaveIntegrationData(data *model.IntegrationData) error
	UpdateIntegrationData(data *model.IntegrationData) error
	DeleteIntegrationData(id string) error
	ListAllIntegrationData() ([]*model.IntegrationData, error)
}

type IntegrationServiceImpl struct {
	integrationRepo *repository.MongoIntegrationRepository
}

func NewIntegrationService(integrationRepo *repository.MongoIntegrationRepository) IntegrationService {
	return &IntegrationServiceImpl{
		integrationRepo: integrationRepo,
	}
}

func (s *IntegrationServiceImpl) ListAllIntegrationData() ([]*model.IntegrationData, error) {
	return s.integrationRepo.ListAllIntegrationData()
}

func (s *IntegrationServiceImpl) GetAllIntegrationData() ([]*model.IntegrationData, error) {
	return s.integrationRepo.ListAllIntegrationData()
}

func (s *IntegrationServiceImpl) GetIntegrationDataByID(id string) (*model.IntegrationData, error) {
	return s.integrationRepo.FindByID(id)
}

func (s *IntegrationServiceImpl) SaveIntegrationData(data *model.IntegrationData) error {
	return s.integrationRepo.Save(data)
}

func (s *IntegrationServiceImpl) UpdateIntegrationData(data *model.IntegrationData) error {
	return s.integrationRepo.Update(data)
}

func (s *IntegrationServiceImpl) DeleteIntegrationData(id string) error {
	return s.integrationRepo.DeleteIntegrationData(id)
}
