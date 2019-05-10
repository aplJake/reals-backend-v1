package controllers

import (
		"encoding/json"
		"fmt"
		"github.com/aplJake/reals-course/server/models"
		"github.com/aplJake/reals-course/server/utils"
		"net/http"
)

// SellerHandler uses for creation of Seller model
// that uses when the initial PropertyListing is added to the database
var NewPropertyListing = func(w http.ResponseWriter, r *http.Request) {
		listing := &models.PropertyListingRequest{}

		fmt.Println("Request 23")
		fmt.Println(r.Body)


		err := json.NewDecoder(r.Body).Decode(listing)
		fmt.Println(listing.UserId)
		if err != nil {
				utils.Respond(w, utils.Message(false, "Failed to get the Seller object"))
				panic(err.Error())
				return
		}

		fmt.Println("Lisiting", listing)


		respond := models.CreateListing(listing)
		utils.Respond(w, respond)
}

var GetSeller = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome23"))
}