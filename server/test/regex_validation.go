package main

import (
	"fmt"
	"regexp"
)

func main() {
	Validate("099994")
	Validate("09Tgv99 hghgA94")
	ValidateEmail("apl.jakegmail.com")
	ValidateEmail("desg.natalia01@gmail.com")
}

func Validate(data string) {
	if m, err := regexp.MatchString("^[a-zA-Z0-9 ]+$", data); m {
		fmt.Println("valid")
	} else {
		panic(err.Error())
	}
}

func ValidateEmail(data string) {
	pattern := "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

	if m, _ := regexp.MatchString(pattern, data); m {
		fmt.Println("valid")
	} else {
		fmt.Println("Invalid")
	}
}
