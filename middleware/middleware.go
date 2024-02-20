package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
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
			c.Next()
		}

		// Extract VTM info from token, internally it checks if token is valid or not
		vtm, vtmErr := token.ExtractVTM(bearerToken)

		if vtmErr == nil {
			c.Set(token.ParamVTM, vtm)
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
