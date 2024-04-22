package utils

import (
	"fmt"
	"net/http"

	"github.com/nateshr/likeminds-authentication/logging"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
)

// Generate Response to be sent on request success
func GenerateResponse(c *gin.Context, dataResponse map[string]interface{}, parseWidgets bool) {
	//Generating Response Object
	response := Response{
		Success: true,
	}

	// Get widgets data
	if parseWidgets {
		ParseAndFetchWidgets(c, GetUserIdFromContext(c), dataResponse)
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

// Validate & Parse Response for request sent internally
func ValidateClientResponseWithoutContext(respBytes []byte, statuscode int, err error) map[string]interface{} {

	//If API fails or any other error
	if err != nil {
		logging.Error(fmt.Sprintf("Error Occured : %s", err.Error()))
		return nil
	}

	//Parse response
	var apiCR api_client.APIClientResponse
	marshal_err := api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)

	if marshal_err != nil {
		//Internal unmarshal error
		logging.Error(fmt.Sprintf("Error while Umarshalling: %s", marshal_err.Error()))
		return nil
	}

	if !apiCR.Success {
		//If internal api returns success as false
		logging.Error(fmt.Sprintf("Error Occured :(%d) %s", statuscode, apiCR.ErrorMessage))
		return nil
	}

	return apiCR.Response
}

// ParseResponse from request sent internally
func ParseResponse(c *gin.Context, respBytes []byte, statusCode int, parseProfileWidgets bool) {

	apiCR := ValidateClientResponse(c, respBytes, statusCode)

	if apiCR != nil {
		GenerateResponse(c, apiCR.Response, parseProfileWidgets)
	}
}

// Method to Send Request
func GetRequestResponseWithoutContext(serviceType ServiceType, url string, requestType RequestType, headers map[string]interface{}, params map[string]string, body interface{}) ([]byte, int, error) {
	//Create internal API client
	client := api_client.NewAPIClient()
	var finalUrl string
	var respBytes []byte
	var statusCode int
	var err error

	switch serviceType {
	case CoreService:
		finalUrl = client.CoreServiceBaseURL + url

	case SubscriptionService:
		finalUrl = client.SubscriptionServiceBaseURL + url

	case SwarmService:
		finalUrl = client.SwarmServiceBaseUrl + url

	case ExternalService:
		finalUrl = url
	}

	switch requestType {
	case GETRequest:

		options := api_client.GetRequestOptions{
			Url:           finalUrl,
			CustomHeaders: headers,
			Params:        params,
		}

		respBytes, statusCode, err = client.GetRequest(&options)

	case POSTRequestRawBody:

		options := api_client.PostRequestOptions{
			Url:           finalUrl,
			CustomHeaders: headers,
			Params:        params,
			Body:          body,
		}

		respBytes, statusCode, err = client.PostRequest(&options, api_client.BodyTypeRaw)

	case POSTRequestFormUrlEncodedBody:

		options := api_client.PostRequestOptions{
			Url:           finalUrl,
			CustomHeaders: headers,
			Params:        params,
			Body:          body,
		}

		respBytes, statusCode, err = client.PostRequest(&options, api_client.BodyTypeFormUrlEncoded)

	case PUTRequest:

		options := api_client.PostRequestOptions{
			Url:           finalUrl,
			CustomHeaders: headers,
			Params:        params,
			Body:          body,
		}

		respBytes, statusCode, err = client.PutRequest(&options)

	case DELETERequest:

		options := api_client.PostRequestOptions{
			Url:           finalUrl,
			CustomHeaders: headers,
			Params:        params,
			Body:          body,
		}

		respBytes, statusCode, err = client.DeleteRequest(&options)

	case PATCHRequest:

		options := api_client.PostRequestOptions{
			Url:           finalUrl,
			CustomHeaders: headers,
			Params:        params,
			Body:          body,
		}

		respBytes, statusCode, err = client.PatchRequest(&options)
	}

	return respBytes, statusCode, err
}

func GetRequestResponse(c *gin.Context, serviceType ServiceType, url string, requestType RequestType, headers map[string]interface{}, params map[string]string, body interface{}) ([]byte, int) {

	respBytes, statusCode, err := GetRequestResponseWithoutContext(serviceType, url, requestType, headers, params, body)
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
	ParseResponse(c, respBytes, statusCode, false)

}

// Utility function to call external API
func CallExternalAPI(url string, method RequestType, headers map[string]interface{}, params map[string]string, body interface{},
) ([]byte, int, error) {

	respBytes, statusCode, err := GetRequestResponseWithoutContext(ExternalService, url, method, headers, params, body)
	if err != nil {
		//If API fails or any other error
		logging.Error(fmt.Sprintf("Error Occured while calling API: %s | status code: %d | error: %s", url, statusCode, err.Error()))
		return nil, api_client.DefaultStatusCode, err
	}

	return respBytes, statusCode, err
}
