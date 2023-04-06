package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func ParseArrayToString(array []interface{}) string {
	/*
		This function is used to parse array to string using json marshal
	*/

	temp_params, err := json.Marshal(array)
	if err != nil {
		log.Println("Error in parsing array to string: ", err.Error())
		return ""
	}

	str := fmt.Sprintf("%v", string(temp_params))

	return str
}

func ParseInterfaceToString(inter interface{}) string {
	/*
		This function is used to parse interface to string using json marshal
	*/

	temp_params, err := json.Marshal(inter)
	if err != nil {
		log.Println("Error in parsing interface to string: ", err.Error())
		return ""
	}

	str := fmt.Sprintf("%v", string(temp_params))

	return str
}
