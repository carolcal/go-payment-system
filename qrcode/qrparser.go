package qrcode

import (
	"fmt"
	"qr-payment/models"
	"strconv"
)

func ParseQrCodeData(qr_code_data string) (*models.QRCodeData, error) {
	data := models.QRCodeData{}
	i := 0

	for i < (len(qr_code_data) - 4) {
		tag := qr_code_data[i:i + 2]
		i += 2

		lenStr := qr_code_data[i:i + 2]
		i += 2

		length, err := strconv.Atoi(lenStr)
		if err != nil || i + length > len(qr_code_data) {
			return nil, fmt.Errorf("qrcode inválido")
		}

		value := qr_code_data[i:i + length]
		i += length

		err = assignTagAndValue(&data, tag, value)
		if err != nil {
			return nil, fmt.Errorf("qrcode inválido")
		}
	}
	PrintQRCodeData(&data)
	return &data, nil
}

func assignTagAndValue(data *models.QRCodeData, tag string, value string) error {
	switch tag {
		case "00":
			data.PayloadFormatIndicator = value
		case "01":
		case "26":
			data.MerchantAccountInfo = value
			pixkey, err := extractPixKey(value)
			if err != nil {
				return err
			}
			data.PixKey = pixkey
		case "52":
			data.MerchantCategoryCode = value
		case "53":
			data.TransactionCurrency = value
		case "54":
			data.TransactionAmount = value
		case "58":
			data.CountryCode = value
		case "59":
			data.MerchantName = value
		case "60":
			data.MerchantCity = value
		case "62":
			data.AdditionalDataField = value
		case "63":
			data.CRC = value
		default:
			return fmt.Errorf("unknown tag")
	}
	return nil
}

func extractPixKey(merchantAccountInfo  string) (string, error) {
	i := 0
	for i < len(merchantAccountInfo) {
		tag := merchantAccountInfo[i:i + 2]
		i += 2

		lenStr := merchantAccountInfo[i:i + 2]
		i += 2

		length, err := strconv.Atoi(lenStr)
		if err != nil || i + length > len(merchantAccountInfo) {
			return "", err
		}

		if tag == "01" {
			return merchantAccountInfo[i:i + length], nil
		}
		i += length

	}
	return "", nil
}

func PrintQRCodeData(data *models.QRCodeData) {
    if data == nil {
        fmt.Println("nil QRCodeData")
        return
    }
    fmt.Println("QRCodeData:")
    fmt.Printf("  PayloadFormatIndicator: %s\n", data.PayloadFormatIndicator)
    fmt.Printf("  MerchantAccountInfo:    %s\n", data.MerchantAccountInfo)
    fmt.Printf("  PixKey:                 %s\n", data.PixKey)
    fmt.Printf("  MerchantCategoryCode:   %s\n", data.MerchantCategoryCode)
    fmt.Printf("  TransactionCurrency:    %s\n", data.TransactionCurrency)
    fmt.Printf("  TransactionAmount:      %s\n", data.TransactionAmount)
    fmt.Printf("  CountryCode:            %s\n", data.CountryCode)
    fmt.Printf("  MerchantName:           %s\n", data.MerchantName)
    fmt.Printf("  MerchantCity:           %s\n", data.MerchantCity)
    fmt.Printf("  AdditionalDataField:    %s\n", data.AdditionalDataField)
    fmt.Printf("  CRC:                    %s\n", data.CRC)
}