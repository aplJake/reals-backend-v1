package controllers

import (
		"encoding/json"
		"github.com/aplJake/reals-course/server/models"
		"github.com/aplJake/reals-course/server/utils"
		"net/http"
)

// SellerHandler uses for creation of Seller model
// that uses when the initial PropertyListing is added to the database
var CreateSellerHandler = func(w http.ResponseWriter, r *http.Request) {
		seller := models.Seller{}

		err := json.NewDecoder(r.Body).Decode(seller)
		if err != nil {
				utils.Respond(w, utils.Message(false, "Failed to get the Seller object"))
				return
		}

		respond := seller.Create()
		utils.Respond(w, respond)
}
