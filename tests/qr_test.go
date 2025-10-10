package tests

import (
	"fmt"
	"qr-payment/internal/qrcode"
	"testing"
)

func TestQRWithValidDataAndValue(t *testing.T) {
	pixKey := "92991514078"
	amount := 100
	name := "Arthur Dent"
	city := "Terra"
	
	qrCode := qrcode.GenerateQRCode(pixKey, amount, name, city)

	result, err := qrcode.ParseQrCodeData(qrCode)
	if err != nil {
		t.Errorf("Expected valid QRCodeData, got error: %v", err)
		return
	}
	if (*result).PixKey != pixKey {
		t.Errorf("Expected PixKey: %s, got: %s", pixKey, (*result).PixKey)
	}
	if (*result).TransactionAmount != fmt.Sprintf("%d.00", amount/100) {
		t.Errorf("Expected TransactionAmount: %d.00, got: %s", amount, (*result).TransactionAmount)
	}
	if (*result).MerchantName != name {
		t.Errorf("Expected MerchantName: %s, got: %s", name, (*result).MerchantName)
	}
	if (*result).MerchantCity != city {
		t.Errorf("Expected MerchantCity: %s, got: %s", city, (*result).MerchantCity)
	}
	
}

func TestQRWithValidDataAndNoValue(t *testing.T) {
	pixKey := "92991514078"
	amount := 0
	name := "Arthur Dent"
	city := "Terra"
	
	qrCode := qrcode.GenerateQRCode(pixKey, amount, name, city)

	result, err := qrcode.ParseQrCodeData(qrCode)
	if err != nil {
		t.Errorf("Expected valid QRCodeData, got error: %v", err)
		return
	}
	if (*result).PixKey != pixKey {
		t.Errorf("Expected PixKey: %s, got: %s", pixKey, (*result).PixKey)
	}
	if (*result).TransactionAmount != "" {
		t.Errorf("Expected empty TransactionAmount, got: %s", (*result).TransactionAmount)
	}
	if (*result).MerchantName != name {
		t.Errorf("Expected MerchantName: %s, got: %s", name, (*result).MerchantName)
	}
	if (*result).MerchantCity != city {
		t.Errorf("Expected MerchantCity: %s, got: %s", city, (*result).MerchantCity)
	}
	
}

func TestQRWithInvalidData(t *testing.T) {
	invalidQrCode := "00020126330014BR.GOV.BCB.PIX0114+55219999999990212Pagamento de Servios5204000053039865406100.005802BR5913Arthur Dent6005Terra62070503***6304"
	
	_, err := qrcode.ParseQrCodeData(invalidQrCode)
	if err == nil {
		t.Errorf("Expected error for invalid QR code, got valid QRCodeData")
	}
}