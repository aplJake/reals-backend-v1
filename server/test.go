package server

import (
		"encoding/json"
		"fmt"
)

type City struct {
		CityName  string `json:"city_id"`
}

type Addresses struct {
		StreetName   string `json:"street_name"`
		StreetNumber string `json:"street_number"`
		City         *City 	`json:"city_struct"`
}

func main() {
		var address Addresses
		m := []byte(`{"street_name": "Puchkinska", "street_number": "23", "city_struct ": {"city_id": "Kharkiv"} }`)
		err := json.Unmarshal(m, &address)
		if err != nil {
				panic(err.Error())
		}

		fmt.Println(address, address.City)
}
