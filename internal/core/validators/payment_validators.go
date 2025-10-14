package validators

import (
	"strconv"
	"time"

	"qr-payment/internal/core/models"
	"qr-payment/internal/infrastructure/repository"
)

type PaymentValidator interface {
	ValidatePayment(qrdata *models.QRCodeData, ppd *models.ProcessPaymentData, payment *models.PaymentData) (float64, error)
	validatePaymentStatus(payment *models.PaymentData) error
	validatePaymentAmount(qrdata *models.QRCodeData, ppd *models.ProcessPaymentData, id string) (float64, error)
}

type paymentValidator struct {
	repo repository.PaymentRepository
}

func NewPaymentValidator(repo repository.PaymentRepository) PaymentValidator {
	return &paymentValidator{
		repo: repo,
	}
}

func (v *paymentValidator) ValidatePayment(qrdata *models.QRCodeData, ppd *models.ProcessPaymentData, payment *models.PaymentData) (float64, error) {
	err := v.validatePaymentStatus(payment)
	if err != nil {
		return 0, err
	}
	
	amount, err := v.validatePaymentAmount(qrdata, ppd, payment.ID)
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func (v *paymentValidator) validatePaymentStatus(payment *models.PaymentData) error {

	switch payment.Status {
	case models.StatusPaid:
		return &models.Err{Op: "ValidatePaymentStatus", Status: models.Conflict, Msg: "Payment already made."}
	case models.StatusFailed:
		return &models.Err{Op: "ValidatePaymentStatus", Status: models.Conflict, Msg: "Payment failed."}
	case models.StatusExpired:
		return &models.Err{Op: "ValidatePaymentStatus", Status: models.Precondition, Msg: "Payment expired."}
	}

	if time.Now().After(payment.ExpiresAt) {
		err := v.repo.UpdatePaymentStatus(payment.ID, models.StatusExpired)
		if err != nil {
			return &models.Err{Op: "ValidatePaymentStatus", Status: models.Dependency, Msg: "Failed to update payment status to expired.", Err: err}
		}
		return &models.Err{Op: "ValidatePaymentStatus", Status: models.Precondition, Msg: "Payment expired."}
	}

	return nil
}

func (v *paymentValidator) validatePaymentAmount(qrdata *models.QRCodeData, ppd *models.ProcessPaymentData, id string) (float64, error) {
	var amount float64

	if qrdata.TransactionAmount == "" {
		if ppd.Amount == nil {
			return 0, &models.Err{Op: "ValidatePaymentAmount", Status: models.Precondition, Msg: "Transaction amount is missing in QR code data, please provide payment value."}
		}
		if *ppd.Amount <= 0 {
			return 0, &models.Err{Op: "ValidatePaymentAmount", Status: models.Invalid, Msg: "Invalid provided amount."}
		}
		amount = *ppd.Amount
		v.repo.UpdatePaymentAmount(id, int(amount * 100))
	} else {
		parsedAmount, err := strconv.ParseFloat(qrdata.TransactionAmount, 64)
		if err != nil || parsedAmount <= 0 {
			v.repo.UpdatePaymentStatus(id, models.StatusFailed)
			return 0, &models.Err{Op: "ValidatePaymentAmount", Status: models.Invalid, Msg: "Invalid transaction amount.", Err: err}
		}
		amount = parsedAmount
	}
	return amount, nil
}
