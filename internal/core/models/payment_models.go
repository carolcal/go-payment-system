package models

import (
	"time"
)

type PaymentStatus string
type TypeUser string

const (
	StatusPending PaymentStatus = "pending"
	StatusPaid    PaymentStatus = "paid"
	StatusFailed  PaymentStatus = "failed"
	StatusExpired PaymentStatus = "expired"
)

const (
	UserReceiver TypeUser = "receiver_id"
	UserPayer    TypeUser = "payer_id"
)

func IsValidTypeUser(userType string) (TypeUser, bool) {
	switch userType {
	case string(UserReceiver):
		return UserReceiver, true
	case string(UserPayer):
		return UserPayer, true
	default:
		return "", false
	}
}

type PaymentData struct {
	ID         string        `json:"id"`
	CreatedAt  time.Time     `json:"created_at"`
	ExpiresAt  time.Time     `json:"expires_at"`
	Amount     int           `json:"amount"`
	Status     PaymentStatus `json:"status"`
	ReceiverId string        `json:"receiver_id"`
	PayerId    string        `json:"payer_id"`
	QRCodeData string        `json:"qr_code_data"`
}

type CreatePaymentData struct {
	Amount     *float64	`json:"amount,omitempty"`
	ReceiverId string	`json:"receiver_id"`
}

type ProcessPaymentData struct {
	QRCodeData string	`json:"qr_code_data"`
	Amount     *float64	`json:"amount,omitempty"`
}

