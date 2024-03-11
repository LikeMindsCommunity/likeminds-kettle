package cache

import (
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/environment"
	"github.com/nateshr/likeminds-authentication/token"
)

func InitRedis() *redis.Client {
	//Initializing Redis
	dsn := environment.GoDotEnvVariable("REDIS_DSN")
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

// Get | get the key object value from cache storage
func Get(client *redis.Client, key string) (string, bool, error) {
	val, err := client.Get(key).Result()
	if err == redis.Nil {
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}
	return val, true, err
}

// Set | set the key with object value into cache storage, set expiration = 0 for no expiry
func Set(client *redis.Client, key string, value interface{}, expiration time.Duration) error {
	err := client.Set(key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// Keys | get all keys matching the pattern
func Keys(client *redis.Client, pattern string) ([]string, error) {
	keys, err := client.Keys(pattern).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// Delete | delete the key from cache storage
func Delete(client *redis.Client, key string) error {
	err := client.Del(key).Err()
	if err != nil {
		return err
	}
	return nil
}

// IsLTMBlacklisted checks if token is blacklisted or not => user is logged out or not
func IsLTMBlacklisted(client *redis.Client, ltm *token.LoginTokenMeta) bool {
	userUniqueID, _ := client.Get(ltm.AccessUuid).Result()
	return userUniqueID != ""
}

// IsRTMBlacklisted checks if token is blacklisted or not => user is logged out or not
func IsRTMBlacklisted(client *redis.Client, rtm *token.RefreshTokenMeta) bool {
	userUniqueID, _ := client.Get(rtm.RefreshUuid).Result()
	return userUniqueID != ""
}

// BlacklistToken Updates user id against uuid in cache when user logs out
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
