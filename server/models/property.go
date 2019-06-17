package models

import (
	"database/sql"
	"errors"
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

	tx, err := db.Begin()
	fmt.Println("Property listing 1")

	if r := handleError(err); r != nil {
		return r
	}
	fmt.Println("Property listing 2")
	fmt.Println(listing.RoomNumber, listing.ConstructionType,
		listing.Area, listing.BathroomNumber, listing.MaxFloorNumber, listing.PropertyFloorNumber)

	// Insert data to the property table
	var insertNewPropertyQ = `
		INSERT INTO property(
		room_number,
		construction_type,
		area,
		bathroom_number,
		max_floor_number,
		property_floor_number
		) VALUES(?,?,?,?,?,?);
	`
	res, err := tx.Exec(insertNewPropertyQ, listing.RoomNumber, listing.ConstructionType,
		listing.Area, listing.BathroomNumber, listing.MaxFloorNumber, listing.PropertyFloorNumber)

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
	res2, err := tx.Exec(insertNewStreetQ, listing.CityID, listing.StreetName,
		listing.StreetNumber)
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

func (l *PropertyListingRequest) UpdateListing() error  {
	db := InitDB()

	tx, err := db.Begin()
	if err != nil {
		return errors.New("DB Transaction begin error")
	}

	var updatePropertyQ = `
		UPDATE property
		SET room_number=?,
		    construction_type=?,
		    area=?,
		    bathroom_number=?,
		    max_floor_number=?,
		    property_floor_number=?
		WHERE property_id=?;
		    
	`
	_, err = tx.Exec(updatePropertyQ, l.RoomNumber, l.ConstructionType, l.Area,
		l.BathroomNumber, l.MaxFloorNumber, l.PropertyFloorNumber, l.PropertyId)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.New("Property update transaction error")
		}
		return errors.New("Property update error")
	}

	var updatePropertyStreetQ = `
		UPDATE addresses
		SET city_id=?,
		    street_name=?,
		    street_number=?
		WHERE addresses_id=?
	`
	_, err = tx.Exec(updatePropertyStreetQ, l.CityID, l.StreetNumber, l.StreetNumber,
		l.AddressID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.New("Addresses update transaction error")
		}
		return errors.New("Addresses update error")
	}

	var updatePropertyListing = `
		UPDATE property_listing
		SET property_id=?, 
		    user_id=?, 
		    addresses_id=?, 
		    listing_description=?, 
		    listing_price=?, 
		    listing_currency=?, 
		    listing_is_active=?
		WHERE property_id=?
	`
	_, err = tx.Exec(updatePropertyListing, l.PropertyId, l.UserId, l.AddressID,
		l.ListingDescription, l.ListingPrice, l.ListingCurrency,
		l.ListingIsActive, l.PropertyId)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.New("Listing update transaction error")
		}
		//panic(err.Error())
		return errors.New("Listing update error")
	}

	if err = tx.Commit(); err != nil {
		return errors.New("Property listing transacion commit error")
	}

	defer db.Close()
	return nil
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

func GetListingsList() ([]PropertyListing, error) {
	fmt.Println("Getting all listings");

	db := InitDB()

	res, err := db.Query(getAllListings)
	handleError(err)

	defer db.Close()


	var listing PropertyListing
	var listingsArr []PropertyListing
	for res.Next() {
		err = res.Scan(&listing.PropertyId, &listing.UserID, &listing.AddressesID, &listing.ListingDescription,
			&listing.ListingPrice, &listing.ListingCurrency, &listing.ListingIsActive, &listing.CreatedAt, &listing.UpdatedAt)
		handleError(err)
		listingsArr = append(listingsArr, listing)
		fmt.Println(listingsArr)
	}

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
	PropertyId          uint   `json:"property_id"`
	UserId              uint   `json:"user_id"`
	ConstructionType    string `json:"construction_type"`
	Area                int    `json:"area,string"`
	RoomNumber          string    `json:"room_number"`
	BathroomNumber      int    `json:"bathroom_number,string"`
	MaxFloorNumber      string `json:"max_floor_number"`
	PropertyFloorNumber string `json:"property_floor_number"`
	// Listing
	ListingDescription string `json:"listing_description"`
	ListingPrice       string `json:"listing_price"`
	ListingCurrency    string `json:"listing_currency"`
	ListingIsActive    bool   `json:"listing_is_active,string"`
	// Address
	CityID		uint	`json:"city_id,string"`
	StreetName   string `json:"street_name"`
	StreetNumber string `json:"street_number"`
	CountryID	uint	`json:"country_id,string"`
	AddressID	uint	`json:"address_id,string"`
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
		&p.Property.RoomNumber, &p.Property.ConstructionType, &p.Property.Area,
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
		from property_listing L
			inner join seller S on L.user_id = S.user_id
			inner join user_profile UP on S.user_id = UP.user_id
			inner join users U on UP.user_id = U.user_id
		where L.property_id = ?;
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

var getListingDataForUpdateQ = `
	SELECT L.user_id,
	       L.listing_description,
	       L.listing_price,
	       L.listing_currency,
	       L.listing_is_active,
	       P.*,
	       C.country_id,
	       A.city_id, A.street_name, A.street_number, A.addresses_id
	FROM property_listing L
		 INNER JOIN property P
			    ON P.property_id = L.property_id
		 INNER JOIN addresses A
			    ON A.addresses_id = L.addresses_id
		 INNER JOIN city C
			    ON A.city_id = C.city_id
	WHERE P.property_id=?;
`
func GetPropertyListing(propertyID string) (*PropertyListingRequest, error)  {
	fmt.Println("Property id 1", propertyID)
	p := &PropertyListingRequest{}
	db := InitDB()
	// We get multiple data from the db
	// For that purpose we use transaction
	res := db.QueryRow(getListingDataForUpdateQ, propertyID)
	err := res.Scan(&p.UserId, &p.ListingDescription, &p.ListingPrice, &p.ListingCurrency, &p.ListingIsActive, &p.PropertyId,
		&p.RoomNumber, &p.ConstructionType, &p.Area, &p.BathroomNumber, &p.MaxFloorNumber,
		&p.PropertyFloorNumber, &p.CountryID, &p.CityID, &p.StreetName, &p.StreetNumber, &p.AddressID)

	if err != nil {
		return nil, errors.New("Error to get proeprty listing from db")
	}

	defer db.Close()
	return p, nil
}

