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

func ParseInterfaceToString(data interface{}) string {
	/*
		This function is used to parse interface to string using json marshal
	*/

	// If inter is string type, then return it
	if _, ok := data.(string); ok {
		return data.(string)
	}

	temp_param, err := json.Marshal(data)
	if err != nil {
		log.Println("Error in parsing interface to string: ", err.Error())
		return ""
	}

	str := fmt.Sprintf("%v", string(temp_param))

	return str
}
