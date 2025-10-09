package repository

import (
	"database/sql"
	"fmt"

	"qr-payment/internal/core/models"
	// "qr-payment/internal/qrcode"
	"qr-payment/internal/utils"
)

type PaymentRepository interface {
	FindById(id string) (*models.PaymentData, error)
	FindByQRCodeData(qr_code_data string) (*models.PaymentData, error)
	FindAll() (map[string]*models.PaymentData, error)
	FindAllByUserId(user_type models.TypeUser, user_id string) (map[string]*models.PaymentData, error)
	Create(u *models.UserData, p *models.PaymentData) error
	UpdatePaymentStatus(id string, status models.PaymentStatus) error
	Delete(id string) error
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}


func (p *paymentRepository) FindById(id string) (*models.PaymentData, error) {
	var payment models.PaymentData

	row := p.db.QueryRow(`SELECT * FROM payments WHERE id=?`, id)
	err := utils.ScanPaymentRow(row, &payment)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pagamento não encontrado")
		}
		return nil, fmt.Errorf("erro ao escanear pagamento: %w", err)
	}

	return &payment, nil
}

func (p *paymentRepository) FindByQRCodeData(qr_code_data string) (*models.PaymentData, error) {
	var payment models.PaymentData

	row := p.db.QueryRow(`SELECT * FROM payments WHERE qr_code_data=?`, qr_code_data)
	err := utils.ScanPaymentRow(row, &payment)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pagamento não encontrado")
		}
		return nil, fmt.Errorf("erro ao escanear pagamento: %w", err)
	}

	return &payment, nil
}


func (p *paymentRepository) FindAll() (map[string]*models.PaymentData, error) {
	allPayments := make(map[string]*models.PaymentData)

	rows, err := p.db.Query(`SELECT * FROM payments`)
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

func (p *paymentRepository) FindAllByUserId(user_type models.TypeUser, user_id string) (map[string]*models.PaymentData, error) {
	allPayments := make(map[string]*models.PaymentData)

	query := fmt.Sprintf(`SELECT * FROM payments WHERE %s=? `, user_type)
	rows, err := p.db.Query(query, user_id)
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

func (p *paymentRepository) Create(u *models.UserData, pd *models.PaymentData) error {
	_, err := p.db.Exec("INSERT INTO payments VALUES(?, ?, ?, ?, ?, ?, ?, ?);", pd.ID, pd.CreatedAt, pd.ExpiresAt, pd.Amount, pd.Status, pd.ReceiverId, "", pd.QRCodeData)
	if err != nil {
		return fmt.Errorf("falha ao criar novo pagamento")
	}

	return nil
}

func (p *paymentRepository) UpdatePaymentStatus(id string, status models.PaymentStatus) error {
	_, err := p.db.Exec("UPDATE payments SET status=? WHERE id=?", status, id)
	if err != nil {
		return fmt.Errorf("falha ao atualizar pagamento para status %s: %w", status, err)
	}

	return nil
}

func (p *paymentRepository) Delete(id string) error {
	_, err := p.db.Exec("DELETE FROM payments WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("falha ao deletar pagamento: %w", err)
	}

	return nil
}
