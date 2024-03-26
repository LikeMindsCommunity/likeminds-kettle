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
	"github.com/nateshr/likeminds-authentication/constants"
	"github.com/nateshr/likeminds-authentication/environment"
	log "github.com/nateshr/likeminds-authentication/logging"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
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

			// // Set API key in request header
			if vtm.ApiKey != "" {
				c.Request.Header["X-Api-Key"] = []string{vtm.ApiKey}
			}
		}
		c.Next()
	}
}

func LTMValidationMiddleware(redisClient *redis.Client, isGuestAccess bool) gin.HandlerFunc {
	return func(c *gin.Context) {
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
				c.Request.Header["X-Api-Key"] = []string{ltm.ApiKey}
			}
		}
		c.Next()
	}
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
			c.Request.Header["X-Api-Key"] = []string{rtm.ApiKey}
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
				c.Request.Header["X-Api-Key"] = []string{ltm.ApiKey}
			}

			c.Next()
		}

		// Extract VTM info from token, internally it checks if token is valid or not
		vtm, vtmErr := token.ExtractVTM(bearerToken)

		if vtmErr == nil {
			c.Set(constants.ParamVTM, vtm)

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
			//Extract RTM from token, internally it checks if token is valid or not
			rtm, err := token.ExtractRTM(logoutRequest.RefreshToken)
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
				//If valid and not blacklisted, set "ltm" and "rtm" in context, to be used in later APIs
				c.Set(constants.ParamLTM, ltm)
				c.Set(constants.ParamRTM, rtm)
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

type CommunityBillingMeta struct {
	TierType int `json:"tier_type"`
}

//	type tierDataMAURPM struct {
//		MAU tierDataType `json:"mau"`
//		RPM tierDataType `json:"rpm"`
//	}
type TierDataType struct {
	MaxRequestLimitValue int    `json:"max_request_limit_value"`
	TTL                  int    `json:"ttl"`
	RateLimitKeyName     string `json:"rate_limit_key_name"`
	ErrorMessage         string `json:"error_message"`
}

type ApitierDataType struct {
	*TierDataType
	TierType      int `json:"tier_type"`
	TierValueType int `json:"tier_value_type"`
}

type TierTypeApiResponse struct {
	Success bool              `json:"success"`
	Data    []ApitierDataType `json:"data"`
}

func RateLimitingMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	if environment.GoDotEnvVariable("IS_RATE_LIMITING_ENABLED") == "true" {
		return func(c *gin.Context) {
			startTime := time.Now()

			headers := utils.CreateHeaders(c, "")
			apiKey := headers[utils.HeadersApiKey].(string)

			// Get community Billing data  function
			// Get CommunityId from apiKey
			communityId, err := utils.FetchCommunityIdFromApiKey(redisClient, apiKey)
			if err != nil {
				log.Error(err)
				return
			}
			cacheKey := fmt.Sprintf(cache.CommunityBillingDataKey, communityId)
			// Get Community Billing data from cache
			communityBillingData, err := FetchCommunityBillingData(redisClient, communityId, cacheKey, headers)
			if err != nil {
				log.Error(err)
				return
			}

			//Get TierData Function
			tierType := communityBillingData.TierType
			tierData, err := FetchTierData(redisClient, communityId, headers, tierType)
			if err != nil {
				log.Error(err)
				return
			}

			for _, limitFactor := range tierData {
				rateLimitCurrentValueKey := limitFactor.RateLimitKeyName + fmt.Sprintf("_%d", communityId)
				rateLimitValue := limitFactor.MaxRequestLimitValue
				rateLimitErrorMessage := limitFactor.ErrorMessage
				rateLimitTTL := limitFactor.TTL
				rateLimitTierValueType := limitFactor.TierValueType

				// Get rate limit current value from cache
				currentValue, exists, err := cache.Get(redisClient, rateLimitCurrentValueKey)
				if err != nil {
					log.Error(err)
					return
				}
				if !exists {
					redisClient.Incr(rateLimitCurrentValueKey)
					redisClient.ExpireAt(rateLimitCurrentValueKey, time.Now().Add(time.Second*time.Duration(rateLimitTTL)))
					currentValue = "1"
				}

				isAllowed, errMessage, err := calculateRateLimitLogic(rateLimitTierValueType, currentValue, rateLimitValue, rateLimitErrorMessage, rateLimitCurrentValueKey, redisClient)
				if err != nil {
					log.Error(err)
					return
				}
				if !isAllowed {
					response := utils.Response{
						Success:      false,
						ErrorMessage: errMessage,
					}
					c.AbortWithStatusJSON(http.StatusTooManyRequests, response)
					return
				}
			}

			fmt.Println("Time taken to execute RateLimitingMiddleware: ", time.Since(startTime))
			c.Next()
		}
	}
	return func(c *gin.Context) {
		c.Next()
	}
}

