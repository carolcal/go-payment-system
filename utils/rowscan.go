package utils

import (
	"database/sql"
	"qr-payment/models"
)

func ScanRow(row *sql.Row, payment *models.PaymentData) error {
	return row.Scan(
		&payment.ID,
		&payment.Amount,
		&payment.Status,
		&payment.CreatedAt,
		&payment.ExpiresAt,
		&payment.QRCodeData,
	)
}

func ScanRows(rows *sql.Rows, payment *models.PaymentData) error {
	return rows.Scan(
		&payment.ID,
		&payment.Amount,
		&payment.Status,
		&payment.CreatedAt,
		&payment.ExpiresAt,
		&payment.QRCodeData,
	)
}
