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
