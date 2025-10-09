package services

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"qr-payment/internal/core/models"
	"qr-payment/internal/infrastructure/repository"
	"qr-payment/internal/qrcode"
	"qr-payment/internal/utils"
)

type PaymentService interface {
    GetPaymentById(ctx context.Context, id string) (*models.PaymentData, error)
    GetPaymentByQRCodeData(ctx context.Context, qr_code_data string) (*models.PaymentData, error)
	GetAllPayments(ctx context.Context) (map[string]*models.PaymentData, error)
	GetAllPaymentsByUserId(ctx context.Context, user_type_param string, user_id string) (map[string]*models.PaymentData, error)
    CreatePayment(ctx context.Context, req *models.CreatePaymentData) (*models.PaymentData, error)
    ValidatePaymentStatus(payment *models.PaymentData) error
    ProcessPayment(ctx context.Context, user_id string, qr_code_data string) error
	RemovePayment(ctx context.Context, id string) error
}

type paymentService struct {
    repo repository.PaymentRepository
    userService UserService
}

func NewPaymentService(repo repository.PaymentRepository, userService UserService) PaymentService {
    return &paymentService{
        repo: repo,
        userService: userService,
    }
}

func (s *paymentService) GetPaymentById(ctx context.Context, id string) (*models.PaymentData, error) {
    return s.repo.FindById(id)
}

func (s *paymentService) GetPaymentByQRCodeData(ctx context.Context, qr_code_data string) (*models.PaymentData, error) {
    return s.repo.FindByQRCodeData(qr_code_data)
}


func (s *paymentService) GetAllPayments(ctx context.Context) (map[string]*models.PaymentData, error) {
    return s.repo.FindAll()
}

func (s *paymentService) GetAllPaymentsByUserId(ctx context.Context, user_type_param string, user_id string) (map[string]*models.PaymentData, error) {
    user_type, isValid := models.IsValidTypeUser(user_type_param)
	if !isValid {
        return nil, fmt.Errorf("user_type deve ser 'receiver_id' ou 'payer_id'")
	}
    return s.repo.FindAllByUserId(user_type, user_id)
}

func (s *paymentService) CreatePayment(ctx context.Context, cpd *models.CreatePaymentData) (*models.PaymentData, error) {
    user, err := s.userService.GetUserById(ctx, cpd.ReceiverId)
	if err != nil {
		return nil, err
	}

    id := utils.GenerateID("pay")
    
	pd := &models.PaymentData{
		ID:         id,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(15 * time.Minute),
		Amount:     int(cpd.Amount * 100),
		ReceiverId: cpd.ReceiverId,
		Status:     models.StatusPending,
	}
    pd.QRCodeData = qrcode.GenerateQRCode(user.CPF, pd.Amount, user.Name, user.City)

    err = s.repo.Create(user, pd)
    if err != nil {
        return nil, err
    }
    return pd, nil
}

func (s *paymentService) ValidatePaymentStatus(payment *models.PaymentData) error {

	switch payment.Status {
	case models.StatusPaid:
		return fmt.Errorf("pagamento já realizado")
	case models.StatusFailed:
		return fmt.Errorf("pagamento falhou anteriormente")
	case models.StatusExpired:
		return fmt.Errorf("pagamento expirado")
	}

	if time.Now().After(payment.ExpiresAt) {
		err := s.repo.UpdatePaymentStatus(payment.ID, models.StatusExpired)
		if err != nil {
			return fmt.Errorf("pagamento expirou: falha ao atualizar pagamento para status expirado: %w", err)
		}
		return fmt.Errorf("pagamento expirou")
	}

	return nil

}

func (s *paymentService) ProcessPayment(ctx context.Context, payer_id string, qr_code_data string) error {
    qrdata, err := qrcode.ParseQrCodeData(qr_code_data)
	if err != nil {
		return err
	}

    amount, err := strconv.ParseFloat(qrdata.TransactionAmount, 64)
    if err != nil || amount <= 0 {
		return fmt.Errorf("invalid transaction amount")
	}

    receiver, err := s.userService.GetUserByNameAndCPF(ctx, qrdata.MerchantName, qrdata.PixKey)
	if err != nil {
		return err
	}
	
	payment, err := s.GetPaymentByQRCodeData(ctx, qr_code_data)
    if err != nil {
        return err
    }
	err = s.ValidatePaymentStatus(payment)
	if err != nil {
		return err
	}

    err = s.userService.UpdateBalance(ctx, payer_id, models.UpdateBalanceData{Diff: -amount})
    if err != nil {
        return err
    }

    err = s.userService.UpdateBalance(ctx, receiver.ID, models.UpdateBalanceData{Diff: amount})
    if err != nil {
        return err
    }

    return s.repo.UpdatePaymentStatus(payment.ID, models.StatusPaid)
}

func (s *paymentService) RemovePayment(ctx context.Context, id string) error {
    _, err := s.repo.FindById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("pagamento não encontrado")
		}
		return fmt.Errorf("erro ao escanear pagamento: %w", err)
	}
    return s.repo.Delete(id)
}
