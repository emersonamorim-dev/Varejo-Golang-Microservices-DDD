package dto

import (
	"Varejo-Golang-Microservices/services/payment-service/domain/model"
	"time"
)

type PaymentDTO struct {
	ID        string             `json:"id"`
	OrderID   string             `json:"orderId"`
	CustomerID string            `json:"customerId"`
	Amount    float64            `json:"amount"`
	Method    PaymentMethodDTO   `json:"method"`
	Status    model.PaymentStatus `json:"status"`
	PaymentDate time.Time         `json:"paymentDate"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

type PaymentMethodDTO struct {
	Type        model.PaymentType `json:"type"`
	CardNumber  string            `json:"cardNumber,omitempty"` 
	Expiry      string            `json:"expiry,omitempty"`    
	CVV         string            `json:"cvv,omitempty"`        
}
