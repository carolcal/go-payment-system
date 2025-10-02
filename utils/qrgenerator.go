package utils

import (
	"fmt"
	"time"
)

func GenerateQRCode(pixKey string, amount int, name string, city string) string {
	result := ""
	result += tlv("00", "01") //payload
	result += getMerchantAccountInfo(pixKey)
	result += tlv("52", "0000") //merchant category code
	result += tlv("53", "986") //currency

	if amount > 0 {
		result += getTransationAmount(amount)
	}
	
	result += tlv("58", "BR") //country code
	result += getMerchantName(name)
	result += getMerchantCity(city)
	result += getAdditionalDataField()
	result += getCRC(result)

	return result
}

// tlv = ID(2) + LEN(2 decimal) + VALUE
func tlv(id string, value string) string {
	l := len([]byte(value)) // bytes length (UTF-8)
	return fmt.Sprintf("%s%02d%s", id, l, value)
}

func getMerchantAccountInfo(pixKey string) string {
	gui := tlv("00", "br.gov.bcb.pix")
	pix := tlv("01", pixKey)
	accountInfo := gui + pix
	return tlv("26", accountInfo)
}

func getTransationAmount(amount int) string {
	// Convert cents to decimal format (e.g., 1050 cents = "10.50")
	dollars := amount / 100
	cents := amount % 100
	amountStr := fmt.Sprintf("%d.%02d", dollars, cents)
	return tlv("54", amountStr)
}

func getMerchantName(name string) string {
	maxLen := 25
	if len(name) > maxLen {
		name = name[:maxLen]
	}
	return tlv("59", name)
}

func getMerchantCity(city string) string {
	maxLen := 15
	if len(city) > maxLen {
		city = city[:maxLen]
	}
	return tlv("60", city)
}

func getAdditionalDataField() string {
	txid := tlv("05", generateTXID())
	return tlv("62", txid)
}

func generateTXID() string {
	now := time.Now()
	formattedTime := now.Format("20060102150405")
	return fmt.Sprintf("42%s", formattedTime)
}

func getCRC(data string) string {
	// Calculate CRC on data + "6304" (CRC field indicator and length)
	dataForCRC := data + "6304"
	crcValue := crc16([]byte(dataForCRC))
	// Ensure CRC is exactly 4 hexadecimal digits
	return fmt.Sprintf("6304%04X", crcValue&0xFFFF)
}

func crc16(data []byte) uint16 {
    var crc uint16 = 0xFFFF
    for _, b := range data {
        crc ^= uint16(b) << 8
        for i := 0; i < 8; i++ {
            if (crc & 0x8000) != 0 {
                crc = (crc << 1) ^ 0x1021
            } else {
                crc <<= 1
            }
        }
    }
    return crc & 0xFFFF
}
