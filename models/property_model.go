package models

import (
		"time"
)

// ConstructionType handles only values: apartment, house
type Property struct {
		PropertyId          uint   `json:"property_id"`
		RoomNumber          int    `json:"room_number"`
		ConstructionType    string `json:"construction_type"`
		KidsAllowed         bool   `json:"kids_allowed"`
		PetsAllowed         bool   `json:"pets_allowed"`
		Area                int    `json:"area"`
		BathroomNumber      int    `json:"bathroom_number"`
		MaxFloorNumber      int    `json:"max_floor_number"`
		PropertyFloorNumber int    `json:"property_floor_number"`
}

// LisitngCurrency field holds by default such vars: usd, hrv, eur
type PropertyListing struct {
		PropertyId         uint       `json:"property_id"`
		ListingDescription string     `json:"listing_description"`
		ListingPrice       int        `json:"listing_price"`
		ListingCurrency    string     `json:"listing_currency"`
		ListingIsActive    bool       `json:"listing_is_active"`
		CreatedAt          time.Time  `json:"created_at"`
		UpdatedAt          *time.Time `json:"updated_at"`

		UserID uint `json:"user_id"`
		User   *User
		//AddressesID uint `json:"addresses_id"`
		//Addresses *Addresses
}

// Creates new PropertyListing
//func (lisitng PropertyListing) Create() map[string]interface{} {
//		// Validate if the property listing is the first in the table
//		// than create
//
//		// Create Property model
//		// Set the data to the database
//		// Create PropertyListing
//		// Set data to the database
//
//		// Validate the input data where the fields are strict necessary
//		// Handle errors
//		// Return response
//}

// Uses for grouping data of Property and PropertyListing
type NewPropertyResponse struct{}

//func (listing) Validate() {
//
//}
