package models

import (
		"fmt"
		"github.com/aplJake/reals-course/server/utils"
		"time"
)

type UserProfile struct {
		ProfileDescription string     `json:"profile_description"`
		CreatedAt          time.Time  `json:"created_at"`
		UpdatedAt          *time.Time `json:"updated_at"`
		UserID             uint       `json:"user_id"`
		User               *User
}

func (user *User) InitProfile() error {
		db := GetDb()
		fmt.Println("User id is ", user.ID)
		_, err := db.Exec("INSERT INTO user_profile(user_id) VALUE(?)", user.ID)
		return err
}

func (profile *UserProfile) Update(data UserProfile) error {
		db := GetDb()
		_, err := db.Exec("UPDATE user_profile SET profile_description=?", data.ProfileDescription)
		return err
}

// Uses to get the response about the received object
func (profile *UserProfile) ReceiveProfile() map[string]interface{} {
		// TODO: validate the received data ?

		dbProfile, err := GetDb().Exec("SELECT * FROM user_profile WHERE user_id=?", profile.UserID)
		if err != nil {
				fmt.Println(err.Error())
				panic(err.Error())
				return utils.Message(false, "Failed to get profile data, data error.")
		}

		response := utils.Message(true, "Profile data was accessed")
		response["profile"] = dbProfile
		return response
}

// Used for profile data
type UserProfileRespond struct {
		UserID      uint       `json:"user_id"`
		UserName           string     `json:"user_name"`
		ProfileDescription string     `json:"profile_description"`
		CreatedAt          time.Time  `json:"created_at"`
		UpdatedAt          *time.Time `json:"updated_at"`
}

func GetUserProfile(u uint) (*UserProfileRespond, error) {
		resProfile := &UserProfileRespond{}
		//fmt.Println("User id23 ", u)
		//row := GetDb().QueryRow("SELECT profile_id, profile_description, created_at FROM user_profile WHERE user_id=?", u)
		row := GetDb().QueryRow(
				`	SELECT u.user_id, u.user_name, p.profile_description, p.created_at, p.updated_at
				FROM user_profile p 
				JOIN listings u ON p.user_id = u.user_id
				WHERE u.user_id = ?;`, u)

		err := row.Scan(&resProfile.UserID, &resProfile.UserName,
				&resProfile.ProfileDescription, &resProfile.CreatedAt, &resProfile.UpdatedAt)
		return resProfile, err
}

// Response function to the handler
func NewUserProfileResponse(profile *UserProfileRespond) (map[string]interface{}, error) {
		// TODO: validate the article in ArticleResponse (check user by userId and so on)
		response := utils.Message(true, "Profile data was accessed")
		response["profile"] = profile
		return response, nil
}

// Response function to the handler
func UpdateUserProfileResponse(profile *UserProfile) (map[string]interface{}, error) {
		fmt.Printf("\nUser id in profile %s\n", profile.UserID)
		fmt.Printf("\nProfile description %s\n", profile.ProfileDescription)

		_, err := db.Exec("UPDATE user_profile SET profile_description=? WHERE user_id=?", profile.ProfileDescription, profile.UserID)
		respond := utils.Message(true, "Profile data was updated.")
		return respond, err

}
