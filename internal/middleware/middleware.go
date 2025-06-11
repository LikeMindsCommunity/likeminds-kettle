package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/internal/cache"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/handlers/token"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	log "github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
	"github.com/nateshr/likeminds-authentication/internal/utils/api_client"
)

func OTMValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract OTM from token, internally it checks if token is valid or not
		otm, err := token.ExtractOTM(c.Request.Header.Get(constants.HeaderAuthorization))

		if otm == nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: constants.ErrorInvalidOTM,
			})
			return

		} else {
			// If valid, set "otm" in context, to be used in later APIs
			c.Set(constants.ParamOTM, otm)

			// Set platform type in request header
			if otm.PlatformType != "" {
				c.Request.Header[ContextPlatformTypeheader] = []string{otm.PlatformType}
			}

			// Set API key in request header
			if otm.ApiKey != "" {
				c.Request.Header[ContextApiKeyHeader] = []string{otm.ApiKey}
			}
		}
		c.Next()
	}
}

func VTMValidationMiddleware(isMandatory bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract VTM from token, internally it checks if token is valid or not
		vtm, err := token.ExtractVTM(c.Request.Header.Get(constants.HeaderAuthorization))

		if vtm == nil && isMandatory {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: constants.ErrorInvalidVTM,
			})
			return

		} else if vtm == nil {
			log.Error(err)
			c.Next()

		} else {
			// If valid, set "vtm" in context, to be used in later APIs
			c.Set(constants.ParamVTM, vtm)

			// Set platform type in request header
			if vtm.PlatformType != "" {
				c.Request.Header[ContextPlatformTypeheader] = []string{vtm.PlatformType}
			}

			// // Set API key in request header
			if vtm.ApiKey != "" {
				c.Request.Header[ContextApiKeyHeader] = []string{vtm.ApiKey}
			}
		}
		c.Next()
	}
}

func LTMValidationMiddleware(redisClient *redis.Client, isGuestAccess bool, isWsAPI bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isWsAPI &&
			(c.Request.URL.Query().Get(utils.HeadersPlatformCode) == utils.PlatformReact ||
				c.Request.URL.Query().Get(utils.HeadersPlatformCode) == utils.PlatformReactNative) {
			setContextHeadersFromQueryParams(c)
		}

		bearerToken := c.Request.Header.Get(constants.HeaderAuthorization)
		//Extract LTM from token, internally it checks if token is valid or not
		ltm, err := token.ExtractLTM(bearerToken)
		if ltm == nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: constants.ErrorInvalidLTM,
			})
			return
		} else {
			// Check if LTM is black listed or not
			if cache.IsLTMBlacklisted(redisClient, ltm) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
					Success:      false,
					ErrorMessage: utils.ErrorDeviceLoggedOut,
				})
				return
			}

			// Check if guest access is given
			if ltm.IsGuest && !isGuestAccess {
				c.AbortWithStatusJSON(http.StatusForbidden, utils.Response{
					Success:      false,
					ErrorMessage: utils.ErrorGuestAccessNotAllowed,
				})
				return
			} else if ltm.IsGuest {
				// Add additional headers
				headers := map[string]string{
					utils.HeaderMemberRole: utils.GuestRole,
				}

				utils.AddHeaders(c, headers)
			}

			// If valid and not blacklisted, set "ltm" in context, to be used in later APIs
			c.Set(constants.ParamLTM, ltm)

			// Set API key in request header
			if ltm.ApiKey != "" {
				c.Request.Header[ContextApiKeyHeader] = []string{ltm.ApiKey}
			}
		}
		c.Next()
	}
}

