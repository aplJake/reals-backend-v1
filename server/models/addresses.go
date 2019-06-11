package models

import (
	"database/sql"
	"errors"
	"log"
)

type Country struct {
		CountryId   uint   `json:"country_id"`
		CountryName string `json:"country_name"`
		CountryCode string `json:"country_code"`
}

var getAllCountriesQ = `
	SELECT * FROM country;
`

var getCountryQ = ` SELECT * FROM country WHERE country_id=?`

func DBGetCountry(countryID string) Country  {
		var country Country

		db = InitDB()
		row := db.QueryRow(getCountryQ, countryID)
		err := row.Scan(&country.CountryId, &country.CountryName, &country.CountryCode)
		handleError(err)

		defer db.Close()
		return country
}

func GetAllCountries() ([]Country, error) {
		db := InitDB()

		res, err := db.Query(getAllCountriesQ)
		handleError(err)

		country := Country{}
		countryArr := []Country{}
		for res.Next() {
				err = res.Scan(&country.CountryId, &country.CountryName, &country.CountryCode)
				handleError(err)
				countryArr = append(countryArr, country)
		}

		defer db.Close()
		return countryArr, nil
}

func (country *Country) Create() error  {
	// Validate Country obj (it must be unique in the database
	if ok := CountryExists(country.CountryName); ok {
		log.Println("Country name is exists")
		return errors.New("County with such name is already exists")
	}

	db := InitDB()

	_, err := db.Exec("INSERT into country(country_name, zip_code) VALUES (?,?)",
		country.CountryName, country.CountryCode)
	defer db.Close()

	if err != nil {
		log.Fatal(err.Error())
		return errors.New("Cannot insert a new country to the db")
	}
	return nil
}

func (country *Country) Update() error  {
	db := InitDB()

	_, err := db.Exec("UPDATE country SET country_name=?, zip_code=? WHERE country_id=?",
		country.CountryName, country.CountryCode, country.CountryId)
	defer db.Close()

	if err != nil {
		log.Fatal(err.Error())
		return errors.New("Cannot update a country in db")
	}
	return nil
}

func CountryExists(countryName string) bool {
	db := InitDB()
	sqlStmt := "SELECT country_name FROM country WHERE country_name=?"
	err := db.QueryRow(sqlStmt, countryName).Scan(&countryName)

	defer db.Close()

	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			log.Print(err)
		}

		return false
	}

	return true
}

func hadleError(err error) {
		if err != nil {
				panic(err.Error())
		}
}

type City struct {
		CityId    uint   `json:"city_id"`
		CityName  string `json:"city_name"`
		CountryId uint   `json:"country_id"`
}
var getAllCitiesQ = `
	SELECT * FROM city;
`
func GetAllCities() ([]City, error) {
	db := InitDB()

	res, err := db.Query(getAllCitiesQ)
	handleError(err)

	city := City{}
	cityArr := []City{}
	for res.Next() {
		err = res.Scan(&city.CityId, &city.CountryId, &city.CityName)
		handleError(err)
		cityArr = append(cityArr, city)
	}

	defer db.Close()
	return cityArr, nil
}
var findCitiesByCountryQ = `
	SELECT * FROM city WHERE country_id=?;
`
func (c *Country) FindCitiesByCountry() ([]City, error)  {
		var db *sql.DB

		db = InitDB()

		res, err := db.Query(findCitiesByCountryQ, c.CountryId)
		handleError(err)

		city := City{}
		citiesArr := []City{}
		for res.Next() {
				err = res.Scan(&city.CityId, &city.CountryId, &city.CityName)
				handleError(err)
				citiesArr = append(citiesArr, city)
		}

		defer db.Close()
		return citiesArr, nil
}

type Addresses struct {
		AddressesId  uint   `json:"addresses_id"`
		StreetName   string `json:"street_name"`
		StreetNumber string `json:"street_number"`
		CityId       uint   `json:"city_id"`
}

// Router JSON REQUEST MODEL
type AddressesRequest struct {
		AddressesID uint `json:"addresses_id,string"`
		CityID		uint	`db:"city_id" json:"city_id,string"`
		StreetName   string `db:"street_name" json:"street_name"`
		StreetNumber string `db:"street_number" json:"street_number"`
		CountryID	uint	`db:"country_id" json:"country_id,string"`
}
