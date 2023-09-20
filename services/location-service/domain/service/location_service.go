package service

import (
	"Varejo-Golang-Microservices/services/location-service/domain/model"
	"Varejo-Golang-Microservices/services/location-service/domain/repository"
)

type LocationService interface {
	GetLocationByID(id string) (*model.Location, error)
	SaveLocation(location *model.Location) error
	UpdateLocation(location *model.Location) error
	DeleteLocation(id string) error
	GetAllLocations() ([]*model.Location, error)
}

type LocationServiceImpl struct {
	locationRepo *repository.MongoLocationRepository
}

func NewLocationService(locationRepo *repository.MongoLocationRepository) LocationService {
	return &LocationServiceImpl{
		locationRepo: locationRepo,
	}
}

func (s *LocationServiceImpl) GetAllLocations() ([]*model.Location, error) {
	return s.locationRepo.ListAll()
}

func (s *LocationServiceImpl) GetLocationByID(id string) (*model.Location, error) {
	return s.locationRepo.FindByID(id)
}

func (s *LocationServiceImpl) SaveLocation(location *model.Location) error {
	return s.locationRepo.Save(location)
}

func (s *LocationServiceImpl) UpdateLocation(location *model.Location) error {
	return s.locationRepo.Update(location)
}

func (s *LocationServiceImpl) DeleteLocation(id string) error {
	return s.locationRepo.Delete(id)
}
