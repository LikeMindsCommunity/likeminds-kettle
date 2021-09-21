package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/otp"
	"github.com/nateshr/likeminds-authentication/user"
	"log"
)

var (
	client *redis.Client
	router = gin.Default()
)

func main() {
	client = cache.InitRedis()
	router.Use(ApiMiddleware(client))
	router.GET("/otp/verify", otp.VerifyOTP)
	router.POST("/user/login", user.Login)

	log.Fatal(router.Run(":8080"))
}

// ApiMiddleware will add the db connection to the context
func ApiMiddleware(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("client", client)
		c.Next()
	}
}
