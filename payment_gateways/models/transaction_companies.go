package models

type TransactionCompanies struct {
	TransactionUuid string `json:"transaction_uuid"`
	CompanyID       uint64 `json:"company_id"`
	CompanyCode     string `json:"company_code"`
	Name            string `json:"name"`
	Address         string `json:"address"`
	Postcode        string `json:"postcode"`
	Email           string `json:"email"`
	MobilePhone     string `json:"mobile_phone"`

	CountryName  string `json:"country_name"`
	CountryCode  string `json:"country_code"`
	ProvinceName string `json:"province_name"`
	ProvinceCode string `json:"province_code"`
	CityName     string `json:"city_name"`
	CityCode     string `json:"city_code"`
}
