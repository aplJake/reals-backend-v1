package models

import (
		"database/sql"
		"fmt"
		"github.com/aplJake/reals-course/server/utils"
		"github.com/dgrijalva/jwt-go"
		"golang.org/x/crypto/bcrypt"
		"log"
		"net/http"
		"strings"
		"time"
)

type User struct {
		ID       uint   `json:"user_id"`
		UserName string `json:"user_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Token    string `sql:"-" json:"token"`
}

// new MD5 string 7278AC6970BEC14904CFC8CB7D4F933C
var signedString = "some_string"

type Token struct {
		UserId  uint
		IsAdmin bool
		jwt.StandardClaims
}

var insertUserQ = `
		INSERT INTO users(
		user_name,
		user_email, 
		user_password)
		VALUES (?,?,?);
		`

// Creation of User and UserProfile
func (user User) Create(w http.ResponseWriter) map[string]interface{} {
		if resp, ok := user.Validate(); !ok {
				return resp
		}

		// CreateSeller hashed password with bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
				panic(err.Error())
		}
		user.Password = string(hashedPassword)

		// Query to add new user to the database
		res, err := GetDb().Exec(insertUserQ, user.UserName, user.Email, user.Password)
		if err != nil {
				panic(err.Error())
		}
		temp, err := res.LastInsertId()
		if err != nil {
				panic(err.Error())
				return utils.Message(false, "Failed to convert user id, convertion error.")
		}
		user.ID = uint(temp)

		// Create JWT Token
		experationTime := time.Now().Add(30 * time.Minute)
		token := CredentialsToken{
				UserId: user.ID,
				IsAdmin: false,
				StandardClaims: jwt.StandardClaims{
						ExpiresAt: experationTime.Unix(),
				},
		}
		jwtToken := token.CreateJWTToken(w, experationTime)

		//// CreateSeller JWT Token
		//fmt.Println("CreateSeller JWT Token User id", user.ID)
		//tk := &Token{
		//		UserId:  user.ID,
		//		IsAdmin: false,
		//}
		//token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
		//tokenString, _ := token.SignedString([]byte(signedString))
		user.Token = jwtToken

		// Link the Profile to the User
		err = user.InitProfile()
		if err != nil {
				panic(err.Error())
				return utils.Message(false, "Failed to create user profile.")

		}

		// Delete password for safe client response
		user.Password = ""

		// Respond to the client about success
		response := utils.Message(true, "Account has been created")
		response["account"] = user

		fmt.Println("*********************************")
		fmt.Printf("SIGN UP is ID=%s \n name=%s \n email=%s \n password=%s \n hashed=%s", user.ID, user.UserName, user.Email, user.Password, user.Token)
		fmt.Println("User ", user)
		fmt.Println(response)
		fmt.Println("*********************************\n")

		return response
}
/*
At login operation we check the login user data with USER TABLE
and then check it in ADMINS TABLE if such user exists in both tables
then lets check its payload to IS_ADMIN: TRUE
*/
func LogIn(email, password string) map[string]interface{} {
		// create the ref to User and search in db the user with some email
		user := &User{}
		row := GetDb().QueryRow("SELECT * FROM users WHERE user_email=?", email)
		err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Token)

		if err != nil {
				if err == sql.ErrNoRows {
						return utils.Message(false, "Email address not found")
				}
				log.Println(err.Error())
				return utils.Message(false, "Connection error. Please retry")
		}

		// Decrypt the User password and login password and compare em
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
				return utils.Message(false, "invalid login creadentials. Please try again.")
		}

		// Delete password for safe client response
		user.Password = ""

		_, status := GetAdmin(user.ID)

		// CreateSeller new bcrypt token and JWT token
		tk := &Token{UserId: user.ID, IsAdmin: status}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
		tokenString, _ := token.SignedString([]byte(signedString))
		user.Token = tokenString

		response := utils.Message(true, "Logged In")
		response["account"] = user
		return response
}

/*
	Email has to contain @ sign
	Password has to be not less than 6
	Email address must be unique
*/
func (user *User) Validate() (map[string]interface{}, bool) {
		var db *sql.DB
		if !strings.Contains(user.Email, "@") {
				return utils.Message(false, "Email address is required"), false
		}
		if len(user.Password) < 6 {
				return utils.Message(false, "Password is required"), false
		}

		temp := &User{}

		db = InitDB()
		row := db.QueryRow("SELECT user_email FROM users WHERE user_email=?", user.Email)
		err := row.Scan(&temp.Email)

		defer db.Close()

		// Check if database has some errors
		// HAS ERROR 		(no such row) that the data is unique and evth is ok
		// HAS NOT ERROR	the data is not unique and it is selected from the db
		if err != sql.ErrNoRows {
				panic(err.Error())
				return utils.Message(false, "Connection error. Please retry!"), false

				if temp.Email != "" {
						fmt.Println(temp.Email)
						return utils.Message(false, "Email address already in use by another user."), false
				}
		}

		return utils.Message(false, "Requirement passed"), true
}

func GetUser(u uint) *User {
		user := &User{}
		row := GetDb().QueryRow("SELECT * FROM listings WHERE user_id=?", u)
		err := row.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
				return nil
		}
		if user.Email == "" { //User not found!
				return nil
		}

		user.Password = ""
		return user
}

type Admin struct {
		UserId    uint   `json:"user_id"`
		AdminRole string `json:"admin_role"`
}

// Creates the initial super admin user in both tables
// Users and admin by transaction
func InitAdmin() map[string]interface{} {
		// Validate if there already exists SUPER_USER
		if resp, ok := ValidateSuperUser(); !ok {
				return resp
		}

		tx, err := GetDb().Begin()
		if err != nil {
				log.Fatal(err)
		}

		// Add admin to User Table
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("adminpassword"), bcrypt.DefaultCost)
		adminPassword := string(hashedPassword)
		res, err := tx.Exec("INSERT INTO listings(user_name, user_email, user_password) VALUES (?,?, ?)",
				"admin", "admin@gmail.com", adminPassword)

		if err != nil {
				err := tx.Rollback()
				log.Fatal(err)
				return utils.Message(false, "Admin superuser cannot be added to listings table")
		}

		// Add admin to Admins Table
		id, _ := res.LastInsertId()
		_, err = tx.Exec("INSERT INTO admins(user_id, admin_role) VALUES (?,?)",
				id, "SUPER_USER")

		if err != nil {
				err := tx.Rollback()
				log.Fatal(err)
				return utils.Message(false, "Admin superuser cannot be added to admins table")
		}

		if err := tx.Commit(); err != nil {
				return utils.Message(false, "Commit error for admin user")
		}

		return utils.Message(true, "Admin super user was successfully created")
}

func ValidateSuperUser() (map[string]interface{}, bool) {
		_, err := GetDb().Exec("SELECT user_id FROM admins WHERE admin_role=?", "SUPER_USER")
		if err != nil {
				resp := utils.Message(false, "Validation error while searching the superuser")
				return resp, false
		}
		resp := utils.Message(true, "Successful validation")
		return resp, true
}
func GetAdmin(u uint) (*Admin, bool) {
		admin := &Admin{}
		row := GetDb().QueryRow("SELECT * FROM admins WHERE user_id=?", u)
		err := row.Scan(&admin.UserId, &admin.AdminRole)
		if err != nil {
				return nil, false
		}
		return admin, true
}
