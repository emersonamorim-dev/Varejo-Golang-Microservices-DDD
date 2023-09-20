package service

import (
	"Varejo-Golang-Microservices/services/payment-service/domain/model"
	"Varejo-Golang-Microservices/services/payment-service/domain/repository"
)

type PaymentService interface {
	GetAllPayments() ([]*model.Payment, error)
	GetPaymentByID(id string) (*model.Payment, error)
	SavePayment(payment *model.Payment) error
	UpdatePayment(payment *model.Payment) error
	DeletePayment(id string) error
}

type PaymentServiceImpl struct {
	paymentRepo *repository.MongoPaymentRepository
}

func NewPaymentService(paymentRepo *repository.MongoPaymentRepository) PaymentService {
	return &PaymentServiceImpl{
		paymentRepo: paymentRepo,
	}
}

func (s *PaymentServiceImpl) GetAllPayments() ([]*model.Payment, error) {
	return s.paymentRepo.GetAllPayments()
}

func (s *PaymentServiceImpl) GetPaymentByID(id string) (*model.Payment, error) {
	return s.paymentRepo.FindByID(id)
}

func (s *PaymentServiceImpl) SavePayment(payment *model.Payment) error {
	return s.paymentRepo.Save(payment)
}

func (s *PaymentServiceImpl) UpdatePayment(payment *model.Payment) error {
	return s.paymentRepo.Update(payment)
}

func (s *PaymentServiceImpl) DeletePayment(id string) error {
	return s.paymentRepo.Delete(id)
}
