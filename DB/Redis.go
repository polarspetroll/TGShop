package DB

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

var redisdb = redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%v:%v", os.Getenv("REDISADDR"), os.Getenv("REDISPORT")),
	Password: os.Getenv("REDISPASSWD"),
	DB:       0,
})

func GetCookie(cookie string) (stat bool, username string) {
	username, err := redisdb.Get(cookie).Result()
	if err == redis.Nil {
		stat = false
	} else if err != nil {
		return false, ""
	} else {
		stat = true
	}
	return stat, username
}

func SetCookie(username string) (cookie string, err error) {
	c := make([]byte, 12)
	rand.Read(c)
	cookie = hex.EncodeToString(c)
	err = redisdb.Set(cookie, username, 10*time.Hour).Err()
	if err != nil {
		cookie = ""
		return cookie, err
	}
	return cookie, err
}

func SetCache(data QueryOutput) {
	key := fmt.Sprintf("product:%v", data.Id[0])
	redisdb.Do("hmset", key, "name", data.Name[0], "price", data.Price[0], "fname", data.Fname[0], "status", data.Stat[0])
	redisdb.Do("expire", key, "1800")
}

func GetCache(id int64) (result []interface{}, err error) {
	key := fmt.Sprintf("product:%v", id)
	r, err := redisdb.Do("hmget", key, "name", "price", "fname", "status").Result()
	if err != nil {
		result[0] = 0
		return result, err
	}
	out := make(map[string]interface{})
	out["res"] = r
	result = out["res"].([]interface{})
	return result, err
}

func SetList(key, value string) {
	redisdb.Do("sadd", key, value)
}

func GetList(key string) (result []interface{}, err error) {
	r, err := redisdb.Do("smembers", key).Result()
	if err != nil {
		result[0] = 0
		return result, err
	}
	out := make(map[string]interface{})
	out["res"] = r
	result = out["res"].([]interface{})
	return result, err

}

func Del(key string) {
	redisdb.Del(key)
}
