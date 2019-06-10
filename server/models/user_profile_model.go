package models

import (
		"database/sql"
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

//func (user *User) InitProfile() error {
//		db := InitDB()
//		fmt.Println("User id is ", user.ID)
//		_, err := db.Exec("INSERT INTO user_profile(user_id, profile_description) VALUE(?, ?)", user.ID, "")
//
//		return err
//}

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
		CreatedAt          *time.Time  `json:"created_at"`
		UpdatedAt          *time.Time `json:"updated_at"`
}

var _ = `
SELECT u.user_id, u.user_name, p.profile_description, p.created_at, p.updated_at
				FROM user_profile p 
				JOIN listings u ON p.user_id = u.user_id
				WHERE u.user_id = ?;
`
var getUserProfileQ = `
		SELECT 
		u.user_id,
		u.user_name,
		p.profile_description,
		cast(p.created_at as datetime),
		cast(p.updated_at as datetime)
		FROM user_profile p 
			JOIN users u 
			    ON p.user_id = u.user_id
		WHERE u.user_id = ?;
`
func GetProfileData(u uint) (UserProfileRespond, error) {
		var profileRes UserProfileRespond

		db := InitDB()
		row := db.QueryRow(getUserProfileQ, u)

		err := row.Scan(&profileRes.UserID, &profileRes.UserName,
				&profileRes.ProfileDescription, &profileRes.CreatedAt, &profileRes.UpdatedAt)
		if err != nil {
				panic(err.Error())
		}
		defer db.Close()

		return profileRes, err
}

var getListingsByProfileQ = `
	SELECT property_id,
	       user_id,
	       addresses_id,
	       listing_description,
	       listing_price,
	       listing_currency,
	       listing_is_active,
	       cast(created_at as datetime),
		   cast(updated_at as datetime)
	FROM property_listing
	WHERE user_id=?;
`

func GetLisitingsByProfile(profileID uint) ([]PropertyListing, error)  {
		var db *sql.DB
		var err error

		db = InitDB()

		res, err := db.Query(getListingsByProfileQ, profileID)
		handleError(err)

		listing := PropertyListing{}
		listingsArr := []PropertyListing{}
		for res.Next() {
				err = res.Scan(&PropertyId, &UserID, &AddressesID, &ListingDescription,
						&ListingPrice, &ListingCurrency, &ListingIsActive, &CreatedAt, &UpdatedAt)
				handleError(err)
				listingsArr = append(listingsArr, listing)
		}

		defer db.Close()

		return listingsArr, err
}
// Response function to the handler
func NewUserProfileResponse(profile UserProfileRespond, listings []PropertyListing) (map[string]interface{}, error) {
		// TODO: validate the article in ArticleResponse (check user by userId and so on)
		response := utils.Message(true, "Profile data was accessed")
		response["profile"] = profile
		response["listings"] = listings
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
