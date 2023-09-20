package service

import (
	"Varejo-Golang-Microservices/services/product-service/domain/model"
	"Varejo-Golang-Microservices/services/product-service/domain/repository"
)

type ProductService interface {
	GetProductByID(id string) (*model.Product, error)
	SaveProduct(product *model.Product) error
	UpdateProduct(product *model.Product) error
	DeleteProduct(id string) error
	ListAllProducts() ([]*model.Product, error)
}

type ProductServiceImpl struct {
	productRepo *repository.MongoProductRepository
}

func NewProductService(productRepo *repository.MongoProductRepository) ProductService {
	return &ProductServiceImpl{
		productRepo: productRepo,
	}
}

func (s *ProductServiceImpl) GetProductByID(id string) (*model.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *ProductServiceImpl) SaveProduct(product *model.Product) error {
	return s.productRepo.SaveProduct(product)
}

func (s *ProductServiceImpl) UpdateProduct(product *model.Product) error {
	return s.productRepo.Update(product)
}

func (s *ProductServiceImpl) DeleteProduct(id string) error {
	return s.productRepo.Delete(id)
}

func (s *ProductServiceImpl) ListAllProducts() ([]*model.Product, error) {
	return s.productRepo.ListAll()
}
