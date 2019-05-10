package models

import (
		"fmt"
		"github.com/aplJake/reals-course/server/utils"
		"log"
		"time"
)

// ConstructionType handles only values: apartment, house
type Property struct {
		PropertyId          uint   `json:"property_id"`
		RoomNumber          int    `json:"room_number,string"`
		ConstructionType    string `json:"construction_type"`
		KidsAllowed         *bool   `json:"kids_allowed"`
		PetsAllowed         *bool   `json:"pets_allowed"`
		Area                int    `json:"area,string"`
		BathroomNumber      int    `json:"bathroom_number,string"`
		MaxFloorNumber      int    `json:"max_floor_number,string"`
		PropertyFloorNumber int    `json:"property_floor_number,string"`
}

// LisitngCurrency field holds by default such vars: usd, hrv, eur
type PropertyListing struct {
		PropertyId         uint   `json:"property_id"`
		ListingDescription string `json:"listing_description"`
		ListingPrice       int    `json:"listing_price,string"`
		ListingCurrency    string `json:"listing_currency"`
		ListingIsActive    *bool   `json:"listing_is_active"`
		// Automation
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`

		UserID uint `json:"user_id"`
		User   *User
		//AddressesID uint `json:"addresses_id"`
		//Addresses *Addresses
}

// Creates new PropertyListing
func CreateListing(listing *PropertyListingRequest) map[string]interface{} {
		//seller := &Seller{}
		//seller.ID = listing.UserId
		fmt.Println("Listing id", listing.UserId)
		//fmt.Println("Seller id", seller.ID)
		seller := GetSeller(listing.UserId)
		// Validate Seller Account
		fmt.Println("Seller fidede ", seller)
		if seller == nil {
				// CreateSeller the SellerAccount
				if resp, ok := CreateSeller(listing.UserId, ""); !ok {
						return resp
				}
		}

		// Validate if the property listing is the first in the table
		//if r, ok := ListingValidate(listing.UserID); !ok {
		//		return r
		//}

		// CreateSeller new Property and property listing with transaction
		fmt.Println("Property listing", listing)

		_, err := GetDb().Exec(`INSERT INTO property(
                    room_number, construction_type, kids_allowed, pets_allowed,
                    area, bathroom_number, max_floor_number, property_floor_number)
                     VALUES(?,?,?,?,?,?,?,?);`, listing.RoomNumber, listing.ConstructionType,
				listing.KidsAllowed, listing.PetsAllowed, listing.Area,
				listing.BathroomNumber, listing.MaxFloorNumber, listing.PropertyFloorNumber)
		if err != nil {
				log.Fatal(err)
		}

		tx, err := GetDb().Begin()
		fmt.Println("Property listing 2")

		if r := handleError(err); r != nil {
				return r
		}
		fmt.Println("Property listing 3")
		fmt.Println(listing.RoomNumber, listing.ConstructionType,
				listing.KidsAllowed, listing.PetsAllowed, listing.Area,
				listing.BathroomNumber, listing.MaxFloorNumber, listing.PropertyFloorNumber)

		// Insert data to the property table
		res, err := tx.Exec(`INSERT INTO property(
                     room_number, construction_type, kids_allowed, pets_allowed,
                     area, bathroom_number, max_floor_number, property_floor_number)
                      VALUES(?,?,?,?,?,?,?,?);`, listing.RoomNumber, listing.ConstructionType,
                      listing.KidsAllowed, listing.PetsAllowed, listing.Area,
                      listing.BathroomNumber, listing.MaxFloorNumber, listing.PropertyFloorNumber)

		//fmt.Println(err.Error())

		if err != nil {
				err := tx.Rollback()
				log.Fatal(err)
				return utils.Message(false, "Property object wasn created")
		}

		fmt.Println("Property listing 5")

		// Fetch the auto increment Property Id
		id, err := res.LastInsertId()
		if r := handleError(err); r != nil {
				return r
		}
		fmt.Println("Index 23 ", id)
		// Insert data to Property Listing
		res, err = tx.Exec(`INSERT INTO property_listing(
                             property_id, user_id, listing_description,
                             listing_price, listing_currency, listing_is_active) VALUES(?,?,?,?,?,?);`,
                             id, seller.ID, listing.ListingDescription, listing.ListingPrice,
                             listing.ListingCurrency, listing.ListingIsActive)
		if err != nil {
				tx.Rollback()
				log.Fatal(err)
				return utils.Message(false, "Property object wasn created")
		}

		// commit the transaction
		if r := handleError(tx.Commit()); r != nil {
				return r
		}

		log.Println("Done added listing")

		return utils.Message(true, "Property listing was successfully added")

		// CreateSeller Property model
		// Set the data to the database
		// CreateSeller PropertyListing
		// Set data to the database

		// Validate the input data where the fields are strict necessary
		// Handle errors
		// Return response
}

func handleError(err error) map[string]interface{} {
		if err != nil {
				log.Fatal(err)
				return utils.Message(false, "Error while addig a new listing")
		}
		return nil
}

// Validate if such seller with UserId exists
func ListingValidate(u uint) (map[string]interface{}, bool) {
		if seller := GetSeller(u); seller == nil {
				r := utils.Message(false, "There is no seller account")
				return r, false
		}

		r := utils.Message(true, "Lising is validated")
		return r, true
}

//Response serializer
// Uses for grouping data of Property and PropertyListing
type PropertyListingRequest struct {
		PropertyId          uint   `json:"property_id"`
		UserId              uint   `json:"user_id"`
		ConstructionType    string `json:"construction_type"`
		Area                int `json:"area,string"`
		RoomNumber          int    `json:"room_number,string"`
		BathroomNumber      int    `json:"bathroom_number,string"`
		MaxFloorNumber      string    `json:"max_floor_number"`
		PropertyFloorNumber string    `json:"property_floor_number"`
		KidsAllowed         *bool   `json:"kids_allowed"`
		PetsAllowed         *bool   `json:"pets_allowed"`
		// Listing
		ListingDescription string `json:"listing_description"`
		ListingPrice       string    `json:"listing_price"`
		ListingCurrency    string `json:"listing_currency"`
		ListingIsActive    *bool   `json:"listing_is_active"`
}

//func (listing) Validate() {
//
//}
