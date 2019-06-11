package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/aplJake/reals-course/server/models"
	"github.com/aplJake/reals-course/server/utils"
	"log"
	"net/http"
)

// Get all Countries
func GetCountries(w http.ResponseWriter, r *http.Request) {
		var countries []models.Country
		var err error

		// Request all the data from the database
		countries, err = models.GetAllCountries()
		if err != nil {
				panic(err.Error())
		}

		resp := utils.Message(true, "Coutries are sended")
		resp["countries"] = countries
		// Respond to the client and ...
		utils.Respond(w, resp)
}

// Add Country
func AddNewCountry(w http.ResponseWriter, r *http.Request) {
	country := &models.Country{}

	// Decode the request to server
	err := json.NewDecoder(r.Body).Decode(&country)
	if err != nil {
		log.Println(err.Error())
		utils.Respond(w, utils.Message(false, "Cannot decode recieved json object"))
		return
	}

	fmt.Print(country)
	// CreateSeller new User and UserProfile
	err = country.Create()
	if err != nil {
		utils.Respond(w, utils.Message(false, "Cannot add new country object to the database"))
		log.Fatal(err.Error())
		return
	}

	resp := utils.Message(true, "New Country was successfully added")
	utils.Respond(w, resp)
}
// Edit Country
func UpdateCountry(w http.ResponseWriter, r *http.Request)  {
	country := &models.Country{}

	// Decode the request to server
	err := json.NewDecoder(r.Body).Decode(&country)
	if err != nil {
		log.Println(err.Error())
		utils.Respond(w, utils.Message(false, "Cannot decode recieved json object"))
		return
	}

	fmt.Print(country)
	// CreateSeller new User and UserProfile
	err = country.Update()
	if err != nil {
		utils.Respond(w, utils.Message(false, "Cannot update country object in the database"))
		log.Fatal(err.Error())
		return
	}

	resp := utils.Message(true, "Country was successfully updated")
	utils.Respond(w, resp)
}

// Remove Country
func DeleteCountry(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("Dele request")
	countryID := r.Context().Value("countryToDeleteID").(string)

	fmt.Println(countryID)

	db := models.InitDB()

	_, err := db.Exec("DELETE FROM country where country_id=?;", countryID)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	resp := utils.Message(true, "Country was successfully removed")
	utils.Respond(w, resp)
}


// CITIES

// Get all Cities
func GetCitiesList(w http.ResponseWriter, r *http.Request) {
	var cities []models.City
	var err error

	// Request all the data from the database
	cities, err = models.GetAllCities()
	if err != nil {
		panic(err.Error())
	}

	resp := utils.Message(true, "Cities are sended")
	resp["cities"] = cities
	// Respond to the client and ...
	utils.Respond(w, resp)
}
// Get City By Country
func GetCitiesByCountry(w http.ResponseWriter, r *http.Request) {
		var cities []models.City
		var err error

		country := r.Context().Value("coutryID").(models.Country)


		// Request all the data from the database
		cities, err = country.FindCitiesByCountry()
		if err != nil {
				panic(err.Error())
		}

		resp := utils.Message(true, "Cities are sended")
		resp["cities"] = cities
		// Respond to the client and ...
		utils.Respond(w, resp)
}



