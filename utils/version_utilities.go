package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func CheckVersion(c *gin.Context, featureVersionCode map[string]int) bool {
	/*
		returns True if,
		  versionCode >= featureVersionCode for the given platform
		returns False for all other cases
	*/

	var isVersionCheck bool = false

	headers := CreateHeaders(c, "")

	platformCode, platformCodeCheck := headers[HeadersPlatformCode].(string)
	versionCode, versionCodeCheck := headers[HeadersVersionCode]

	if platformCodeCheck && versionCodeCheck {
		featureVersionCodeForPlatform, ok := featureVersionCode[platformCode]
		versionCode, isVersionCodeConverted := versionCode.(string)

		if ok && isVersionCodeConverted {
			versionCode, versionCodeConversionErr := strconv.Atoi(versionCode)

			if versionCodeConversionErr == nil {
				isVersionCheck = versionCode >= featureVersionCodeForPlatform
			}
		}
	}

	return isVersionCheck

}

func ApiRevampV1Check(c *gin.Context) bool {
	/*
		Api Revamp V1 Check

		returns True if, x-accept-version == v1
	*/

	acceptVersion := c.GetHeader(HeadersAcceptVersion)

	return acceptVersion == "v1"

}
