package service

import (
	"Varejo-Golang-Microservices/services/customer-service/domain/model"
	"Varejo-Golang-Microservices/services/customer-service/domain/repository"
)

type CustomerService interface {
	GetCustomerByID(id string) (*model.Customer, error)
	GetAllCustomers() ([]model.Customer, error)
	SaveCustomer(customer *model.Customer) error
	UpdateCustomer(customer *model.Customer) error
	DeleteCustomer(id string) error
}

type CustomerServiceImpl struct {
	customerRepo *repository.MongoCustomerRepository
}

func (s *CustomerServiceImpl) AddCustomer(customer *model.Customer) error {
	return s.SaveCustomer(customer)
}

func NewCustomerService(customerRepo *repository.MongoCustomerRepository) CustomerService {
	return &CustomerServiceImpl{
		customerRepo: customerRepo,
	}
}

func (s *CustomerServiceImpl) GetAllCustomers() ([]model.Customer, error) {
	return s.customerRepo.GetAll()
}

func (s *CustomerServiceImpl) GetCustomerByID(id string) (*model.Customer, error) {
	return s.customerRepo.FindByID(id)
}

func (s *CustomerServiceImpl) SaveCustomer(customer *model.Customer) error {
	return s.customerRepo.Save(customer)
}

func (s *CustomerServiceImpl) UpdateCustomer(customer *model.Customer) error {
	return s.customerRepo.Update(customer)
}

func (s *CustomerServiceImpl) DeleteCustomer(id string) error {
	return s.customerRepo.Delete(id)
}
