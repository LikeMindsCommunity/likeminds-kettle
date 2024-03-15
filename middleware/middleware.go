package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/cache"
	log "github.com/nateshr/likeminds-authentication/logging"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func OTMValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract OTM from token, internally it checks if token is valid or not
		otm, err := token.ExtractOTM(c.Request.Header.Get(token.HeaderAuthorization))

		if otm == nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidOTM,
			})
			return

		} else {
			// If valid, set "otm" in context, to be used in later APIs
			c.Set(token.ParamOTM, otm)

			// Set API key in request header
			if otm.ApiKey != "" {
				c.Request.Header["X-Api-Key"] = []string{otm.ApiKey}
			}
		}
		c.Next()
	}
}

func VTMValidationMiddleware(isMandatory bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract VTM from token, internally it checks if token is valid or not
		vtm, err := token.ExtractVTM(c.Request.Header.Get(token.HeaderAuthorization))

		if vtm == nil && isMandatory {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidVTM,
			})
			return

		} else if vtm == nil {
			log.Error(err)
			c.Next()

		} else {
			// If valid, set "vtm" in context, to be used in later APIs
			c.Set(token.ParamVTM, vtm)

			// // Set API key in request header
			if vtm.ApiKey != "" {
				c.Request.Header["X-Api-Key"] = []string{vtm.ApiKey}
			}
		}
		c.Next()
	}
}

func LTMValidationMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get(token.HeaderAuthorization)
		//Extract LTM from token, internally it checks if token is valid or not
		ltm, err := token.ExtractLTM(bearerToken)
		if ltm == nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidLTM,
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
			//If valid and not blacklisted, set "ltm" in context, to be used in later APIs
			c.Set(token.ParamLTM, ltm)

			// Set API key in request header
			if ltm.ApiKey != "" {
				c.Request.Header["X-Api-Key"] = []string{ltm.ApiKey}
			}
		}
		c.Next()
	}
}

func RTMValidationMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Extract RTM from token, internally it checks if token is valid or not
		rtm, err := token.ExtractRTM(c.Request.Header.Get(token.HeaderAuthorization))
		if rtm == nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidRTM,
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
			c.Set(token.ParamRTM, rtm)
		}
		// Set API key in request header
		if rtm.ApiKey != "" {
			c.Request.Header["X-Api-Key"] = []string{rtm.ApiKey}
		}
		c.Next()
	}
}

func LTMorVTMValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from headers
		bearerToken := c.Request.Header.Get(token.HeaderAuthorization)

		if bearerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidLTMorVTM,
			})
			return
		}

		// Extract LTM info from token, internally it checks if token is valid or not
		ltm, ltmErr := token.ExtractLTM(bearerToken)

		if ltmErr == nil {
			c.Set(token.ParamLTM, ltm)

			// Set API key in request header
			if ltm.ApiKey != "" {
				c.Request.Header["X-Api-Key"] = []string{ltm.ApiKey}
			}

			c.Next()
		}

		// Extract VTM info from token, internally it checks if token is valid or not
		vtm, vtmErr := token.ExtractVTM(bearerToken)

		if vtmErr == nil {
			c.Set(token.ParamVTM, vtm)

			// Set API key in request header
			if vtm.ApiKey != "" {
				c.Request.Header["X-Api-Key"] = []string{vtm.ApiKey}
			}

			c.Next()

		} else {
			log.Error(ltmErr)
			log.Error(vtmErr)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidLTMorVTM,
			})
			return
		}
	}
}

func LogoutValidationMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Extract LTM from token, internally it checks if token is valid or not
		ltm, err := token.ExtractLTM(c.Request.Header.Get(token.HeaderAuthorization))
		if ltm == nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidLTM,
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
			//Extract RTM from token, internally it checks if token is valid or not
			rtm, err := token.ExtractRTM(logoutRequest.RefreshToken)
			if rtm == nil {
				log.Error(err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
					Success:      false,
					ErrorMessage: token.ErrorInvalidRTM,
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
				//If valid and not blacklisted, set "ltm" and "rtm" in context, to be used in later APIs
				c.Set(token.ParamLTM, ltm)
				c.Set(token.ParamRTM, rtm)
			}
		}
		c.Next()
	}
}

