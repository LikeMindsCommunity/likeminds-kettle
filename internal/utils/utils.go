package utils

import (
	"fmt"
	"net/http"

	"github.com/nateshr/likeminds-authentication/internal/logging"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/utils/api_client"
)

/*
When calling any internal LM API but don't want any context handling from Kettle to internal API => handling response bytes, status code, error on your own:
- Call GetRequestResponseWithoutContext to hit API and get response bytes, status code, error
- Then follow any of the following:
  - Unmarshall response bytes by custom logic and send response to client
  - Call ValidateClientResponseWithoutContext to get *api_client.APIClientResponse.Response and send response to client

When calling any internal LM API, want context handling of error => handle response bytes on your own except when !(*api_client.APIClientResponse.Success), don't want to handle error except unmarshall
- Call GetRequestResponse to get response bytes
- Call ValidateClientResponse to get *api_client.APIClientResponse, send response to client when !(*api_client.APIClientResponse.Success), send error to client when unable to unmarshall
- Call GenerateResponse to send response to client

When calling any internal LM API, want context handling of response bytes, error
- Call GetRequestResponse to get response bytes
- Call ParseResponse to send response to client or send error to client
*/

// GetRequestResponseWithoutContext used to hit LM API and return response bytes, status code, error
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

	case PandemoniumService:
		finalUrl = client.PandemoniumServiceBaseUrl + url
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

// GetRequestResponse used to hit LM API using GetRequestResponseWithoutContext. Return: response bytes, status code. Handle error using *gin.Context: error returned by GetRequestResponseWithoutContext
func GetRequestResponse(c *gin.Context, serviceType ServiceType, url string, requestType RequestType, headers map[string]interface{}, params map[string]string, body interface{}) ([]byte, int) {
	responseBytes, statusCode, err := GetRequestResponseWithoutContext(serviceType, url, requestType, headers, params, body)
	if err != nil {
		//If API fails or any other error
		GeneralAPIError(c, err.Error())
		return nil, api_client.DefaultStatusCode
	}

	return responseBytes, statusCode
}

// SendRequest used to hit LM API using GetRequestResponse. Handle response bytes, status code, error internally using *gin.Context: call ParseResponse
func SendRequest(c *gin.Context, serviceType ServiceType, url string, requestType RequestType, headers map[string]interface{}, params map[string]string, body interface{}) {
	responseBytes, statusCode := GetRequestResponse(c, serviceType, url, requestType, headers, params, body)
	if responseBytes == nil {
		return
	}

	//Parse response
	ParseResponse(c, responseBytes, statusCode, false)
}

// CallExternalAPI used to hit external API using GetRequestResponseWithoutContext. Return: status code, response bytes, error
func CallExternalAPI(url string, method RequestType, headers map[string]interface{}, params map[string]string, body interface{},
) ([]byte, int, error) {
	responseBytes, statusCode, err := GetRequestResponseWithoutContext(ExternalService, url, method, headers, params, body)
	if err != nil {
		//If API fails or any other error
		logging.Error(fmt.Sprintf("Error Occured while calling API: %s | status code: %d | error: %s", url, statusCode, err.Error()))
		return nil, api_client.DefaultStatusCode, err
	}

	return responseBytes, statusCode, err
}

// ValidateClientResponse called after GetRequestResponse to unmarshall response bytes. Return: *api_client.APIClientResponse. Handle error internally using *gin.Context: unmarshall error. Handle response internally using *gin.Context: !(*api_client.APIClientResponse.Success)
func ValidateClientResponse(c *gin.Context, responseBytes []byte, statusCode int) *api_client.APIClientResponse {
	//Parse response
	var apiClientResponse api_client.APIClientResponse
	err := api_client.UnmarshalAPIClientResponse(responseBytes, &apiClientResponse)
	if err != nil {
		//Internal unmarshal error
		GeneralAPIError(c, err.Error())
		return nil
	}

	if !apiClientResponse.Success {
		//If internal api returns success as false
		c.JSON(statusCode, apiClientResponse)
		return nil
	}
	return &apiClientResponse
}

// GenerateResponse called after ValidateClientResponse. Handle response internally using *gin.Context: parse widgets and create new Response
func GenerateResponse(c *gin.Context, responseMap map[string]interface{}, parseWidgets bool) {
	//Generating Response Object
	response := Response{
		Success: true,
	}

	// Get widgets data
	if parseWidgets {
		ParseAndFetchWidgets(c, GetUserIdFromContext(c), responseMap)
	}

	//Removing Blank Data Key
	if len(responseMap) > 0 {
		response.Data = responseMap
	}

	c.JSON(http.StatusOK, response)
}

// ParseResponse is called after GetRequestResponse to ValidateClientResponse and GenerateResponse
func ParseResponse(c *gin.Context, responseBytes []byte, statusCode int, parseProfileWidgets bool) {
	apiClientResponse := ValidateClientResponse(c, responseBytes, statusCode)

	if apiClientResponse != nil {
		GenerateResponse(c, apiClientResponse.Response, parseProfileWidgets)
	}
}

// ValidateClientResponseWithoutContext called after GetRequestResponseWithoutContext to unmarshal response bytes. Return: *api_client.APIClientResponse.Response. Log error internally: unmarshall error, !(*api_client.APIClientResponse.Success)
func ValidateClientResponseWithoutContext(responseBytes []byte, statusCode int, err error) map[string]interface{} {
	//If API fails or any other error
	if err != nil {
		logging.Error(fmt.Sprintf("Error Occured : %s", err.Error()))
		return nil
	}

	//Parse response
	var apiClientResponse api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(responseBytes, &apiClientResponse)

	if err != nil {
		//Internal unmarshal error
		logging.Error(fmt.Sprintf("Error while Umarshalling: %s", err.Error()))
		return nil
	}

	if !apiClientResponse.Success {
		//If internal api returns success as false
		logging.Error(fmt.Sprintf("Error Occured :(%d) %s", statusCode, apiClientResponse.ErrorMessage))
		return nil
	}

	return apiClientResponse.Response
}
