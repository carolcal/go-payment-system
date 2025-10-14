package services

import (
	"context"
	"time"

	"qr-payment/internal/core/models"
	"qr-payment/internal/core/validators"
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
    ProcessPayment(ctx context.Context, user_id string, ppd *models.ProcessPaymentData) error
	RemovePayment(ctx context.Context, id string) error
}

type paymentService struct {
    repo repository.PaymentRepository
	val validators.PaymentValidator
    userService UserService
}

func NewPaymentService(repo repository.PaymentRepository, val validators.PaymentValidator, userService UserService) PaymentService {
    return &paymentService{
        repo: repo,
		val: val,
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
		return nil, &models.Err{Op: "PaymentService.GetAllPaymentsByUserId", Status: models.Invalid, Msg: "'user_type' must be 'receiver_id' or 'payer_id'."}
	}
    return s.repo.FindAllByUserId(user_type, user_id)
}

func (s *paymentService) CreatePayment(ctx context.Context, cpd *models.CreatePaymentData) (*models.PaymentData, error) {
    user, err := s.userService.GetUserById(ctx, cpd.ReceiverId)
	if err != nil {
		return nil, err
	}

    id := utils.GenerateID("pay")
	amount := 0
	if cpd.Amount != nil {
		if *cpd.Amount <= 0 {
			return nil, &models.Err{Op: "PaymentService.CreatePayment", Status: models.Invalid, Msg: "Amount must be a positive number."}
		}
		amount = int(*cpd.Amount * 100)
	}
    
	pd := &models.PaymentData{
		ID:         id,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(15 * time.Minute),
		Amount:     amount,
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

func (s *paymentService) ProcessPayment(ctx context.Context, payer_id string, ppd *models.ProcessPaymentData) error {
    qrdata, err := qrcode.ParseQrCodeData(ppd.QRCodeData)
	if err != nil {
		return err
	}

	amount, err := s.val.ValidatePaymentAmount(qrdata, ppd)
	if err != nil {
		return err
	}

    receiver, err := s.userService.GetUserByNameAndCPF(ctx, qrdata.MerchantName, qrdata.PixKey)
	if err != nil {
		return err
	}

	payment, err := s.GetPaymentByQRCodeData(ctx, ppd.QRCodeData)
    if err != nil {
        return err
    }
	err = s.val.ValidatePaymentStatus(payment)
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

	err = s.repo.UpdatePaymentPayerId(payment.ID, payer_id)
	if err != nil {
		return err
	}

    return s.repo.UpdatePaymentStatus(payment.ID, models.StatusPaid)
}

func (s *paymentService) RemovePayment(ctx context.Context, id string) error {
    _, err := s.repo.FindById(id)
	if err != nil {
		return err
	}
    return s.repo.Delete(id)
}
