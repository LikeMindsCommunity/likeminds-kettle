package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ParseStringArrayToString(array []string) string {
	/*
		This function is used to parse String array to json string using json marshal
	*/

	temp_params, _ := json.Marshal(array)

	str := fmt.Sprintf("%v", string(temp_params))

	return str
}

func CheckIfStringExistsInArray(array []string, str string) bool {
	/*
		This function is used to check if string exists in array
	*/

	for _, a := range array {
		if a == str {
			return true
		}
	}

	return false
}

// GetBooleanFromString is used to get boolean from string | returns false by default
func GetBooleanFromString(str string) bool {
	return str == "true"
}

// Get first name from string
func GetFirstNameFromName(name string) string {
	var firstName string

	nameList := strings.Split(name, " ")

	if len(nameList) > 0 {
		firstName = nameList[0]
	}

	return firstName
}
