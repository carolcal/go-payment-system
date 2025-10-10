package models

type QRCodeData struct {
	PayloadFormatIndicator	string
	MerchantAccountInfo		string
	PixKey					string
	MerchantCategoryCode	string
	TransactionCurrency		string
	TransactionAmount		string
	CountryCode				string
	MerchantName			string
	MerchantCity			string
	AdditionalDataField  	string
	CRC						string
}