package qrcode

import (
	"fmt"
)

func GetMerchantAccountInfo(pixKey string) string {
	gui := TLV("00", "br.gov.bcb.pix")
	pix := TLV("01", pixKey)
	accountInfo := gui + pix
	return TLV("26", accountInfo)
}

func GetTransationAmount(amount int) string {
	dollars := amount / 100
	cents := amount % 100
	amountStr := fmt.Sprintf("%d.%02d", dollars, cents)
	return TLV("54", amountStr)
}

func GetMerchantName(name string) string {
	maxLen := 25
	if len(name) > maxLen {
		name = name[:maxLen]
	}
	return TLV("59", name)
}

func GetMerchantCity(city string) string {
	maxLen := 15
	if len(city) > maxLen {
		city = city[:maxLen]
	}
	return TLV("60", city)
}

func GetAdditionalDataField() string {
	txid := TLV("05", generateTXID())
	return TLV("62", txid)
}

func GetCRC(data string) string {
	dataForCRC := data + "6304"
	crcValue := crc16([]byte(dataForCRC))
	return fmt.Sprintf("6304%04X", crcValue&0xFFFF)
}


