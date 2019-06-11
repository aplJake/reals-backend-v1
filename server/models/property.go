package models

import (
	"database/sql"
	"fmt"
	"github.com/aplJake/reals-course/server/utils"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

// ConstructionType handles only values: apartment, house
type Property struct {
	PropertyId          uint   `json:"property_id"`
	RoomNumber          string `json:"room_number,string"`
	ConstructionType    string `json:"construction_type"`
	KidsAllowed         bool   `json:"kids_allowed"`
	PetsAllowed         bool   `json:"pets_allowed"`
	Area                int    `json:"area,string"`
	BathroomNumber      int    `json:"bathroom_number,string"`
	MaxFloorNumber      int    `json:"max_floor_number,string"`
	PropertyFloorNumber int    `json:"property_floor_number,string"`
}

// LisitngCurrency field holds by default such vars: usd, hrv, eur
type PropertyListing struct {
	PropertyId         uint      `json:"property_id"`
	AddressesID        uint      `json:"addresses_id"`
	UserID             uint      `json:"user_id"`
	ListingDescription string    `json:"listing_description"`
	ListingPrice       int       `json:"listing_price,string"`
	ListingCurrency    string    `json:"listing_currency"`
	ListingIsActive    bool      `json:"listing_is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	//Addresses *Addresses
}

// Creates new PropertyListing
func CreateListing(listing *PropertyListingRequest) map[string]interface{} {
	var db *sql.DB

	db = InitDB()
	//seller := &Seller{}
	//seller.ID = listing.UserId
	fmt.Println("Listing id", listing.UserId)
	//fmt.Println("Seller id", seller.ID)

	exists, err := SellerIsExists(listing.UserId)
	if err != nil {
		panic(err.Error())
	}

	seller := &Seller{}
	if !exists {
		resp, seller := CreateSeller(listing.UserId, "")
		if seller == nil {
			return resp
		}
	} else {
		seller = GetSeller(listing.UserId)
	}

	//seller := GetSeller(listing.UserId)
	// Validate Seller Account

	//fmt.Println("Seller fidede ", seller)
	//if seller == nil {
	//		// CreateSeller the SellerAccount
	//		if resp, ok := CreateSeller(listing.UserId, ""); !ok {
	//				return resp
	//		}
	//}

	// TODO: ADD PROPERTY VALIDATION
	// Validate if the property listing is the first in the table
	//if r, ok := ListingValidate(listing.UserID); !ok {
	//		return r
	//}

	// CreateSeller new Property and property listing with transaction

	//fmt.Println("Property listing", listing)
	//
	//_, err := db.Exec(`INSERT INTO property(
	//            room_number, construction_type, kids_allowed, pets_allowed,
	//            area, bathroom_number, max_floor_number, property_floor_number)
	//             VALUES(?,?,?,?,?,?,?,?);`, listing.RoomNumber, listing.ConstructionType,
	//		listing.KidsAllowed, listing.PetsAllowed, listing.Area,
	//		listing.BathroomNumber, listing.MaxFloorNumber, listing.PropertyFloorNumber)
	//if err != nil {
	//		log.Fatal(err)
	//}

	tx, err := db.Begin()
	fmt.Println("Property listing 1")

	if r := handleError(err); r != nil {
		return r
	}
	fmt.Println("Property listing 2")
	fmt.Println(listing.RoomNumber, listing.ConstructionType,
		listing.KidsAllowed, listing.PetsAllowed, listing.Area,
		listing.BathroomNumber, listing.MaxFloorNumber, listing.PropertyFloorNumber)

	// Insert data to the property table
	var insertNewPropertyQ = `
				INSERT INTO property(
                room_number,
				construction_type,
				kids_allowed,
				pets_allowed,
                area,
				bathroom_number,
				max_floor_number,
				property_floor_number
				) VALUES(?,?,?,?,?,?,?,?);
		`
	res, err := tx.Exec(insertNewPropertyQ, listing.RoomNumber, listing.ConstructionType,
		listing.KidsAllowed, listing.PetsAllowed, listing.Area,
		listing.BathroomNumber, listing.MaxFloorNumber, listing.PropertyFloorNumber)

	//fmt.Println(err.Error())

	if err != nil {
		panic(err.Error())
		if err := tx.Rollback(); err != nil {
			panic(err.Error())
		}
		return utils.Message(false, "Property object wasn created")
	}

	//TODO: add transaction for adress adding
	var insertNewStreetQ = `
		INSERT INTO addresses(
				city_id,
		        street_name,
		        street_number
		) VALUES (?,?,?);
		`
	res2, err := tx.Exec(insertNewStreetQ, listing.AddressesRequest.CityID, listing.AddressesRequest.StreetName,
		listing.AddressesRequest.StreetNumber)
	if err != nil {
		panic(err.Error())
		if err := tx.Rollback(); err != nil {
			panic(err.Error())
		}
		return utils.Message(false, "Property object wasn created")
	}

	countryId, err := res2.LastInsertId()
	if r := handleError(err); r != nil {
		return r
	}
	fmt.Println("Property listing 5")

	// Fetch the auto increment Property Id
	id, err := res.LastInsertId()
	if r := handleError(err); r != nil {
		return r
	}
	fmt.Println("Index 23 ", id)

	//var _  = `INSERT INTO property_listing(
	//                     property_id, user_id, listing_description,
	//                     listing_price, listing_currency, listing_is_active) VALUES(?,?,?,?,?,?);`
	var insertNewListingQ = `
				INSERT INTO property_listing(
						property_id, 
				        user_id, 
				        addresses_id, 
				        listing_description, 
				        listing_price, 
				    	listing_currency, 
						listing_is_active
				) VALUES (?,?,?,?,?,?,?);
		`
	fmt.Println("Error code", id, seller.ID, countryId, listing.ListingDescription,
		listing.ListingPrice, listing.ListingCurrency, listing.ListingIsActive)
	// Insert data to Property Listing
	res, err = tx.Exec(insertNewListingQ, id, seller.ID, countryId, listing.ListingDescription,
		listing.ListingPrice, listing.ListingCurrency, listing.ListingIsActive)
	if err != nil {
		panic(err.Error())
		if err := tx.Rollback(); err != nil {
			panic(err.Error())
		}
		return utils.Message(false, "Property object wasn created")
	}

	// commit the transaction
	if r := handleError(tx.Commit()); r != nil {
		return r
	}

	defer db.Close()
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
		panic(err.Error())
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

var getAllListings = `
	SELECT * FROM property_listing;
`

func GetAllListings() ([]PropertyListing, error) {
	db := InitDB()

	res, err := db.Query(getAllListings)
	handleError(err)

	var listing PropertyListing
	var listingsArr []PropertyListing
	for res.Next() {
		err = res.Scan(&listing.PropertyId, &listing.UserID, &listing.AddressesID, &listing.ListingDescription,
			&listing.ListingPrice, &listing.ListingCurrency, &listing.ListingIsActive, &listing.CreatedAt, &listing.UpdatedAt)
		handleError(err)
		listingsArr = append(listingsArr, listing)
	}
	defer db.Close()

	return listingsArr, nil
}

var getListingsByTypeQ = `
	SELECT L.*
	FROM property_listing L
		INNER JOIN property P 
		ON P.property_id = L.property_id
	WHERE P.construction_type=?;
`

func GetListingsByType(propertyType string) ([]PropertyListing, error) {
	db := InitDB()

	res, err := db.Query(getListingsByTypeQ, propertyType)
	handleError(err)

	var listing PropertyListing
	var listingsArr []PropertyListing
	for res.Next() {
		err = res.Scan(&listing.PropertyId, &listing.UserID, &listing.AddressesID, &listing.ListingDescription,
			&listing.ListingPrice, &listing.ListingCurrency, &listing.ListingIsActive, &listing.CreatedAt, &listing.UpdatedAt)
		handleError(err)
		listingsArr = append(listingsArr, listing)
	}
	defer db.Close()

	return listingsArr, nil
}

//Response serializer
// Uses for grouping data of Property and PropertyListing
type PropertyListingRequest struct {
	PropertyId          uint   `db:"property_id" json:"property_id"`
	UserId              uint   `db:"user_id"json:"user_id"`
	ConstructionType    string `db:"construction_type" json:"construction_type"`
	Area                int    `db:"area" json:"area,string"`
	RoomNumber          int    `db:"room_number" json:"room_number,string"`
	BathroomNumber      int    `db:"bathroom_number" json:"bathroom_number,string"`
	MaxFloorNumber      string `db:"max_floor_number" json:"max_floor_number"`
	PropertyFloorNumber string `db:"property_floor_number" json:"property_floor_number"`
	KidsAllowed         bool   `db:"kids_allowed" json:"kids_allowed,string"`
	PetsAllowed         bool   `db:"pets_allowed" json:"pets_allowed,string"`
	// Listing
	ListingDescription string `db:"listing_description" json:"listing_description"`
	ListingPrice       string `db:"listing_price" json:"listing_price"`
	ListingCurrency    string `db:"listing_currency" json:"listing_currency"`
	ListingIsActive    bool   `db: listing_is_active" json:"listing_is_active,string"`
	// Address
	AddressesRequest *AddressesRequest `json:"addresses"`
}

type PropertyPageData struct {
	Listing  PropertyListing `json:"property_listing"`
	Property Property        `json:"property"`
	Address  Addresses       `json:"address"`
}

// Transaction utils
type DB struct {
	*sql.DB
}

type Tx struct {
	*sql.Tx
}

func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

var getPropertyListingDataQ = `
		SELECT L.user_id,
			   L.addresses_id,
			   L.listing_description,
			   L.listing_price,
			   L.listing_currency,
			   L.listing_is_active,
			   L.created_at,
			   L.updated_at,
			   P.*,
			   A.city_id, A.street_name, A.street_number
		FROM property_listing L
				 inner JOIN property P
							ON P.property_id = L.property_id
				 inner JOIN addresses A
							ON A.addresses_id = L.addresses_id
		WHERE P.property_id=?;
`

func (tx *Tx) GetPropertyListing(propertyID string) (*PropertyPageData, error) {
	listing := &PropertyPageData{}
	res := tx.QueryRow(getPropertyListingDataQ, propertyID)
	err := res.Scan(&listing)
	return listing, err
}

func GetPropertyPageData(propertyID string) (PropertyPageData, error) {
	p := PropertyPageData{}
	db := InitDB()
	// We get multiple data from the db
	// For that purpose we use transaction
	res := db.QueryRow(getPropertyListingDataQ, propertyID)
	err := res.Scan(&p.Listing.UserID, &p.Address.AddressesId, &p.Listing.ListingDescription, &p.Listing.ListingPrice, &p.Listing.ListingCurrency, &p.Listing.ListingIsActive, &p.Listing.CreatedAt, &p.Listing.UpdatedAt, &p.Property.PropertyId,
		&p.Property.RoomNumber, &p.Property.ConstructionType, &p.Property.KidsAllowed, &p.Property.PetsAllowed, &p.Property.Area,
		&p.Property.BathroomNumber, &p.Property.MaxFloorNumber, &p.Property.PropertyFloorNumber, &p.Address.CityId, &p.Address.StreetName, &p.Address.StreetNumber)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	return p, nil
}

// Consists from UserProfile data and QueueData
type PropertyProfileData struct {
	ProfileDescription sql.NullString `json:"profile_description"`
	TelephoneNumber    sql.NullString `json:"telephone_number"`
	Username           string         `json:"user_name"`
	UserEmail          string         `json:"user_email"`
}

type PropertyQueueData struct {
	UserName  string    `json:"user_name"`
	QueueTime time.Time `json:"queue_time"`
}

type PropertyCtxData struct {
	Queue   []PropertyQueueData `json:"queue"`
	Profile PropertyProfileData `json:"profile"`
}

var getPropertyQueueDataQ = `
		select U.user_name,
			   Q.queue_time
		from property_queue Q
				 inner join users U on Q.user_id = U.user_id
		where Q.property_id = ?;
`

func GetProperyQueue(propertyID string) ([]PropertyQueueData, error) {
	db := InitDB()

	res, err := db.Query(getPropertyQueueDataQ, propertyID)
	handleError(err)

	var qData PropertyQueueData
	var qDataArr []PropertyQueueData

	defer db.Close()
	for res.Next() {
		err = res.Scan(&qData.UserName, &qData.QueueTime)
		handleError(err)
		qDataArr = append(qDataArr, qData)
	}
	return qDataArr, err
}

var getPropertyProfileDataQ = `
		select UP.profile_description,
			   S.telephone_number,
			   U.user_name,
			   U.user_email
		from property_queue Q
			inner join seller S on Q.user_id = S.user_id
			inner join user_profile UP on S.user_id = UP.user_id
			inner join users U on UP.user_id = U.user_id
		where Q.property_id = ?;
`

func GetPropertyProfileData(propertyID string) (PropertyProfileData, error) {
	var pData PropertyProfileData
	db := InitDB()

	res := db.QueryRow(getPropertyProfileDataQ, propertyID)
	err := res.Scan(&pData.ProfileDescription, &pData.TelephoneNumber,
		&pData.Username, &pData.UserEmail)
	handleError(err)
	defer db.Close()

	return pData, err
}