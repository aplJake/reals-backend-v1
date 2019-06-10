package controllers

import (
		"github.com/aplJake/reals-course/server/models"
		"github.com/aplJake/reals-course/server/utils"
		"net/http"
)

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
