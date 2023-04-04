package utils

import (
	"encoding/json"
	"fmt"
)

func ParseArrayToString(array []interface{}) string {
	/*
		This function is used to parse array to string using json marshal
	*/

	temp_params, _ := json.Marshal(array)

	str := fmt.Sprintf("%v", string(temp_params))

	return str
}

func ParseStringArrayToString(array []string) string {
	/*
		This function is used to parse String array to json string using json marshal
	*/

	temp_params, _ := json.Marshal(array)

	str := fmt.Sprintf("%v", string(temp_params))

	return str
}
