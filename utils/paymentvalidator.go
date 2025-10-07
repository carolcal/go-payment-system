package utils

import (
	"fmt"
	"strconv"
	"time"

	"database/sql"

	"qr-payment/models"
)

func ValidateBalance(user_id string, qrdata *models.QRCodeData, db *sql.DB) (int, int, error) {
	amount, err := strconv.ParseFloat(qrdata.TransactionAmount, 64)
	if err != nil || amount <= 0 {
		return 0, 0, fmt.Errorf("invalid transaction amount")
	}

	var balance int
	err = db.QueryRow(`SELECT balance FROM users WHERE id=? LIMIT 1;`, user_id).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, 0, fmt.Errorf("user not found")
	}
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch balance: %w", err)
	}

	new_balance := balance - int(amount*100)
	if new_balance < 0 {
		return 0, 0, fmt.Errorf("not enough balance to make this transaction")
	}

	return new_balance, int(amount*100), nil
}

func ValidatePixKey(qrdata *models.QRCodeData, db *sql.DB) (string, int, error) {
	var receiver_id string
	var balance int
	err := db.QueryRow(`SELECT id, balance FROM users WHERE name=? AND cpf=? LIMIT 1;`, qrdata.MerchantName, qrdata.PixKey).Scan(&receiver_id, &balance)
	if err == sql.ErrNoRows {
		return "", 0, fmt.Errorf("user not found")
	}
	if err != nil {
		return "", 0, fmt.Errorf("failed to fetch id and balance: %w", err)
	}
	return receiver_id, balance, nil
}

func ValidatePaymentStatus(qr_code_data string, db *sql.DB) (*models.PaymentData, error) {
	payment := models.PaymentData{}
	row := db.QueryRow(`SELECT * FROM payments WHERE qr_code_data=?`, qr_code_data)
	err := ScanPaymentRow(row, &payment)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pagamento não encontrado")
		}
		return nil, fmt.Errorf("erro ao escanear pagamento: %w", err)
	}

	switch payment.Status {
	case models.StatusPaid:
		return nil, fmt.Errorf("pagamento já realizado")
	case models.StatusFailed:
		return nil, fmt.Errorf("pagamento falhou anteriormente")
	case models.StatusExpired:
		return nil, fmt.Errorf("pagamento expirado")
	}

	if time.Now().After(payment.ExpiresAt) {
		_, err := db.Exec("UPDATE payments SET status=? WHERE id=?", models.StatusExpired, payment.ID)
		if err != nil {
			return nil, fmt.Errorf("falha ao atualizar pagamento para status expirado: %w", err)
		}
		return nil, fmt.Errorf("pagamento expirou")
	}

	return &payment, nil

}
