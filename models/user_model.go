package models

import (
		"database/sql"
		"fmt"
		"github.com/aplJake/reals-course/server/utils"
		"github.com/dgrijalva/jwt-go"
		"golang.org/x/crypto/bcrypt"
		"strings"
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

// Creation of User and UserProfile
func (user User) Create() map[string]interface{} {
		if resp, ok := user.Validate(); !ok {
				return resp
		}

		// Create hashed password with bcrypt
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)

		// Query to add new user to the database
		res, err := GetDb().Exec("INSERT INTO users(user_name, user_email, user_password) VALUE(?,?,?)",
				user.UserName, user.Email, user.Password)

		temp, err := res.LastInsertId()
		if err != nil {
				panic(err.Error())
				return utils.Message(false, "Failed to convert user id, convertion error.")
		}
		user.ID = uint(temp)

		if err != nil {
				panic(err.Error())
		}

		if user.ID < 0 {
				return utils.Message(false, "Failed to create account, connection error.")
		}

		// Create JWT Token
		fmt.Println("Create JWT Token User id", user.ID)
		tk := &Token{
				UserId:  user.ID,
				IsAdmin: false,
		}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
		tokenString, _ := token.SignedString([]byte(signedString))
		user.Token = tokenString

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

func LogIn(email, password string) map[string]interface{} {
		// create the ref to User and search in db the user with some email
		user := &User{}
		row := GetDb().QueryRow("SELECT * FROM users WHERE user_email=?", email)
		err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Token)

		if err != nil {
				if err == sql.ErrNoRows {
						return utils.Message(false, "Email address not found")
				}
				return utils.Message(false, "Connection error. Please retry")
		}

		// Decrypt the User password and login password and compare em
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
				return utils.Message(false, "invalid login creadentials. Please try again.")
		}

		// Delete password for safe client response
		user.Password = ""

		// Create new bcrypt token and JWT token
		tk := &Token{UserId: user.ID}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
		tokenString, _ := token.SignedString([]byte(signedString))
		user.Token = tokenString

		response := utils.Message(true, "Logged In")
		response["account"] = user

		fmt.Println("*********************************")
		fmt.Printf("SIGN IN is ID=%s \n name=%s \n email=%s \n password=%s \n hashed=%s", user.ID, user.UserName, user.Email, user.Password, user.Token)
		fmt.Println("User ", user)
		fmt.Println(response)
		fmt.Println("*********************************\n")

		return response
}

/*
	Email has to contain @ sign
	Password has to be not less than 6
	Email address must be unique
*/
func (user *User) Validate() (map[string]interface{}, bool) {
		if !strings.Contains(user.Email, "@") {
				return utils.Message(false, "Email address is required"), false
		}
		if len(user.Password) < 6 {
				return utils.Message(false, "Password is required"), false
		}

		temp := &User{}
		row := GetDb().QueryRow("SELECT user_email FROM users WHERE user_email=?", user.Email)
		err := row.Scan(&temp.Email)

		// Check if database has some errors
		// HAS ERROR 		(no such row) that the data is unique and evth is ok
		// HAS NOT ERROR	the data is not unique and it is selected from the db
		if err != sql.ErrNoRows {
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
		row := GetDb().QueryRow("SELECT * FROM users WHERE user_id=?", u)
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