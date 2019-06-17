package controllers

import (
		"encoding/json"
		"fmt"
		"github.com/aplJake/reals-course/server/models"
		"github.com/aplJake/reals-course/server/utils"
		"net/http"
)

// GetProfile returns the profile data by request with context
// so the UserProfile must be in the context
func GetProfile(w http.ResponseWriter, r *http.Request) {
		// Context has the userId information that we use for
		// handle data from db by this id
		fmt.Println("1", r)
		profile := r.Context().Value("userId").(models.UserProfileRespond)
		fmt.Println("2", profile)

		listings, err := models.GetLisitingsByProfile(profile.UserID)

		respond, err := models.NewUserProfileResponse(profile, listings)
		if err != nil {
				utils.Respond(w, utils.Message(false, "Failed to get the UserProfile response"))
				return
		}

		utils.Respond(w, respond)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
		//profile := r.Context().Value("profile").(*models.UserProfile)
		profile := &models.UserProfile{}

		err := json.NewDecoder(r.Body).Decode(profile)
		if err != nil {
				utils.Respond(w, utils.Message(false, "Failed to get the UserProfile response"))
				return
		}
		fmt.Println("PUT Profile", profile.ProfileDescription)

		respond, err := models.UpdateUserProfileResponse(profile)
		if err != nil {
				utils.Respond(w, utils.Message(false, "Failed to get the UserProfile response"))
				return
		}

		utils.Respond(w, respond)
}

type Message struct {
		Message string `json:"message"`
		SecondMessage	string	`json:"second_message"`
}

func AddAds(w http.ResponseWriter, r *http.Request) {
		msg := &models.UserProfile{}

		err := json.NewDecoder(r.Body).Decode(msg)
		if err != nil {
				utils.Respond(w, utils.Message(false, "Failed to get the UserProfile response"))
				return
		}
		fmt.Println("Decoded 23", msg.ProfileDescription)
		//fmt.Println("Decoded ", msg.SecondMessage)

}