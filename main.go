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
	router.POST("/user/logout", LTMValidationMiddleware(client), user.Logout)

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
		vtm, err := token.ExtractVerifyTokenMeta(c.Request)
		if vtm == nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.AuthenticationResponse{
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

func LTMValidationMiddleware(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ltm, err := token.ExtractLoginTokenMeta(c.Request)
		if ltm == nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.AuthenticationResponse{
				Success:      false,
				ErrorMessage: "Invalid token!",
			})
			return
		} else {
			c.Set("ltm", ltm)
			if cache.IsTokenBlacklisted(client, ltm) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, utils.AuthenticationResponse{
					Success:      false,
					ErrorMessage: "Device logged out! Please login again",
				})
				return
			}
		}
		c.Next()
	}
}
