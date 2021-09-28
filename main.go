package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/otp"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
	"log"
	"net/http"
)

var (
	client *redis.Client
	router = gin.Default()
)

func main() {
	client = cache.InitRedis()
	router.Use(ApiMiddleware(client))
	router.GET("/otp/verify", otp.VerifyOTP)
	router.POST("/user/login", VTMValidationMiddleware(), user.Login)
	router.POST("/user/logout", AccessLTMValidationMiddleware(client), user.Logout)
	router.POST("/user/refresh", RefreshLTMValidationMiddleware(), user.Refresh)

	log.Fatal(router.Run(":8080"))
}

// ApiMiddleware will add the db connection to the context
func ApiMiddleware(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redis_client", client)
		c.Next()
	}
}

func VTMValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		vtm, err := token.ExtractVTM(c.Request)
		if vtm == nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusOK, utils.AuthenticationResponse{
				Success:      false,
				ErrorMessage: "Invalid token!",
			})
			return
		} else {
			c.Set("vtm", vtm)
		}
		c.Next()
	}
}

func AccessLTMValidationMiddleware(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ltm, err := token.ExtractAccessLTM(c.Request)
		if ltm == nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusOK, utils.AuthenticationResponse{
				Success:      false,
				ErrorMessage: "Invalid token!",
			})
			return
		} else {
			c.Set("access_ltm", ltm)
			if cache.IsTokenBlacklisted(client, ltm) {
				c.AbortWithStatusJSON(http.StatusOK, utils.AuthenticationResponse{
					Success:      false,
					ErrorMessage: "Device logged out! Please login again",
				})
				return
			}
		}
		c.Next()
	}
}

func RefreshLTMValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ltm, err := token.ExtractRefreshLTM(c.Request)
		if ltm == nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusOK, utils.AuthenticationResponse{
				Success:      false,
				ErrorMessage: "Invalid token!",
			})
			return
		} else {
			c.Set("refresh_ltm", ltm)
		}
		c.Next()
	}
}
