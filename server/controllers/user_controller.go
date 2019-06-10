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
		resp := user.Create(w)
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


func GetUsers(w http.ResponseWriter, r *http.Request) {
		db := models.GetDb()

		queryRes, err := db.Query("SELECT * FROM listings ORDER BY user_id DESC")
		if err != nil {
				panic(err.Error())
		}

		var user models.User
		var users []models.User
		for queryRes.Next() {
				err := queryRes.Scan(&user.ID, &user.UserName, &user.Email, &user.Password)
				if err != nil {
						panic(err.Error())
				}

				users = append(users, user)
		}

		//defer db.Close()
		resp := utils.Message(true, "Users are recieved successfully")
		resp["listings"] = users
		utils.Respond(w, resp)
}

func CreateNewAdminUser(w http.ResponseWriter, r *http.Request) {
		user := &models.User{}

		// Decode the request to server
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
				utils.Respond(w, utils.Message(false, "Invalid request"))
				return
		}

		fmt.Println(user)
		// CreateSeller new Admin
		// Check if there exists such user with such id
		_, err = models.GetDb().Exec("SELECT * FROM listings WHERE user_id=?", user.ID)
		if err != nil {
				resp := utils.Message(false, "There doesn`t exist such user with this ID")
				utils.Respond(w, resp)
				return
		}

		// Insert to the Admin table a new User
		_, err = models.GetDb().Exec("INSERT INTO admins(user_id, admin_role) VALUES (?,?)",
				user.ID, "ADMIN_USER")


		resp := utils.Message(true, "A new admin user was successfully created")
		utils.Respond(w, resp)
}