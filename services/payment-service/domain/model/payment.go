package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	OrderID     string             `json:"orderId" bson:"orderId"`
	CustomerID  string             `json:"customerId" bson:"customerId"`
	Amount      float64            `json:"amount" bson:"amount"`
	Method      PaymentMethod      `json:"method" bson:"method"`
	Status      PaymentStatus      `json:"status" bson:"status"`
	PaymentDate time.Time          `json:"paymentDate" bson:"paymentDate"`
}

type PaymentMethod struct {
	Type       PaymentType `json:"type" bson:"type"`
	CardNumber string      `json:"cardNumber,omitempty" bson:"cardNumber,omitempty"`
	Expiry     string      `json:"expiry,omitempty" bson:"expiry,omitempty"`
	CVV        string      `json:"cvv,omitempty" bson:"cvv,omitempty"`
}

type PaymentType string

const (
	CreditCard PaymentType = "CREDIT_CARD"
	DebitCard  PaymentType = "DEBIT_CARD"
	PayPal     PaymentType = "PAYPAL"
)

type PaymentStatus string

const (
	Unpaid    PaymentStatus = "UNPAID"
	Processed PaymentStatus = "PROCESSED"
	Failed    PaymentStatus = "FAILED"
	Refunded  PaymentStatus = "REFUNDED"
)
