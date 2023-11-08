package utils

import (
	"encoding/json"
	"fmt"

	log "github.com/nateshr/likeminds-authentication/logging"
)

// This function is used to parse array to string using json marshal
func ParseArrayToString(array []interface{}) string {

	temp_params, err := json.Marshal(array)
	if err != nil {

		log.Error(fmt.Sprintf("Error in parsing array to string: %s", err.Error()))
		return ""
	}

	str := fmt.Sprintf("%v", string(temp_params))

	return str
}

// This function is used to parse interface to string using json marshal
func ParseInterfaceToString(data interface{}) string {

	if data == nil {
		return ""
	}

	// If data is string type, then return it
	if _, ok := data.(string); ok {
		return data.(string)
	}

	temp_param, err := json.Marshal(data)
	if err != nil {
		log.Error(fmt.Sprintf("Error in parsing interface to string: %s", err.Error()))
		return ""
	}

	str := fmt.Sprintf("%v", string(temp_param))

	return str
}
