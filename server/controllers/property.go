package controllers

import (
		"fmt"
		"github.com/aplJake/reals-course/server/models"
		"github.com/aplJake/reals-course/server/utils"
		"net/http"
)

func GetAllListings(w http.ResponseWriter, r *http.Request) {
		var (
				listings []models.PropertyListing
				err      error
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
				err      error
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
				err      error
		)
		// Request all the data from the database
		listings, err = models.GetListingsByType("house")
		if err != nil {
				panic(err.Error())
		}

		resp := utils.Message(true, "Listings are sended")
		resp["listings"] = listings
		// Respond to the client and ...
		utils.Respond(w, resp)
}

func GetPropertyPageData(w http.ResponseWriter, r *http.Request) {
		var (
				listingPageData models.PropertyPageData
				//err error
		)
		listingPageData = r.Context().Value("propertyID").(models.PropertyPageData)
		fmt.Println(listingPageData)
		resp := utils.Message(true, "Listing page data was successfully sended")
		resp["listing_data"] = listingPageData
		// Respond to the client and ...
		utils.Respond(w, resp)
}

func GetPropertyQueue(w http.ResponseWriter, r *http.Request) {
		var propertyPageData models.PropertyCtxData

		propertyPageData = r.Context().Value("propertyData").(models.PropertyCtxData)
		resp := utils.Message(true, "Queues are sended")
		resp["property_queue_data"] = propertyPageData
		utils.Respond(w, resp)
}