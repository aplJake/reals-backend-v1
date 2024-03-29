package models

import (
	"database/sql"
	"fmt"
	"github.com/aplJake/reals-course/server/utils"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"regexp"
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
	UserId   uint
	IsAdmin  bool
	UserType string
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
		UserId:   user.ID,
		IsAdmin:  false,
		UserType: "USER",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: experationTime.Unix(),
		},
	}
	jwtToken := token.CreateJWTToken(w, experationTime)

	user.Token = jwtToken

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
func LogIn(email, password string, w http.ResponseWriter) map[string]interface{} {
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

	admin, status := GetAdmin(user.ID)

	// CreateSeller new bcrypt token and JWT token
	experationTime := time.Now().Add(30 * time.Minute)

	token := CredentialsToken{
		UserId:   user.ID,
		IsAdmin:  status,
		UserType: admin.AdminRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: experationTime.Unix(),
		},
	}

	jwtToken := token.CreateJWTToken(w, experationTime)
	user.Token = jwtToken

	//token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	//tokenString, _ := token.SignedString([]byte(signedString))

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
	// Check username
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9 ]+$", user.UserName); !ok {
		return utils.Message(false, "Username must contain only letters and numbers"), false
	}

	// Check email
	pattern := "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	if ok, _ := regexp.MatchString(pattern, user.Email); !ok {
		return utils.Message(false, "Change your email address"), false
	}

	// Check password
	if len(user.Password) < 6 {
		return utils.Message(false, "Minimum password length must be 6 signs"), false
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
		//panic(err.Error())
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

type Admin struct {
	UserId    uint   `json:"user_id"`
	AdminRole string `json:"admin_role"`
}

// Creates the initial super admin user in both tables
// Users and admin by transaction
func InitAdmin() map[string]interface{} {
	// Validate if there already exists SUPER_USER
	if ok := AdminExists("SUPER_USER"); ok {
		fmt.Println(ok)
		resp := utils.Message(true, "Super user is already created")
		return resp
	}

	tx, err := GetDb().Begin()
	if err != nil {
		panic(err.Error())
		//log.Fatal(err)
	}

	// Add admin to User Table
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("adminpassword"), bcrypt.DefaultCost)
	adminPassword := string(hashedPassword)
	res, err := tx.Exec("INSERT INTO users(user_name, user_email, user_password) VALUES (?,?, ?)",
		"admin", "admin@gmail.com", adminPassword)

	if err != nil {
		err := tx.Rollback()
		panic(err.Error())
		//log.Fatal(err)
		return utils.Message(false, "Admin superuser cannot be added to countries table")
	}

	// Add admin to Admins Table
	id, _ := res.LastInsertId()
	_, err = tx.Exec("INSERT INTO admins(user_id, admin_role) VALUES (?,?)",
		id, "SUPER_USER")

	if err != nil {
		err := tx.Rollback()
		panic(err.Error())
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

func AdminExists(adminRole string) bool {
	db := InitDB()
	sqlStmt := "SELECT admin_role FROM admins WHERE admin_role=?"
	err := db.QueryRow(sqlStmt, adminRole).Scan(&adminRole)

	defer db.Close()

	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			log.Print(err)
		}

		return false
	}

	return true
}
