package redis

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

const (
	USER_INFO_CACHE = time.Minute * 5
)

var ErrNil = redis.Nil.Error()

func Set(client *redis.Client, key string, value interface{}) error {
	_, err := client.Set(key, value, USER_INFO_CACHE).Result()
	if err != nil {
		return err
	}

	return nil
}

func Get(client *redis.Client, key string) (string, error) {
	res, err := client.Get(key).Result()
	if err != nil {
		return "", err
	}

	if err == redis.Nil {
		return "", errors.New("redis nil")
	}

	return res, nil
}

func Del(client *redis.Client, key string) error {
	_, err := client.Del(key).Result()
	if err != nil {
		return err
	}

	return nil
}
