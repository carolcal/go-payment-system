package models

import (
	"time"
)

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "pending"
	StatusPaid      PaymentStatus = "paid"
	StatusFailed    PaymentStatus = "failed"
	StatusExpired   PaymentStatus = "expired"
)

type PaymentData struct {
	ID			string			`json:"id"`
	Amount		int				`json:"amount"`
	Status		PaymentStatus	`json:"status"`
	CreatedAt	time.Time		`json:"created_at"`
	ExpiresAt	time.Time		`json:"expires_at"`
	QRCodeData	string			`json:"qr_code_data"`
}

type CreatePaymentData struct {
	Amount		float64			`json:"amount"`
}

