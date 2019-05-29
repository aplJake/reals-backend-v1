package controllers

import (
		"github.com/aplJake/reals-course/server/models"
		"github.com/aplJake/reals-course/server/utils"
		"net/http"
)

func GetAllListings(w http.ResponseWriter, r *http.Request) {
		var (
				listings []models.PropertyListing
			 	err error
		)
		// Request all the data from the database
		listings, err = models.GetAllListings()
		if err != nil {
				panic(err.Error())
		}

		resp := utils.Message(true, "Listings are sended")
		resp["listings"] = listings
		// Respond to the client and ...
		utils.Respond(w, resp)
}

func GetApartmentListings(w http.ResponseWriter, r *http.Request) {
		var (
				listings []models.PropertyListing
				err error
		)
		// Request all the data from the database
		listings, err = models.GetListingsByType("apartment")
		if err != nil {
				panic(err.Error())
		}

		resp := utils.Message(true, "Listings are sended")
		resp["listings"] = listings
		// Respond to the client and ...
		utils.Respond(w, resp)
}

func GetHomeListings(w http.ResponseWriter, r *http.Request) {
		var (
				listings []models.PropertyListing
				err error
		)
		// Request all the data from the database
		listings, err = models.GetListingsByType("home")
		if err != nil {
				panic(err.Error())
		}

		resp := utils.Message(true, "Listings are sended")
		resp["listings"] = listings
		// Respond to the client and ...
		utils.Respond(w, resp)
}