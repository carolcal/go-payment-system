package tests

import (
	"fmt"
	"testing"
	"qr-payment/qrcode"
)

func TestQR(t *testing.T) {
	// Test QR code generation
	pixKey := "92991514078"
	amount := 100 // 1.00 in cents
	name := "Arthur Dent"
	city := "Terra"
	
	qrCode := qrcode.GenerateQRCode(pixKey, amount, name, city)
	fmt.Printf("Generated QR Code: %s\n", qrCode)
	fmt.Printf("Length: %d\n", len(qrCode))
	
	// Let's also test with 0 amount
	qrCode2 := qrcode.GenerateQRCode(pixKey, 0, name, city)
	fmt.Printf("Generated QR Code (0 amount): %s\n", qrCode2)
	fmt.Printf("Length: %d\n", len(qrCode2))
}