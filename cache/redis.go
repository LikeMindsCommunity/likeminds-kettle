package cache

import (
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/token"
	"os"
	"time"
)

const ParamRedisClient = "redis_client"

func InitRedis() *redis.Client {
	//Initializing Redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return client
}

//IsLTMBlacklisted checks if token is blacklisted or not => user is logged out or not
func IsLTMBlacklisted(client *redis.Client, ltm *token.LoginTokenMeta) bool {
	userUniqueID, _ := client.Get(ltm.AccessUuid).Result()
	if userUniqueID != "" {
		return true
	}
	return false
}

//IsRTMBlacklisted checks if token is blacklisted or not => user is logged out or not
func IsRTMBlacklisted(client *redis.Client, rtm *token.RefreshTokenMeta) bool {
	userUniqueID, _ := client.Get(rtm.RefreshUuid).Result()
	if userUniqueID != "" {
		return true
	}
	return false
}

//BlacklistToken Updates user id against uuid in cache when user logs out
func BlacklistToken(client *redis.Client, ltm *token.LoginTokenMeta, rtm *token.RefreshTokenMeta) error {
	atExpires := time.Unix(ltm.AccessTokenExpires, 0)
	rtExpires := time.Unix(rtm.RefreshTokenExpires, 0)
	now := time.Now()

	errAccess := client.Set(ltm.AccessUuid, ltm.UserUniqueID, atExpires.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(rtm.RefreshUuid, ltm.UserUniqueID, rtExpires.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
