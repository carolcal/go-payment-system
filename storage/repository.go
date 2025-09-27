package storage

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"qr-payment/models"
	"qr-payment/utils"
)


func GetPaymentById(id string, db *sql.DB) (*models.PaymentData, error) {
	var payment models.PaymentData

	row := db.QueryRow(`SELECT * FROM payments WHERE id=?`, id)
	err := utils.ScanRow(row, &payment)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("error scaning payment: %w", err)
	}

	return &payment, nil
}

func GetAllPayments(db *sql.DB) (map[string]*models.PaymentData, error) {
	allPayments := make(map[string]*models.PaymentData)

	rows, err := db.Query(`SELECT * FROM payments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment models.PaymentData
		err := utils.ScanRows(rows, &payment)
		if err != nil {
			return nil, fmt.Errorf("error scaning payment: %w", err)
		}
		allPayments[payment.ID] = &payment
	}

	return allPayments, nil
}

func CreatePayment(p *models.PaymentData, db *sql.DB) error {
	id := utils.GenerateID()
	p.ID = id
	p.CreatedAt = time.Now()
	p.ExpiresAt = time.Now().Add(15 * time.Minute)
	p.Status = models.StatusPending
	qrcodeBytes, err := utils.GenerateQRCode(id)
	if err != nil {
		return err
	}
	p.QRCodeData = "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrcodeBytes)

	_, err = db.Exec("INSERT INTO payments VALUES(?, ?, ?, ?, ?, ?);", p.ID, p.Amount, p.Status, p.CreatedAt, p.ExpiresAt, p.QRCodeData)
	if (err != nil) {
		return fmt.Errorf("failed to create new payment")
	}

	return nil
}

func MakePayment(id string, db *sql.DB) error {
	var payment models.PaymentData

	row := db.QueryRow(`SELECT * FROM payments WHERE id=?`, id)
	err := utils.ScanRow(row, &payment)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("payment not found")
		}
		return fmt.Errorf("error scanning payment: %w", err)
	}

	switch payment.Status {
		case models.StatusPaid:
			return fmt.Errorf("payment already completed")
		case models.StatusFailed:
			return fmt.Errorf("payment failed previously")
		case models.StatusExpired:
			return fmt.Errorf("payment has expired")
	}

	if time.Now().After(payment.ExpiresAt) {
		_, err := db.Exec("UPDATE payments SET status=? WHERE id=?", models.StatusExpired, id)
		if err != nil {
			return fmt.Errorf("failed to update payment status to expired: %w", err)
		}
		return fmt.Errorf("payment has expired")
	}

	_, err = db.Exec("UPDATE payments SET status=? WHERE id=?", models.StatusPaid, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status to paid: %w", err)
	}

	return nil
}

func RemovePayment(id string, db *sql.DB) error {
	var payment models.PaymentData

	row := db.QueryRow(`SELECT * FROM payments WHERE id=?`, id)
	err := utils.ScanRow(row, &payment)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("payment not found")
		}
		return fmt.Errorf("error scanning payment: %w", err)
	}

	_, err = db.Exec("DELETE FROM payments WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}

	return nil
}
