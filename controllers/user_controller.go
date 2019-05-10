package controllers

import (
		"encoding/json"
		"fmt"
		"github.com/aplJake/reals-course/server/models"
		"github.com/aplJake/reals-course/server/utils"
		"net/http"
)

// Authentication control
var GetUser = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))

}

var UserSignUp = func(w http.ResponseWriter, r *http.Request) {
		user := &models.User{}

		// Decode the request to server
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
				utils.Respond(w, utils.Message(false, "Invalid request"))
				return
		}
		// CreateSeller new User and UserProfile
		resp := user.Create()
		fmt.Println("Response", resp)
		utils.Respond(w, resp)
}

var UserSignIn = func(w http.ResponseWriter, r *http.Request) {
		user := &models.User{}

		// Decode the request body from client to struct
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
				utils.Respond(w, utils.Message(false, "Invalid request"))
				return
		}

		// Respond to the client and ...
		resp := models.LogIn(user.Email, user.Password)
		utils.Respond(w, resp)
}
