package models

import (
		"database/sql"
		"fmt"
		"github.com/aplJake/reals-course/server/utils"
)

type Seller struct {
		ID              uint   `json:"user_id"`
		TelephoneNumber string `json:"telephone_number"`
}

func CreateSeller(id uint, phone string) (map[string]interface{}, bool) {
		//if resp, ok := seller.Validate(); !ok {
		//		return resp, false
		//}
		var db *sql.DB
		db = InitDB()

		fmt.Println(" Seller User id ", id, " phone ", phone)
		_, err := db.Exec("INSERT INTO seller(user_id, telephone_number) VALUE (?,?)",
				id, phone)
		defer db.Close()
		if err != nil {
				return utils.Message(false, "Couldn`t insert a new seller"), false
		}
		return utils.Message(true, "A new seller was created"), true
}

func (seller Seller) Validate() (map[string]interface{}, bool) {
		if len(seller.TelephoneNumber) == 0 {
				return utils.Message(false, "Telephone number is required"), false
		}

		return utils.Message(true, "Requirements passed"), true
}

func GetSeller(u uint) *Seller {
		seller := &Seller{}
		row := GetDb().QueryRow("SELECT * FROM seller WHERE user_id=?", u)
		err := row.Scan(&seller.ID, &seller.TelephoneNumber)
		if err != nil {
				return nil
		}
		return seller
}