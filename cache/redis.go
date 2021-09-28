package cache

import (
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/token"
	"os"
	"time"
)

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

//IsTokenBlacklisted checks if token is blacklisted or not => user is logged out or not
func IsTokenBlacklisted(client *redis.Client, loginTokenMeta *token.LoginTokenMeta) bool {
	userid, _ := client.Get(loginTokenMeta.AccessUuid).Result()
	if userid != "" {
		return true
	}
	return false
}

//BlacklistToken Updates user id against uuid in cache when user logs out
func BlacklistToken(client *redis.Client, ltm *token.LoginTokenMeta) error {
	at := time.Unix(ltm.ATExpires, 0)
	now := time.Now()

	errAccess := client.Set(ltm.AccessUuid, ltm.UserID, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	return nil
}
