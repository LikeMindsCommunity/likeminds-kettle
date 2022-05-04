package main

import (
	"github.com/nateshr/likeminds-authentication/sdk"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/home"
	"github.com/nateshr/likeminds-authentication/otp"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

var (
	client *redis.Client
	router = gin.Default()
)

func main() {
	client = cache.InitRedis()
	router.Use(ApiMiddleware(client))
	router.GET("/otp/generate", otp.GenerateOTP)
	router.GET("/otp/verify", otp.VerifyOTP)
	router.POST("/user/login", VTMValidationMiddleware(), user.Login)
	router.POST("/user/refresh", RTMValidationMiddleware(), user.Refresh)
	router.POST("/user/logout", LogoutValidationMiddleware(client), user.Logout)
	router.POST("/user/merge_account", LTMValidationMiddleware(client), user.MergeAccount)
	router.POST("/home/fetch_communities", LTMValidationMiddleware(client), home.FetchCommunities)
	router.POST("/sdk/initiate", sdk.InitiateSDK)

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

func LTMValidationMiddleware(client *redis.Client) gin.HandlerFunc {
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
