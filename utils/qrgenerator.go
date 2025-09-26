package utils

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(id string) ([]byte, error) {
	var png []byte
	url := "http://localhost:8080/payments/" + id + "/pay"
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate qrcode")
	}
	return png, nil
}