func calculateRateLimitLogic(rateLimitTierValueType int, currentValue string, rateLimitValue int, rateLimitErrorMessage string, rateLimitCurrentValueKey string, redisClient *redis.Client) (bool, string, error) {
	currentValueInt, err := strconv.Atoi(currentValue)
	if err != nil {
		return false, "", err
	}

	switch rateLimitTierValueType {
	//RPM
	case 1:
		// If rate limit current value is less than rate limit value
		if currentValueInt > rateLimitValue {
			return false, rateLimitErrorMessage, nil
		}
		redisClient.Incr(rateLimitCurrentValueKey)
	}
	return true, "", nil
}

func FetchCommunityBillingData(redisClient *redis.Client, communityId int, cacheKey string, headers map[string]interface{}) (CommunityBillingMeta, error) {
	// Get Community Billing data from cache
	value, valueExists, err := cache.Get(redisClient, cacheKey)
	//If error continue to next middleware
	if err != nil {
		return CommunityBillingMeta{}, err
	}

	communityBillingMeta := CommunityBillingMeta{}

	// communityBillingDataApi
	if !valueExists {
		// Get Value from API
		respBytes, _, err := utils.GetRequestResponseWithoutContext(utils.SubscriptionService, fmt.Sprintf("/%s/%d", utils.BillingPlanEnpoint, communityId), utils.GETRequest, headers, nil, nil)
		if err != nil {
			return communityBillingMeta, err
		}

		err = json.Unmarshal(respBytes, &communityBillingMeta)
		if err != nil {
			return communityBillingMeta, err
		}

		communityBillingMetaForCache, err := json.Marshal(communityBillingMeta)

		if err != nil {
			return communityBillingMeta, err
		}
		// Update value in Cache
		err = cache.Set(redisClient, cacheKey, communityBillingMetaForCache, time.Hour*cache.CommunityBillingDataTTL)
		if err != nil {
			return communityBillingMeta, err
		}

	} else {
		// Unmarshal value from cache
		err := json.Unmarshal([]byte(value), &communityBillingMeta)
		if err != nil {
			return communityBillingMeta, err
		}
	}
	return communityBillingMeta, nil
}

func FetchTierData(redisClient *redis.Client, communityId int, headers map[string]interface{}, tierType int) ([]ApitierDataType, error) {

	cacheKey := fmt.Sprintf("tier_data_%d", tierType)
	value, exists, err := cache.Get(redisClient, cacheKey)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	tierData := []ApitierDataType{}

	if !exists {
		params := map[string]string{
			utils.ParamTierType: strconv.Itoa(tierType),
		}
		// Get data from skulk service
		respBytes, _, err := utils.GetRequestResponseWithoutContext(utils.SubscriptionService, utils.TierEndpoint, utils.GETRequest, headers, params, nil)
		if err != nil {
			return nil, err
		}

		apiResponse := TierTypeApiResponse{}
		err = json.Unmarshal(respBytes, &apiResponse)
		fmt.Println(apiResponse)
		if err != nil {
			return nil, err
		}

		cacheDataVal, err := json.Marshal(apiResponse.Data)
		if err != nil {
			return nil, err
		}

		// Save in cache
		err = cache.Set(redisClient, cacheKey, cacheDataVal, time.Hour*cache.TierDataTTL)
		if err != nil {
			return nil, err
		}
		tierData = apiResponse.Data
	} else {
		err := json.Unmarshal([]byte(value), &tierData)
		if err != nil {
			return nil, err
		}
	}
	return tierData, nil
}
