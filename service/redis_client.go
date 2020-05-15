package service

import (
	"os"

	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		logrus.Fatal(err)
	}

	return client
}
