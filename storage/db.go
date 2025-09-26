package storage

import (
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"qr-payment/models"
	"qr-payment/utils"
)


var paymentsDB = make(map[string]*models.PaymentData)

var mu sync.Mutex

func GetPayment(id string) (*models.PaymentData, error) {
	mu.Lock()
	defer mu.Unlock()
	payment, exists := paymentsDB[id]
	if !exists {
		return nil, fmt.Errorf("payment not found")
	}
	return payment, nil
}

func GetPayments() (map[string]*models.PaymentData, error) {
	mu.Lock()
	defer mu.Unlock()
	return paymentsDB, nil
}

func CreatePayment(p *models.PaymentData) error {
	mu.Lock()
	defer mu.Unlock()
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
	paymentsDB[id] = p
	return nil
}

func MakePayment(id string) error {
	mu.Lock()
	defer mu.Unlock()
	payment, exists := paymentsDB[id]
	if !exists {
		return fmt.Errorf("payment not found")
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
		payment.Status = models.StatusExpired
		return fmt.Errorf("payment has expired")
	}
	
	payment.Status = models.StatusPaid
	return nil
}

func RemovePayment(id string) (error) {
	mu.Lock()
	defer mu.Unlock()
	_, exists := paymentsDB[id]
	if !exists {
		return fmt.Errorf("payment not found")
	}
	delete(paymentsDB, id)
	return nil
}