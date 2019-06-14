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
		resp := models.LogIn(user.Email, user.Password, w)
		utils.Respond(w, resp)
}

var getUsersQ = `
	SELECT U.user_id,
		   U.user_name,
		   U.user_email
	FROM users U
	WHERE U.user_id not in (
		SELECT user_id FROM admins
		)
	ORDER BY user_id
		DESC;
`
func GetUsers(w http.ResponseWriter, r *http.Request) {
		db := models.InitDB()

		queryRes, err := db.Query(getUsersQ)
		if err != nil {
				panic(err.Error())
		}

		defer db.Close()

		var user models.User
		var users []models.User
		for queryRes.Next() {
				err := queryRes.Scan(&user.ID, &user.UserName, &user.Email)
				if err != nil {
						panic(err.Error())
				}

				users = append(users, user)
		}

		//defer db.Close()
		resp := utils.Message(true, "Users are recieved successfully")
		resp["users"] = users
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

		db := models.InitDB()

		_, err = db.Exec("SELECT * FROM users WHERE user_id=?", user.ID)
		if err != nil {
				resp := utils.Message(false, "There doesn`t exist such user with this ID")
				utils.Respond(w, resp)
				return
		}

		// Insert to the Admin table a new User
		_, err = db.Exec("INSERT INTO admins(user_id, admin_role) VALUES (?,?)",
				user.ID, "ADMIN_USER")

		defer db.Close()

		resp := utils.Message(true, "A new admin user was successfully created")
		utils.Respond(w, resp)
}

var getAdminsQ = `
	SELECT U.user_id,
		   U.user_name,
		   U.user_email
	FROM users U
		LEFT JOIN admins a on U.user_id = a.user_id
	WHERE a.admin_role != 'SUPER_USER'
	ORDER BY user_id
		DESC;
`

func GetAdmins(w http.ResponseWriter, r *http.Request) {
	db := models.InitDB()

	queryRes, err := db.Query(getAdminsQ)
	if err != nil {
		panic(err.Error())
	}

	var user models.User
	var users []models.User
	for queryRes.Next() {
		err := queryRes.Scan(&user.ID, &user.UserName, &user.Email)
		if err != nil {
			panic(err.Error())
		}

		users = append(users, user)
	}

	defer db.Close()
	resp := utils.Message(true, "Admins are recieved successfully")
	resp["admins"] = users
	utils.Respond(w, resp)
}

var deleteAdminsQ = `
	DELETE FROM admins where user_id=?;
`

func DeleteAdminUser(w http.ResponseWriter, r *http.Request) {
	adminID := r.Context().Value("adminToDeleteID").(string)

	fmt.Println(adminID)

	db := models.InitDB()

	_, err := db.Exec(deleteAdminsQ, adminID)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	resp := utils.Message(true, "Admin was successfully removed")
	utils.Respond(w, resp)
}