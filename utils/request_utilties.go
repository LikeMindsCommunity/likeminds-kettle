package utils

import (
	"encoding/json"
	"fmt"
	"strings"

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

// Exposed Method to parse a String Array from query params
func ParseStringArrayFromParam(param string) []string {
	response := []string{}

	// Removal of square braces from array string
	if len(param) > 0 && param[0] == '[' {
		param = param[1:]
	}

	if len(param) > 0 && param[len(param)-1] == ']' {
		param = param[:len(param)-1]
	}

	// Removal of extra spaces from the array string
	param = strings.TrimSpace(param)

	if len(param) > 0 {
		paramValues := strings.Split(param, ",")

		for _, value := range paramValues {
			value = strings.TrimSpace(value)

			// Removal of quotes from each string from array
			if len(value) > 0 && value[0] == '"' {
				value = value[1:]
			}

			if len(value) > 0 && value[len(value)-1] == '"' {
				value = value[:len(value)-1]
			}

			response = append(response, value)
		}
	}

	return response
}
