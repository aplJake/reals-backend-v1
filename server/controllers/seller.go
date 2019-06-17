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
func NewPropertyListing(w http.ResponseWriter, r *http.Request) {
		//var listing models.PropertyListingRequest
		listing := &models.PropertyListingRequest{}

		fmt.Println("Request 23")
		fmt.Println(r.Body)


		err := json.NewDecoder(r.Body).Decode(&listing)
		fmt.Println(listing.UserId)
		if err != nil {
				utils.Respond(w, utils.Message(false, "Failed to decode the listing data"))
				panic(err.Error())
				return
		}

		fmt.Println("Lisiting", listing)


		respond := models.CreateListing(listing)
		utils.Respond(w, respond)
}

func GetSeller(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome23"))
}

// GET request handler to get the initial info
// about the property on update page
func GetPropertyListingUpdate(w http.ResponseWriter, r *http.Request)  {
	var (
		listing *models.PropertyListingRequest
		err error
	)
	propertyID := r.Context().Value("propertyID").(string)

	fmt.Println("Property id 1", propertyID)

	// Request all the data from the database
	listing, err = models.GetPropertyListing(propertyID)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Listing", listing)

	resp := utils.Message(true, "Listings are sended")
	resp["listing"] = listing
	// Respond to the client and ...
	utils.Respond(w, resp)
}

// PUT request to submit the updating info
// about the property
func PropertyListingUpdate(w http.ResponseWriter, r *http.Request) {
	listing := &models.PropertyListingRequest{}

	fmt.Print("Update property listings ")

	err := json.NewDecoder(r.Body).Decode(&listing)
	fmt.Println(listing.UserId)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Failed to decode the listing update data"))
		panic(err.Error())
		return
	}

	err = listing.UpdateListing()
	if err != nil {
		panic(err.Error())
	}
	utils.Respond(w, utils.Message(true, "Property listing was sucessfully updated"))
}

// DELETE request handler to remove
// the property from user profile
func PropertiesListingDelete(w http.ResponseWriter, r *http.Request) {

}
