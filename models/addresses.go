package models

type Country struct {
		CountryId   uint   `json:"country_id"`
		CountryName string `json:"country_name"`
		ZipCode     string `json:"zip_code"`
}

type City struct {
		CityId    uint   `json:"city_id"`
		CityName  string `json:"city_name"`
		CountryId uint   `json:"country_id"`
}

type Addresses struct {
		AddressesId  uint   `json:"addresses_id"`
		StreetName   string `json:"street_name"`
		StreetNumber string `json:"street_number"`
		CityId       uint   `json:"city_id"`
}

// Router JSON REQUEST MODEL
type AddressesRequest struct {
		StreetName   string `json:"street_name"`
		StreetNumber string `json:"street_number"`
		CityName     string `json:"city_name"`
		CountryName  string `json:"country_name"`
		ZipCode      string `json:"zip_code"`
}