func setContextHeadersFromQueryParams(c *gin.Context) {
	c.Request.Header.Set(constants.HeaderAuthorization, c.Request.URL.Query().Get(constants.HeaderAuthorization))
	c.Request.Header.Set(utils.HeadersMemberId, c.Request.URL.Query().Get(utils.HeadersMemberId))
	c.Request.Header.Set(utils.HeadersVersionCode, c.Request.URL.Query().Get(utils.HeadersVersionCode))
	c.Request.Header.Set(utils.HeadersPlatformCode, c.Request.URL.Query().Get(utils.HeadersPlatformCode))
	c.Request.Header.Set(utils.HeadersPlatformType, c.Request.URL.Query().Get(utils.HeadersPlatformType))
	c.Request.Header.Set(utils.HeadersSdkSource, c.Request.URL.Query().Get(utils.HeadersSdkSource))
	c.Request.Header.Set(utils.HeadersDeviceId, c.Request.URL.Query().Get(utils.HeadersDeviceId))
	c.Request.Header.Set(utils.HeadersApiKey, c.Request.URL.Query().Get(utils.HeadersApiKey))
	c.Request.Header.Set(utils.HeadersAcceptVersion, c.Request.URL.Query().Get(utils.HeadersAcceptVersion))
	c.Request.Header.Set(utils.HeadersApiVersion, c.Request.URL.Query().Get(utils.HeadersApiVersion))
	c.Request.Header.Set(utils.HeaderMemberRole, c.Request.URL.Query().Get(utils.HeaderMemberRole))
	c.Request.Header.Set(utils.HeaderContentType, c.Request.URL.Query().Get(utils.HeaderContentType))
}

func RTMValidationMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Extract RTM from token, internally it checks if token is valid or not
		rtm, err := token.ExtractRTM(c.Request.Header.Get(constants.HeaderAuthorization))
		if rtm == nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: constants.ErrorInvalidRTM,
			})
			return
		} else {
			//Check if RTM is black listed or not
			if cache.IsRTMBlacklisted(redisClient, rtm) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
					Success:      false,
					ErrorMessage: utils.ErrorDeviceLoggedOut,
				})
				return
			}
			//If valid and not blacklisted, set "rtm" in context, to be used in later APIs
			c.Set(constants.ParamRTM, rtm)
		}
		// Set API key in request header
		if rtm.ApiKey != "" {
			c.Request.Header[ContextApiKeyHeader] = []string{rtm.ApiKey}
		}
		c.Next()
	}
}

func LTMorVTMValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from headers
		bearerToken := c.Request.Header.Get(constants.HeaderAuthorization)

		if bearerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: constants.ErrorInvalidLTMorVTM,
			})
			return
		}

		// Extract LTM info from token, internally it checks if token is valid or not
		ltm, ltmErr := token.ExtractLTM(bearerToken)

		if ltmErr == nil {
			c.Set(constants.ParamLTM, ltm)

			// Set API key in request header
			if ltm.ApiKey != "" {
				c.Request.Header[ContextApiKeyHeader] = []string{ltm.ApiKey}
			}

			c.Next()
			return
		}

		// Extract VTM info from token, internally it checks if token is valid or not
		vtm, vtmErr := token.ExtractVTM(bearerToken)

		if vtmErr == nil {
			c.Set(constants.ParamVTM, vtm)

			// Set API key in request header
			if vtm.ApiKey != "" {
				c.Request.Header[ContextApiKeyHeader] = []string{vtm.ApiKey}
			}

			c.Next()

		} else {
			log.Error(ltmErr)
			log.Error(vtmErr)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: constants.ErrorInvalidLTMorVTM,
			})
			return
		}
	}
}

func LogoutValidationMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		//Extract LTM from token, internally it checks if token is valid or not
		ltm, err := token.ExtractLTM(c.Request.Header.Get(constants.HeaderAuthorization))
		if ltm == nil {

			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: constants.ErrorInvalidLTM,
			})
			return
		} else {

			//Check if LTM is black listed or not
			if cache.IsLTMBlacklisted(redisClient, ltm) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
					Success:      false,
					ErrorMessage: utils.ErrorDeviceLoggedOut,
				})
				return
			}

			//Get RTM token from body
			var logoutRequest user.LogoutRequest
			if err := c.ShouldBindJSON(&logoutRequest); err != nil {
				c.JSON(http.StatusUnprocessableEntity, utils.Response{
					Success:      false,
					ErrorMessage: utils.ErrorInvalidRequest,
				})
				return
			}

			//Extract and validate RTM token
			rtm, err := extractAndValidateRTMToken(logoutRequest.RefreshToken, redisClient)
			if rtm == nil {
				log.Error(err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
					Success:      false,
					ErrorMessage: constants.ErrorInvalidRTM,
				})
				return
			} else {
				//If valid and not blacklisted, set "ltm" and "rtm" in context, to be used in later APIs
				c.Set(constants.ParamLTM, ltm)
				c.Set(constants.ParamRTM, rtm)
			}
		}
		c.Next()
	}
}

