package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
)

const HeadersMemberId = "x-member-id"
const HeadersVersionCode = "x-version-code"
const HeadersPlatformCode = "x-platform-code"
const HeadersDeviceId = "x-device-id"
const HeadersApiKey = "x-api-key"

const GETMethod = 0
const POSTMethod = 1
const PUTMethod = 2
const DELETEMethod = 3

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
		str       = fmt.Sprintf("%."+strconv.Itoa(prc)+"f", num)
	)
	return strings.TrimRight(strings.TrimRight(str, zero), dot)
}

//CreateHeaders Used to create headers for our internal APIs
func CreateHeaders(c *gin.Context, userUniqueID string) map[string]interface{} {
	headers := make(map[string]interface{})
	if len(userUniqueID) > 0 {
		headers[HeadersMemberId] = userUniqueID
	}
	headers[HeadersPlatformCode] = c.GetHeader(HeadersPlatformCode)
	headers[HeadersVersionCode] = c.GetHeader(HeadersVersionCode)
	headers[HeadersDeviceId] = c.GetHeader(HeadersDeviceId)
	headers[HeadersApiKey] = c.GetHeader(HeadersApiKey)
	return headers
}

//Generate Response to be sent on request success
func GenerateResponse(c *gin.Context, dataResponse map[string]interface{}) {
	//
	response := Response{
		Success: true,
	}

	if len(dataResponse) > 0 {
		response.Data = dataResponse
	}

	c.JSON(http.StatusOK, response)
}

func ValidateClientResponse(c *gin.Context, respBytes []byte) *api_client.APIClientResponse {
	//Parse response
	var apiCR api_client.APIClientResponse
	err := api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
	if err != nil {
		//Internal unmarshal error
		GeneralAPIError(c, err.Error())
		return nil
	}

	if !apiCR.Success {
		//If internal api returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return nil
	}

	return &apiCR
}

//Parse Response from request sent internally
func ParseResponse(c *gin.Context, respBytes []byte) {

	apiCR := ValidateClientResponse(c, respBytes)

	if apiCR != nil {
		GenerateResponse(c, apiCR.Response)
	}
}
