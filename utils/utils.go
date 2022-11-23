package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
)

//RequestParamsToMap Converts a struct to a map while maintaining the json alias as keys
func RequestParamsToMap(obj interface{}) (newMap map[string]string, err error) {
	data, err := json.Marshal(obj) // Convert to a json string
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
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
	headers[HeadersAcceptVersion] = c.GetHeader(HeadersAcceptVersion)
	return headers
}

//Generate Response to be sent on request success
func GenerateResponse(c *gin.Context, dataResponse map[string]interface{}) {
	//Generating Response Object
	response := Response{
		Success: true,
	}

	//Removing Blank Data Key
	if len(dataResponse) > 0 {
		response.Data = dataResponse
	}

	c.JSON(http.StatusOK, response)
}

func ValidateClientResponse(c *gin.Context, respBytes []byte, statusCode int) *api_client.APIClientResponse {
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
		c.JSON(statusCode, apiCR)
		return nil
	}

	return &apiCR
}

//ParseResponse from request sent internally
func ParseResponse(c *gin.Context, respBytes []byte, statusCode int) {

	apiCR := ValidateClientResponse(c, respBytes, statusCode)

	if apiCR != nil {
		GenerateResponse(c, apiCR.Response)
	}
}

func GetRequestResponse(c *gin.Context, serviceType ServiceType, url string, requestType RequestType, headers map[string]interface{}, params map[string]string, body interface{}) ([]byte, int) {
	//Create internal API client
	client := api_client.NewAPIClient()
	var baseUrl string
	var respBytes []byte
	var statusCode int
	var err error

	switch serviceType {
	case CoreService:
		baseUrl = client.CoreServiceBaseURL

	case SubscriptionService:
		baseUrl = client.SubscriptionServiceBaseURL
	}

	switch requestType {
	case GETRequest:

		options := api_client.GetRequestOptions{
			Url:           baseUrl + url,
			CustomHeaders: headers,
			Params:        params,
		}

		respBytes, statusCode, err = client.GetRequest(&options)

	case POSTRequestRawBody:

		options := api_client.PostRequestOptions{
			Url:           baseUrl + url,
			CustomHeaders: headers,
			Params:        params,
			Body:          body,
		}

		respBytes, statusCode, err = client.PostRequest(&options, api_client.BodyTypeRaw)

	case POSTRequestFormUrlEncodedBody:

		options := api_client.PostRequestOptions{
			Url:           baseUrl + url,
			CustomHeaders: headers,
			Params:        params,
			Body:          body,
		}

		respBytes, statusCode, err = client.PostRequest(&options, api_client.BodyTypeFormUrlEncoded)

	case PUTRequest:

		options := api_client.PostRequestOptions{
			Url:           baseUrl + url,
			CustomHeaders: headers,
			Params:        params,
			Body:          body,
		}

		respBytes, statusCode, err = client.PutRequest(&options)

	case DELETERequest:

		options := api_client.PostRequestOptions{
			Url:           baseUrl + url,
			CustomHeaders: headers,
			Params:        params,
			Body:          body,
		}

		respBytes, statusCode, err = client.DeleteRequest(&options)
	}

	if err != nil {
		//If API fails or any other error
		GeneralAPIError(c, err.Error())
		return nil, api_client.DefaultStatusCode
	}

	return respBytes, statusCode
}

func SendRequest(c *gin.Context, serviceType ServiceType, url string, requestType RequestType, headers map[string]interface{}, params map[string]string, body interface{}) {
	respBytes, statusCode := GetRequestResponse(c, serviceType, url, requestType, headers, params, body)
	if respBytes == nil {
		return
	}

	//Parse response
	ParseResponse(c, respBytes, statusCode)

}
