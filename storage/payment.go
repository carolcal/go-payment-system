package storage

import (
	"fmt"
	"time"
	"database/sql"

	"qr-payment/models"
	"qr-payment/qrcode"
	"qr-payment/utils"
)

func GetPaymentById(id string, db *sql.DB) (*models.PaymentData, error) {
	var payment models.PaymentData

	row := db.QueryRow(`SELECT * FROM payments WHERE id=?`, id)
	err := utils.ScanPaymentRow(row, &payment)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pagamento não encontrado")
		}
		return nil, fmt.Errorf("erro ao escanear pagamento: %w", err)
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
		err := utils.ScanPaymentRows(rows, &payment)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear pagamento: %w", err)
		}
		allPayments[payment.ID] = &payment
	}

	return allPayments, nil
}

func GetAllPaymentsByUserId(user_type models.TypeUser, user_id string, db *sql.DB) (map[string]*models.PaymentData, error) {
	allPayments := make(map[string]*models.PaymentData)

	query := fmt.Sprintf(`SELECT * FROM payments WHERE %s=? `, user_type)
	rows, err := db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment models.PaymentData
		err := utils.ScanPaymentRows(rows, &payment)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear pagamento: %w", err)
		}
		allPayments[payment.ID] = &payment
	}

	return allPayments, nil
}

func CreatePayment(p *models.PaymentData, db *sql.DB) error {
	id := utils.GenerateID("pay")
	p.ID = id
	p.CreatedAt = time.Now()
	p.ExpiresAt = time.Now().Add(15 * time.Minute)
	p.Status = models.StatusPending
	p.QRCodeData = qrcode.GenerateQRCode("92991514078", p.Amount, "Arthur Dent", "Terra")

	_, err := db.Exec("INSERT INTO payments VALUES(?, ?, ?, ?, ?, ?, ?, ?);", p.ID, p.CreatedAt, p.ExpiresAt, p.Amount, p.Status, p.ReceiverId, "", p.QRCodeData)
	if err != nil {
		return fmt.Errorf("falha ao criar novo pagamento")
	}

	return nil
}

func MakePayment(id string, db *sql.DB) error {
	var payment models.PaymentData

	row := db.QueryRow(`SELECT * FROM payments WHERE id=?`, id)
	err := utils.ScanPaymentRow(row, &payment)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("pagamento não encontrado")
		}
		return fmt.Errorf("erro ao escanear pagamento: %w", err)
	}

	switch payment.Status {
	case models.StatusPaid:
		return fmt.Errorf("pagamento já realizado")
	case models.StatusFailed:
		return fmt.Errorf("pagamento falhou anteriormente")
	case models.StatusExpired:
		return fmt.Errorf("pagamento expirado")
	}

	if time.Now().After(payment.ExpiresAt) {
		_, err := db.Exec("UPDATE payments SET status=? WHERE id=?", models.StatusExpired, id)
		if err != nil {
			return fmt.Errorf("falha ao atualizar pagamento para status expirado: %w", err)
		}
		return fmt.Errorf("pagamento expirou")
	}

	_, err = db.Exec("UPDATE payments SET status=? WHERE id=?", models.StatusPaid, id)
	if err != nil {
		return fmt.Errorf("falha ao atualizar pagamento para status paga: %w", err)
	}

	return nil
}

func RemovePayment(id string, db *sql.DB) error {
	var payment models.PaymentData

	row := db.QueryRow(`SELECT * FROM payments WHERE id=?`, id)
	err := utils.ScanPaymentRow(row, &payment)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("pagamento não encontrado")
		}
		return fmt.Errorf("erro ao escanear pagamento: %w", err)
	}

	_, err = db.Exec("DELETE FROM payments WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("falha ao deletar pagamento: %w", err)
	}

	return nil
}
