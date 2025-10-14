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
	UpdatePaymentPayerId(id string, payer_id string) error
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
			return nil, &models.Err{Op: "PaymentRepository.FindById", Status: models.NotFound, Msg: "Payment Not Found."}
		}
		return nil, &models.Err{Op: "PaymentRepository.FindById", Status: models.Dependency, Msg: "Error when scanning payment.", Err: err}
	}

	return &payment, nil
}

func (p *paymentRepository) FindByQRCodeData(qr_code_data string) (*models.PaymentData, error) {
	var payment models.PaymentData

	row := p.db.QueryRow(`SELECT * FROM payments WHERE qr_code_data=?`, qr_code_data)
	err := utils.ScanPaymentRow(row, &payment)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.Err{Op: "PaymentRepository.FindByQRCodeData", Status: models.NotFound, Msg: "Payment Not Found."}
		}
		return nil, &models.Err{Op: "PaymentRepository.FindByQRCodeData", Status: models.Dependency, Msg: "Error when scanning payment.", Err: err}
	}

	return &payment, nil
}


func (p *paymentRepository) FindAll() (map[string]*models.PaymentData, error) {
	allPayments := make(map[string]*models.PaymentData)

	rows, err := p.db.Query(`SELECT * FROM payments`)
	if err != nil {
		return nil, &models.Err{Op: "PaymentRepository.FindAll", Status: models.Dependency, Msg: "Error executing payment query.", Err: err}
	}
	defer rows.Close()

	for rows.Next() {
		var payment models.PaymentData
		err := utils.ScanPaymentRows(rows, &payment)
		if err != nil {
			return nil, &models.Err{Op: "PaymentRepository.FindAll", Status: models.Dependency, Msg: "Error when scanning payment.", Err: err}
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
		return nil, &models.Err{Op: "PaymentRepository.FindAllByUserId", Status: models.Dependency, Msg: "Error executing payment query.", Err: err}
	}
	defer rows.Close()

	for rows.Next() {
		var payment models.PaymentData
		err := utils.ScanPaymentRows(rows, &payment)
		if err != nil {
			return nil, &models.Err{Op: "PaymentRepository.FindAllByUserId", Status: models.Dependency, Msg: "Error when scanning payment.", Err: err}
		}
		allPayments[payment.ID] = &payment
	}

	return allPayments, nil
}

func (p *paymentRepository) Create(u *models.UserData, pd *models.PaymentData) error {
	_, err := p.db.Exec("INSERT INTO payments VALUES(?, ?, ?, ?, ?, ?, ?, ?);", pd.ID, pd.CreatedAt, pd.ExpiresAt, pd.Amount, pd.Status, pd.ReceiverId, "", pd.QRCodeData)
	if err != nil {
		return &models.Err{Op: "PaymentRepository.Create", Status: models.Dependency, Msg: "Failed to create new payment.", Err: err}
	}

	return nil
}

func (p *paymentRepository) UpdatePaymentStatus(id string, status models.PaymentStatus) error {
	_, err := p.db.Exec("UPDATE payments SET status=? WHERE id=?", status, id)
	if err != nil {
		return &models.Err{Op: "PaymentRepository.UpdatePaymentStatus", Status: models.Dependency, Msg: "Fail to update payment.", Err: err}
	}

	return nil
}

func (p *paymentRepository) UpdatePaymentPayerId(id string, payer_id string) error {
	_, err := p.db.Exec("UPDATE payments SET payer_id=? WHERE id=?", payer_id, id)
	if err != nil {
		return &models.Err{Op: "PaymentRepository.UpdatePaymentPayerId", Status: models.Dependency, Msg: "Fail to update payment.", Err: err}
	}

	return nil
}

func (p *paymentRepository) Delete(id string) error {
	_, err := p.db.Exec("DELETE FROM payments WHERE id=?", id)
	if err != nil {
		return &models.Err{Op: "PaymentRepository.Delete", Status: models.Dependency, Msg: "Fail to delete payment.", Err: err}
	}

	return nil
}
