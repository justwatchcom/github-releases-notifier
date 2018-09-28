package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"os"
)

//Client redis Client
var Client *redis.Client

//SetKey set key in redis
func SetKey(key string, value string) {
	var err error

	if GetValue(key) != "" {
		Client.Del(key)
	}

	err = Client.Set(key, value, 0).Err()
	check(err)
}

//GetValue get value from redis
func GetValue(key string) string {
	return Client.Get(key).Val()
}

//СonnectToRedis connect to redis
func СonnectToRedis() redis.Client {
	var host = os.Getenv("REDIS_HOST")
	var port = os.Getenv("REDIS_PORT")

	Client = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       11,
	})

	pong, err := Client.Ping().Result()
	fmt.Println(pong, err)

	return *Client
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