func extractAndValidateRTMToken(refreshToken string, redisClient *redis.Client) (*constants.RefreshTokenMeta, error) {

	//Extract RTM from token, internally it checks if token is valid or not
	rtm, err := token.ExtractRTM(refreshToken)
	if err != nil {
		return nil, err
	}

	//Check if RTM is black listed or not
	if cache.IsRTMBlacklisted(redisClient, rtm) {
		return nil, fmt.Errorf(utils.ErrorDeviceLoggedOut)
	}

	return rtm, nil
}

// Internal service validation middelware
func InternalServiceValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request is from internal service
		platformType := c.Request.Header.Get(utils.HeadersPlatformType)
		if !utils.CheckIfStringExistsInArray([]string{string(utils.PlatformSwarmService), string(utils.PlatformCaravanService), string(utils.PlatformSkulkService)}, platformType) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: utils.ErrorMemeberAccessFail,
			})
			return
		}
		c.Next()
	}
}

// ApiMiddleware will add the db connection to the context
func ApiMiddleware(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(cache.ParamRedisClient, client)
		c.Next()
	}
}

// AddResponseHeadersMiddleware | adds necessary API response headers
func AddResponseHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		addTransportSecurityHeaders(c)
		addAPIContentResponseHeaders(c)
		addCacheControlHeaders(c)
		c.Next()
	}
}

func addTransportSecurityHeaders(c *gin.Context) {
	c.Header(StrictTransportSecurityHeaderKey, StrictTransportSecurityHeaderValue)
}

func addAPIContentResponseHeaders(c *gin.Context) {
	c.Header(ContentTypeOptionsHeaderKey, ContentTypeOptionsHeaderValue)
}

func addCacheControlHeaders(c *gin.Context) {
	c.Header(CacheControlHeaderKey, CacheControlHeaderValue)
}

// GuestAccessCheckMiddleware | restrict guest access on endpoints
func GuestAccessCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Create internal API client
		client := api_client.NewAPIClient()
		options := api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + user.UserFetchEndPoint,
			CustomHeaders: utils.CreateHeaders(c, c.GetHeader(utils.HeadersMemberId)),
		}
		//Send request
		respBytes, statusCode, err := client.GetRequest(&options)
		if err != nil {
			//If API fails or any other error
			utils.GeneralAPIError(c, err.Error())
			return
		}
		//Parse response
		var apiCR api_client.APIClientResponse
		err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
		if err != nil {
			//Internal unmarshal error
			utils.GeneralAPIError(c, err.Error())
			return
		}

		if !apiCR.Success {
			//If api/user/fetch returns success as false
			c.JSON(statusCode, apiCR)
			return
		}

		isGuest := apiCR.Response[user.ResponseUser].(map[string]interface{})[user.ResponseUserIsGuest].(bool)
		if isGuest {
			type GuestAccessDeniedResponseData struct {
				Route string `json:"route"`
			}
			response := utils.Response{
				Success:      false,
				ErrorMessage: utils.ErrorGuestAccessDenied,
				Data:         GuestAccessDeniedResponseData{Route: user.GuestLoginRoute},
			}

			// If user is guest returns success as false
			utils.APIError(c, http.StatusUnauthorized, response)
			return
		}
		c.Next()
	}
}

// CustomRecoveryMiddleware is a Gin middleware to handle panics
func CustomRecoveryMiddleware(c *gin.Context, err interface{}) {

	// Send Internal server error with error_message
	utils.GeneralAPIError(c, utils.ErrorSomethingWentWrong)
}

