package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const HeadersMemberId = "x-member-id"

//RequestParamsToMap Converts a struct to a map while maintaining the json alias as keys
func RequestParamsToMap(obj interface{}) (newMap map[string]string, err error) {
	data, err := json.Marshal(obj) // Convert to a json string
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}

//FormatFloat convert float upto prc
func FormatFloat(num float64, prc int) string {
	var (
		zero, dot = "0", "."
		str = fmt.Sprintf("%."+strconv.Itoa(prc)+"f", num)
	)
	return strings.TrimRight(strings.TrimRight(str, zero), dot)
}
