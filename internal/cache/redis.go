package cache

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/environment"
	"github.com/nateshr/likeminds-authentication/internal/logging"
)

func InitRedis() *redis.Client {
	//Initializing Redis
	dsn := environment.GoDotEnvVariable("REDIS_DSN")
	password := environment.GoDotEnvVariable("REDIS_PASSWORD")

	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: dsn,
	})

	serverEnv := environment.GoDotEnvVariable(environment.EnvServerEnviornment)
	if serverEnv == "load" {
		client.Options().Password = password
		client.Options().TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

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
		logging.Info(fmt.Sprint("cache miss for key: ", key))
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}
	logging.Info(fmt.Sprint("cache hit for key: ", key))
	return val, true, err
}

// Set | set the key with object value into cache storage, set expiration = 0 for no expiry
func Set(client *redis.Client, key string, value interface{}, expiration time.Duration) error {
	err := client.Set(key, value, expiration).Err()
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("cache set for key: %s with expiry: %v", key, expiration.String()))
	return nil
}

// GetKeys | get all keys matching the pattern
func GetKeys(client *redis.Client, pattern string) ([]string, error) {
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

	logging.Info(fmt.Sprint("cache deleted for key: ", key))
	return nil
}

// GetFromMultipleKeys | get multiple keys from cache storage
func GetFromMultipleKeys(client *redis.Client, keys ...string) ([]interface{}, error) {
	val, err := client.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}
	for i, v := range val {
		if v == nil {
			logging.Info(fmt.Sprint("cache miss for key: ", keys[i]))
		} else {
			logging.Info(fmt.Sprint("cache hit for key: ", keys[i]))
		}
	}
	return val, nil
}

// SetMultipleValues | set multiple keys with object value into cache storage
func SetMultipleValues(client *redis.Client, values ...interface{}) error {
	err := client.MSet(values...).Err()
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprint("cache set for keys: ", values))
	return nil
}

// Increment | increment the key value by 1 or set to 1 if key does not exist
func Increment(client *redis.Client, key string) error {
	err := client.Incr(key).Err()
	if err != nil {
		logging.Error(fmt.Sprint("cache increment failed for key: ", key, " | err: ", err.Error()))
		return err
	}

	logging.Info(fmt.Sprint("cache incremented for key: ", key))
	return nil
}

// Expire | set the expiration time for the key
func ExpireAt(client *redis.Client, key string, expiration time.Time) error {
	err := client.ExpireAt(key, expiration).Err()
	if err != nil {
		logging.Error(fmt.Sprint("cache expiry set failed for key: ", key, " | err: ", err.Error()))
		return err
	}

	logging.Info(fmt.Sprintf("cache expiry set for key: %s with expiry: %v", key, expiration.String()))
	return nil
}

// FetchTTL | get the TTL for key
func FetchTTL(client *redis.Client, key string) (time.Duration, error) {
	ttl, err := client.TTL(key).Result()
	if err != nil {
		logging.Error(fmt.Sprint("cache ttl fetch failed for key: ", key, " | err: ", err.Error()))
		return 0, err
	}

	logging.Info(fmt.Sprintf("cache ttl fetch for key: %s with ttl: %v", key, ttl.String()))
	return ttl, nil
}

// IsLTMBlacklisted checks if token is blacklisted or not => user is logged out or not
func IsLTMBlacklisted(client *redis.Client, ltm *constants.LoginTokenMeta) bool {
	userUniqueID, _ := client.Get(ltm.AccessUuid).Result()
	return userUniqueID != ""
}

// IsRTMBlacklisted checks if token is blacklisted or not => user is logged out or not
func IsRTMBlacklisted(client *redis.Client, rtm *constants.RefreshTokenMeta) bool {
	userUniqueID, _ := client.Get(rtm.RefreshUuid).Result()
	return userUniqueID != ""
}

// BlacklistToken Updates user id against uuid in cache when user logs out
func BlacklistToken(client *redis.Client, ltm *constants.LoginTokenMeta, rtm *constants.RefreshTokenMeta) error {
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

// DeleteCacheByKeyPatterns deletes keys from Redis based on key patterns.
func DeleteCacheByKeyPatterns(redisClient *redis.Client, keyPatterns []string) (map[string][]string, error) {
	// Initialize a map to store key patterns with their corresponding matched keys.
	result := make(map[string][]string)

	for _, keyPattern := range keyPatterns {
		// Get all keys matching the pattern
		keys, err := GetKeys(redisClient, keyPattern)
		if err != nil {
			logging.Error(fmt.Sprintf("Error fetching keys for pattern %s: %v", keyPattern, err))
			return nil, err
		}

		// Store the keys in the result map, even if no keys are found.
		result[keyPattern] = keys

		if len(keys) == 0 {
			logging.Info(fmt.Sprintf("No cache keys found for pattern: %s", keyPattern))
			continue
		}
		// Delete all keys
		for _, key := range keys {
			err = Delete(redisClient, key)
			if err != nil {
				logging.Error(fmt.Sprintf("Error deleting key %s: %v", key, err))
				return nil, err
			}
			logging.Info(fmt.Sprintf("Successfully deleted cache for key: %s", key))
		}

	}
	return result, nil
}
