package utils

import (
	"fmt"
	"log"
	"net/http"
	"io"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(id string) ([]byte, error) {
	var png []byte
	
	publicIP := getMyPublicIP()

	url := fmt.Sprintf("http://%s:8080/payment/%s/pay", publicIP, id)
	fmt.Printf("Generated URL: %s\n", url)
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate qrcode")
	}
	return png, nil
} 

func getMyPublicIP() string {
	ip, err := http.Get("https://api.ipify.org")
	if err != nil {
		log.Fatalf("Error making HTTP request: %v", err)
	}
	defer ip.Body.Close()

	ipBytes, err := io.ReadAll(ip.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	return string(ipBytes)
}