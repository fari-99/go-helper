package models

type TransactionAddress struct {
	TransactionUuid   string `json:"transaction_uuid"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Address           string `json:"address"`
	AdditionalAddress string `json:"additional_address"`
	Postcode          string `json:"postcode"`
	EmailAddress      string `json:"email_address"`
	Phone             string `json:"phone"`

	CountryName  string `json:"country_name"`
	CountryCode  string `json:"country_code"`
	ProvinceName string `json:"province_name"`
	ProvinceCode string `json:"province_code"`
	CityName     string `json:"city_name"`
	CityCode     string `json:"city_code"`
}
