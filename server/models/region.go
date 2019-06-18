package models

import "C"
import (
	"database/sql"
	"errors"
	"log"
)

type Country struct {
	CountryId   uint   `json:"country_id,string"`
	CountryName string `json:"country_name"`
	CountryCode string `json:"country_code"`
}

var getAllCountriesQ = `
	SELECT * FROM country ORDER BY country_id DESC;
`

var getCountryQ = ` SELECT * FROM country WHERE country_id=?`

func DBGetCountry(countryID string) Country {
	var country Country

	db = InitDB()
	row := db.QueryRow(getCountryQ, countryID)
	err := row.Scan(&country.CountryId, &country.CountryName, &country.CountryCode)
	handleError(err)

	defer db.Close()
	return country
}

func GetAllCountries(onlyWithCities bool) ([]Country, error) {
	db := InitDB()
	var query string

	if onlyWithCities {
		query = `SELECT * FROM country ORDER BY country_id DESC;`
	} else {
		query = `SELECT DISTINCT CON.country_id,
				        CON.country_name,
				        CON.zip_code
				FROM country CON
				    INNER JOIN city C on CON.country_id = C.country_id;
				`
	}
	res, err := db.Query(query)
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

func (country *Country) Create() error {
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

func (country *Country) Update() error {
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

func (city *City) Update() error {
	db := InitDB()

	_, err := db.Exec("UPDATE city SET country_id=?, city_name=? WHERE city_id=?",
		city.CountryId, city.CityName, city.CityId)
	defer db.Close()

	if err != nil {
		log.Fatal(err.Error())
		return errors.New("Cannot update a city in db")
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
	CityId    uint   `json:"city_id,string"`
	CityName  string `json:"city_name"`
	CountryId uint   `json:"country_id,string"`
}

var getAllCitiesQ = `
	SELECT * FROM city ORDER BY city_id DESC;
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

func (c *Country) FindCitiesByCountry() ([]City, error) {
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

func (city *City) Create() error {
	// Validate Country obj (it must be unique in the database
	if ok := CityExists(city.CityName); ok {
		log.Println("Country name is exists")
		return errors.New("County with such name is already exists")
	}

	db := InitDB()

	_, err := db.Exec("INSERT into city(country_id, city_name) VALUES (?, ?)",
		city.CountryId, city.CityName)
	defer db.Close()

	if err != nil {
		log.Fatal(err.Error())
		return errors.New("Cannot insert a new city to the db")
	}
	return nil
}

func CityExists(cityName string) bool {
	db := InitDB()
	sqlStmt := "SELECT city_name FROM city WHERE city_name=?"
	err := db.QueryRow(sqlStmt, cityName).Scan(&cityName)

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

type Regions struct {
	RegionID   uint   `json:"region_id,string"`
	RegionName string `json:"region_name"`
	CityId     uint   `json:"city_id,string"`
}

// Router JSON REQUEST MODEL
type AddressesRequest struct {
	RegionID   uint   `json:"region_id,string"`
	CityID     uint   `json:"city_id,string"`
	RegionName string `json:"region_name"`
	CountryID  uint   `json:"country_id,string"`
}

func (region *Regions) Create() error {
	// Validate Country obj (it must be unique in the database
	if ok := RegionExists(region.RegionName); ok {
		log.Println("Region name already exists")
		return errors.New("Region with such name already exists")
	}

	db := InitDB()

	_, err := db.Exec("INSERT into regions(city_id, region_name) VALUES (?, ?)",
		region.CityId, region.RegionName)
	defer db.Close()

	if err != nil {
		log.Fatal(err.Error())
		return errors.New("Cannot insert a new region to the db")
	}
	return nil
}

func (region *Regions) Update() error {
	db := InitDB()

	_, err := db.Exec("UPDATE regions SET city_id=?, region_name=? WHERE region_id=?",
		region.CityId, region.RegionName, region.RegionID)
	defer db.Close()

	if err != nil {
		log.Fatal(err.Error())
		return errors.New("Cannot update a region in db")
	}
	return nil
}

func DeleteRegion(id string) error {
	db := InitDB()
	_, err := db.Exec("DELETE FROM regions where region_id=?;", id)
	if err != nil {
		return errors.New("Delete region in the parent table")
	}

	defer db.Close()
	return nil
}

func RegionExists(cityName string) bool {
	db := InitDB()
	sqlStmt := "SELECT city_name FROM city WHERE city_name=?"
	err := db.QueryRow(sqlStmt, cityName).Scan(&cityName)

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

func GetAllRegions() ([]Regions, error) {
	db := InitDB()

	res, err := db.Query("SELECT * FROM regions")
	handleError(err)
	defer db.Close()

	region := Regions{}
	regionsArr := []Regions{}
	for res.Next() {
		err = res.Scan(&region.RegionID, &region.CityId, &region.RegionName)
		handleError(err)
		regionsArr = append(regionsArr, region)
	}

	return regionsArr, nil
}