func EnableCorsMiddleware() gin.HandlerFunc {

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders(
		"x-member-id",
		"x-api-version",
		"x-platform-code",
		"x-platform-type",
		"x-version-code",
		"x-sdk-source",
		"x-accept-version",
		"x-username",
		"x-password",
		"x-device-id",
		"x-api-key",
		"Authorization",
	)

	return cors.New(config)
}

// Method to process API request to log
func processRequest(c *gin.Context) interface{} {
	requestBodyData := gin.H{}

	// Reading request body
	requestBody, err := ioutil.ReadAll(c.Request.Body)

	// Updating request body after read
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Unmarshalling request body
	if err == nil {
		_ = json.Unmarshal(requestBody, &requestBodyData)
	}

	return gin.H{
		"host":         c.Request.Host,
		"absolute_uri": c.Request.RequestURI,
		"method":       c.Request.Method,
		"headers":      c.Request.Header,
		"body":         requestBodyData,
	}
}

// responseBodyWriter | Custom Response Writer
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write | Custom Write Method for responseBodyWriter
func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// sanitizeRequestHeaders removes sensitive headers and returns a sanitized copy
func sanitizeRequestHeaders(requestData gin.H) {
	if headers, ok := requestData["headers"].(http.Header); ok {
		headersCopy := make(map[string]string)
		for key, values := range headers {
			if key != constants.HeaderAuthorization {
				if len(values) > 0 {
					// Remove hyphens from header keys
					sanitizedKey := strings.ReplaceAll(key, "-", "_")
					headersCopy[sanitizedKey] = values[0]
				}
			}
		}
		requestData["headers"] = headersCopy
	}
}

// addUserContextHeaders adds user-specific headers from LTM token
func addUserContextHeadersFromToken(data gin.H, authHeader string) {
	ltm, _ := token.ExtractLTM(authHeader)
	if ltm != nil {
		headers := data["request"].(gin.H)["headers"].(map[string]string)
		if ltm.UserUniqueID != "" {
			headers["X_Member_Id"] = ltm.UserUniqueID
		}
		if ltm.ApiKey != "" {
			headers["X_Api_Key"] = ltm.ApiKey
		}
	}
}

// addUserContextHeaders adds user-specific headers from LTM in context
func addUserContextHeadersFromContext(data gin.H, ltm *constants.LoginTokenMeta) {
	headers := data["request"].(gin.H)["headers"].(map[string]string)
	if ltm.UserUniqueID != "" {
		headers["X_Member_Id"] = ltm.UserUniqueID
	}
	if ltm.ApiKey != "" {
		headers["X_Api_Key"] = ltm.ApiKey
	}
}

// LoggingMiddleware will log the request and response of API
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.RequestURI == "/" {

			c.Next()

		} else {

			data := gin.H{}

			// Starting time
			startTime := time.Now()

			// Implementing custom response body writer in the context
			w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
			c.Writer = w

			// Updating Request Data
			data["request"] = processRequest(c)

			sanitizeRequestHeaders(data["request"].(gin.H))

			// Processing request
			c.Next()

			authToken := c.Request.Header.Get(constants.HeaderAuthorization)
			if authToken == "" {
				authToken, _ := c.Get(constants.ParamLTM)
				if authToken != nil {
					addUserContextHeadersFromContext(data, authToken.(*constants.LoginTokenMeta))
				}
			} else {
				addUserContextHeadersFromToken(data, authToken)
			}

			// End Time
			endTime := time.Now()

			response := gin.H{}
			statusCode := c.Writer.Status()

			// Unmarshalling Request Response
			_ = json.Unmarshal(w.body.Bytes(), &response)

			// Updating Request Response
			data["response"] = gin.H{
				"http_response_code": statusCode,
				"content":            response,
			}

			// We are logging the response "content" for successful requests as well. Earlier, we were only logged for error responses.

			// Updating Request Meta Data
			data["meta"] = gin.H{
				"latency":   endTime.Sub(startTime),
				"client_ip": c.ClientIP(),
			}

			if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
				// Logging the generated request data as Info
				log.InfoWithFields(data)
			} else {
				// Logging the generated request data as Error
				log.ErrorWithFields(data)
			}

			c.Next()
		}
	}
}
