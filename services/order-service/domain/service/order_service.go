package service

import (
	"Varejo-Golang-Microservices/services/order-service/domain/model"
	"Varejo-Golang-Microservices/services/order-service/domain/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService interface {
	GetOrderByID(id string) (*model.Order, error)
	SaveOrder(order *model.Order) error
	GetAllOrders() ([]*model.Order, error)
	UpdateOrderStatus(id string, status model.OrderStatus) error
	DeleteOrder(id string) error
}

type OrderServiceImpl struct {
	orderRepo *repository.MongoOrderRepository
}

func NewOrderService(orderRepo *repository.MongoOrderRepository) OrderService {
	return &OrderServiceImpl{
		orderRepo: orderRepo,
	}
}

func (s *OrderServiceImpl) GetOrderByID(id string) (*model.Order, error) {
	return s.orderRepo.FindByID(id)
}

func (s *OrderServiceImpl) SaveOrder(order *model.Order) error {
	return s.orderRepo.Save(order)
}

func (s *OrderServiceImpl) GetAllOrders() ([]*model.Order, error) {
	return s.orderRepo.GetAll()
}

func (s *OrderServiceImpl) UpdateOrderStatus(id string, status model.OrderStatus) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	orderToUpdate := &model.Order{ID: objID, Status: status}
	return s.orderRepo.Update(orderToUpdate)
}

func (s *OrderServiceImpl) UpdateOrder(order *model.Order) error {
	return s.orderRepo.Update(order)
}

func (s *OrderServiceImpl) DeleteOrder(id string) error {
	return s.orderRepo.Delete(id)
}