// Internal service validation middelware
func InternalServiceValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request is from internal service
		platformType := c.Request.Header.Get(utils.HeadersPlatformType)
		if !utils.CheckIfStringExistsInArray([]string{string(utils.PlatformSwarmService), string(utils.PlatformCaravanService)}, platformType) {
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

			// Processing request
			c.Next()

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

			if statusCode < http.StatusBadRequest {
				data["response"].(gin.H)["content"] = gin.H{}
			}

			// Updating Request Meta Data
			data["meta"] = gin.H{
				"latency":   endTime.Sub(startTime),
				"client_ip": c.ClientIP(),
			}

			// Marshalling the final Data
			marshelledData, _ := json.Marshal(data)

			if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
				// Logging the generated request data as Info
				log.Info(string(marshelledData))
			} else {
				// Logging the generated request data as Error
				log.Error(string(marshelledData))
			}

			c.Next()
		}
	}
}

// Ratelimiting Middleware will limit API requests for a given community
func RatelimitingMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	// TODO - check for all the return statement and add/remove according to the flows
	return func(c *gin.Context) {
		headers := utils.CreateHeaders(c, "")
		apiKey := headers[utils.HeadersApiKey]

		// Community Billing Value from Cache
		communityBillingDataMapFromCache := map[string]interface{}{}
		// // communityBillingDataNew, err := json.Marshal(communityBillingDataMapFromCache)
		// if err != nil {
		// 	log.Error(err)
		// 	return
		// }

		// Get community billing data from cache
		communityBillingData, keyExists, err := cache.Get(redisClient, cache.CommunityBillingDataKey)
		if err != nil {
			log.Error(err)
			return
		}

		// If key doesn't exist in cache
		// if !keyExists {
		// 	// Set a blank map[string]interface{} in cache

		// 	err := cache.Set(redisClient, cache.CommunityBillingDataKey, communityBillingDataNew, cache.CommunityBillingDataTTL*time.Hour)
		// 	if err != nil {
		// 		log.Error(err)
		// 		return
		// 	}
		// } else {
		if keyExists {
			// Unmarshal the value from the cache
			err = json.Unmarshal([]byte(communityBillingData), &communityBillingDataMapFromCache)
			if err != nil {
				log.Error(err)
				return
			}
		}

		// Get Community Billing Data for a specific api key
		apiKeyBillingValue, ok := communityBillingDataMapFromCache[apiKey.(string)]

		// If value doesn't exist for given api key
		if !ok {
			// Fetch CommunityID for the API key
			communityIdFromCache, err := utils.FetchCommunityIdFromApiKey(redisClient, apiKey.(string))
			if err != nil {
				log.Error(err)
				return
			}

			// Fetch from Subscription service
			respBytes, statusCode := utils.GetRequestResponse(c, utils.SubscriptionService, fmt.Sprintf("%s/%d", utils.BillingPlanEnpoint, communityIdFromCache), utils.GETRequest, headers, nil, nil)

			// Validate response from Skulk service api call
			apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
			if apiCR == nil {
				return
			}

			// Parsing the data from api call
			skulkResponse := apiCR.Response
			//Extract Billing Data from Skulk Response
			if billingData, ok := skulkResponse[cache.BillingDataKey]; ok {
				// Set billing data for a specific api key in cache
				communityBillingDataMapFromCache[apiKey.(string)] = billingData
				cache.Set(redisClient, cache.CommunityBillingDataKey, communityBillingDataMapFromCache, 86400)
				apiKeyBillingValue = billingData
			}

		}

		// Setting tier_type for the corresponding API key
		tierType, ok := apiKeyBillingValue.(map[string]interface{})[cache.TierTypeKey]

		if !ok {
			// TODO - update error log
			log.Error("")
			return
		}

		// Tier Data Value from Cache
		tierDataMapFromCache := map[string]interface{}{}

		// Get tier_data from cache
		tierData, tierDataExists, err := cache.Get(redisClient, cache.TierDataKey)
		if err != nil {
			log.Error(err)
			return
		}

		// If data doesn't exist
		// if !tierDataExists {
		// 	// Set default value in cache
		// 	err = cache.Set(redisClient, cache.TierDataKey, map[string]interface{}{}, 86400*30)
		// 	if err != nil {
		// 		log.Error(err)
		// 		return
		// 	}

		// } else
		if tierDataExists {
			// Unmarshal value from cache
			err = json.Unmarshal([]byte(tierData), &tierDataMapFromCache)
			if err != nil {
				log.Error(err)
				return
			}
		}

		// Get Tier Data for a specific tier type
		tierValue, ok := tierDataMapFromCache[tierType.(string)]

		// If data doesn't exist for given tier type
		if !ok {
			// Get data from skulk service

			// prepare params for api call
			params := map[string]string{
				utils.ParamTierType: c.Query(utils.ParamTierType),
			}

			// Get data from api call
			respBytes, statusCode := utils.GetRequestResponse(c, utils.SubscriptionService, utils.TierEndpoint, utils.GETRequest, headers, params, nil)

			// Validate response from api call
			apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
			if apiCR == nil {
				return
			}

			// parsing the data from api call
			skulkResponse := apiCR.Response
			if tierData, ok := skulkResponse[cache.TierDataKey]; ok {
				// Set tier data for a specific tier type in cache
				tierDataMapFromCache[tierType.(string)] = tierData
				cache.Set(redisClient, cache.TierDataKey, tierDataMapFromCache, cache.TierDataTTL*time.Hour)
				tierValue = tierData
			}
		}

		// Iterate over the factors in tier value
		for _, limitFactor := range tierValue.([]map[string]interface{}) {
			rateLimitCurrentValueKey := limitFactor[RateLimitKeyNameKey].(string)
			rateLimitTTLValue := limitFactor[RateLimitTTLKey].(int)
			rateLimitMaxRequestValue := limitFactor[RateLimitMaxRequestLimitValueKey].(int)
			rateLimitErrorMessage := limitFactor[RateLimitErrorMessageKey].(string)
			rateLimitCurrentValueFromCache := 0

			// Get rate limit current value from cache
			rateLimitCurrentValue, keyExists, err := cache.Get(redisClient, rateLimitCurrentValueKey)
			if err != nil {
				log.Error(err)
				return
			}
			//convert RateLimit String value to Int
			rateLimitValue, err := strconv.Atoi(rateLimitCurrentValue)
			if err != nil {
				log.Error(err)
				return
			}

			// If key doesn't exists
			// if !keyExists {
			// 	// Set a new key in cache with rateLimitCurrentValueFromCache and ttl=rateLimitTTLValue
			// 	err = cache.Set(redisClient, rateLimitCurrentValueKey, rateLimitCurrentValueFromCache, time.Duration(rateLimitTTLValue))
			// 	if err != nil {
			// 		log.Error(err)
			// 		return
			// 	}
			// } else
			if keyExists {
				// Unmarshal the value from cache to rateLimitCurrentValue
				err = json.Unmarshal([]byte(rateLimitCurrentValue), &rateLimitCurrentValueFromCache)
				if err != nil {
					log.Error(err)
				}
			}

			if rateLimitCurrentValueFromCache <= rateLimitMaxRequestValue {
				//increment rateLimit Value and Set in cache
				rateLimitValue++
				err = cache.Set(redisClient, rateLimitCurrentValueKey, rateLimitValue, time.Duration(rateLimitTTLValue))
				if err != nil {
					log.Error(err)
					return
				}
			} else {
				response := utils.Response{
					Success:      false,
					ErrorMessage: rateLimitErrorMessage,
				}
				//Send Rate Limit Error
				utils.APIError(c, http.StatusTooManyRequests, response)
			}
		}
		c.Next()
	}
}
