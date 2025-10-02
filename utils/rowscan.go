package utils

import (
	"database/sql"
	"qr-payment/models"
)

func ScanPaymentRow(row *sql.Row, payment *models.PaymentData) error {
	return row.Scan(
		&payment.ID,
		&payment.Amount,
		&payment.Status,
		&payment.CreatedAt,
		&payment.ExpiresAt,
		&payment.QRCodeData,
	)
}

func ScanPaymentRows(rows *sql.Rows, payment *models.PaymentData) error {
	return rows.Scan(
		&payment.ID,
		&payment.Amount,
		&payment.Status,
		&payment.CreatedAt,
		&payment.ExpiresAt,
		&payment.QRCodeData,
	)
}

func ScanUserRow(row *sql.Row, user *models.UserData) error {
	return row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Name,
		&user.CPF,
		&user.Balance,
	)
}

func ScanUserRows(rows *sql.Rows, user *models.UserData) error {
	return rows.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Name,
		&user.CPF,
		&user.Balance,
	)
}