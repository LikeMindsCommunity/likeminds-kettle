package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/community"
	"github.com/nateshr/likeminds-authentication/home"
	"github.com/nateshr/likeminds-authentication/otp"
	"github.com/nateshr/likeminds-authentication/sdk"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
	"github.com/nateshr/likeminds-authentication/web"
)

var (
	client *redis.Client
	router = gin.Default()
)

func main() {
	client = cache.InitRedis()
	router.Use(ApiMiddleware(client))
	router.GET("", web.Home)
	router.GET("/otp/generate", otp.GenerateOTP)
	router.GET("/otp/verify", otp.VerifyOTP)
	router.POST("/user/login", VTMValidationMiddleware(), user.Login)
	router.POST("/user/refresh", RTMValidationMiddleware(), user.Refresh)
	router.POST("/user/logout", LogoutValidationMiddleware(client), user.Logout)
	router.POST("/user/merge_account", LTMValidationMiddleware(client, true), user.MergeAccount)
	router.POST("/home/fetch_communities", LTMValidationMiddleware(client, true), home.FetchCommunities)
	router.POST("/sdk/initiate", APIKeyValidationMiddleware(), LTMValidationMiddleware(client, false), sdk.InitiateSDK)
	router.POST("/chatroom/schedule_follow", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.ScheduleFollow)
	router.POST("/chatroom/pin", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.PinChatroom)
	router.GET("/chatroom/get_tagging_list", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.GetTaggingList)
	router.POST("/user/bot", APIKeyValidationMiddleware(), user.CreateBot)
	router.PUT("/user/bot", APIKeyValidationMiddleware(), user.EditBot)
	router.GET("/user/bot", APIKeyValidationMiddleware(), user.GetBot)
	router.POST("/sdk/project", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), sdk.CreateProject)
	router.GET("/sdk/project", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), sdk.GetProject)
	router.DELETE("/sdk/project", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), sdk.DeleteProject)
	router.GET("/chatroom", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.GetChatroom)
	router.POST("/chatroom", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.CreateChatroom)
	router.PUT("/chatroom", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.EditChatroom)
	router.GET("/community/member", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetMember)
	router.GET("/chatroom/participants", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.GetParticipants)
	router.POST("/chatroom/participants", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.AddParticipants)
	log.Fatal(router.Run(":8080"))
}

// ApiMiddleware will add the db connection to the context
func ApiMiddleware(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(cache.ParamRedisClient, client)
		c.Next()
	}
}

func VTMValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Extract VTM from token, internally it checks if token is valid or not
		vtm, err := token.ExtractVTM(c.Request.Header.Get(token.HeaderAuthorization))
		if vtm == nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidVTM,
			})
			return
		} else {
			//If valid, set "vtm" in context, to be used in later APIs
			c.Set(token.ParamVTM, vtm)
		}
		c.Next()
	}
}

func LTMValidationMiddleware(client *redis.Client, emptyBearerTokenCheck bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get(token.HeaderAuthorization)
		//If bearer token is empty, let it pass through
		if !emptyBearerTokenCheck && len(bearerToken) == 0 {
			c.Next()
		}
		//Extract LTM from token, internally it checks if token is valid or not
		ltm, err := token.ExtractLTM(bearerToken)
		if ltm == nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidLTM,
			})
			return
		} else {
			//Check if LTM is black listed or not
			if cache.IsLTMBlacklisted(client, ltm) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
					Success:      false,
					ErrorMessage: utils.ErrorDeviceLoggedOut,
				})
				return
			}
			//If valid and not blacklisted, set "ltm" in context, to be used in later APIs
			c.Set(token.ParamLTM, ltm)
		}
		c.Next()
	}
}

func RTMValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Extract RTM from token, internally it checks if token is valid or not
		rtm, err := token.ExtractRTM(c.Request.Header.Get(token.HeaderAuthorization))
		if rtm == nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidRTM,
			})
			return
		} else {
			//Check if RTM is black listed or not
			if cache.IsRTMBlacklisted(client, rtm) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
					Success:      false,
					ErrorMessage: utils.ErrorDeviceLoggedOut,
				})
				return
			}
			//If valid and not blacklisted, set "rtm" in context, to be used in later APIs
			c.Set(token.ParamRTM, rtm)
		}
		c.Next()
	}
}

func LogoutValidationMiddleware(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Extract LTM from token, internally it checks if token is valid or not
		ltm, err := token.ExtractLTM(c.Request.Header.Get(token.HeaderAuthorization))
		if ltm == nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidLTM,
			})
			return
		} else {
			//Check if LTM is black listed or not
			if cache.IsLTMBlacklisted(client, ltm) {
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
				log.Print(err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
					Success:      false,
					ErrorMessage: token.ErrorInvalidRTM,
				})
				return
			} else {
				//Check if RTM is black listed or not
				if cache.IsRTMBlacklisted(client, rtm) {
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

const SDKAuthenticateEndPoint = "/api/sdk/authenticate"
const ResponseCommunityId = "community_id"

func APIKeyValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Create internal API client
		client := api_client.NewAPIClient()
		options := api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + SDKAuthenticateEndPoint,
			CustomHeaders: utils.CreateHeaders(c, ""),
		}
		//Send request
		respBytes, err := client.GetRequest(&options)
		if err != nil {
			//If API fails or any other error
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: err.Error(),
			})
			return
		}
		//Parse response
		var apiCR api_client.APIClientResponse
		err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
		if err != nil {
			//Internal unmarshal error
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: err.Error(),
			})
			return
		}

		if !apiCR.Success {
			//If api/sdk/authenticate returns success as false
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: utils.ErrorInvalidAPIKey,
			})
			return
		}
		c.Set(ResponseCommunityId, apiCR.Response[ResponseCommunityId])
		c.Next()
	}
}
