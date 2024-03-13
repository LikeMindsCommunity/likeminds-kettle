package utils

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// RequestParamsToMap Converts a struct to a map while maintaining the json alias as keys
func RequestParamsToMap(obj interface{}) (newMap map[string]string, err error) {
	data, err := json.Marshal(obj) // Convert to a json string
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}

// CreateHeaders Used to create headers for our internal APIs
func CreateHeaders(c *gin.Context, userUniqueID string) map[string]interface{} {
	headers := make(map[string]interface{})
	if len(userUniqueID) > 0 {
		headers[HeadersMemberId] = userUniqueID
	}
	headers[HeadersPlatformCode] = c.GetHeader(HeadersPlatformCode)
	headers[HeadersVersionCode] = c.GetHeader(HeadersVersionCode)
	headers[HeadersSdkSource] = c.GetHeader(HeadersSdkSource)
	headers[HeadersDeviceId] = c.GetHeader(HeadersDeviceId)
	headers[HeadersApiKey] = c.GetHeader(HeadersApiKey)
	headers[HeadersAcceptVersion] = c.GetHeader(HeadersAcceptVersion)
	headers[HeadersApiVersion] = c.GetHeader(HeadersApiVersion)
	headers[HeaderMemberRole] = c.GetHeader(HeaderMemberRole)
	return headers
}

// Add more headers to already existing headers
func AddHeaders(c *gin.Context, headerValMap map[string]string) {

	if headerValMap != nil {

		for header, val := range headerValMap {
			c.Request.Header.Add(header, val)
		}
	}
}
