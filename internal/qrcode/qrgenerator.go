package qrcode

func GenerateQRCode(pixKey string, amount int, name string, city string) string {
	result := ""
	result += TLV("00", "01") //payload
	result += GetMerchantAccountInfo(pixKey)
	result += TLV("52", "0000") //merchant category code
	result += TLV("53", "986") //currency

	if amount > 0 {
		result += GetTransationAmount(amount)
	}
	
	result += TLV("58", "BR") //country code
	result += GetMerchantName(name)
	result += GetMerchantCity(city)
	result += GetAdditionalDataField()
	result += GetCRC(result)

	return result
}
