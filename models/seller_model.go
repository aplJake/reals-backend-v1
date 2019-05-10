package models

import (
		"github.com/aplJake/reals-course/server/utils"
)

type Seller struct {
		ID              uint   `json:"user_id"`
		TelephoneNumber string `json:"telephone_number"`
}

func (seller Seller) Create() map[string]interface{} {
		if resp, ok := seller.Validate(); !ok {
				return resp
		}

		_, err := GetDb().Exec("INSERT INTO seller(user_id, telephone_number) VALUE (?,?)",
				seller.ID, seller.TelephoneNumber)
		if err != nil {
				return utils.Message(false, "Couldn`t insert a new seller")
		}
		return utils.Message(true, "A new seller was created")
}

func (seller Seller) Validate() (map[string]interface{}, bool) {
		if len(seller.TelephoneNumber) == 0 {
				return utils.Message(false, "Telephone number is required"), false
		}

		return utils.Message(true, "Requirements passed"), true
}
