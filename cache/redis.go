package cache

import (
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/token"
	"net/http"
	"os"
	"strconv"
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

//TODO
func FetchLogoutUserID(client *redis.Client, loginTokenMeta *token.LoginTokenMeta) (uint64, error) {
	userid, err := client.Get(loginTokenMeta.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

// CreateLogoutUserID When logout api is called
func CreateLogoutUserID(client *redis.Client, r *http.Request) error {
	ltm, err := token.ExtractLoginTokenMeta(r)
	if err != nil {
		return err
	}

	at := time.Unix(ltm.AtExpires, 0)
	now := time.Now()

	errAccess := client.Set(ltm.AccessUuid, ltm.UserID, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	return nil
}
