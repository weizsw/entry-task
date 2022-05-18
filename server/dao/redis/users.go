package redis

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

const (
	USER_TOKEN_FMT   = "user:token:%s"
	USER_PASSWD_FMT  = "user:passwd:%s"
	USER_INFO_CACHE  = time.Minute * 5
	USER_TOKEN_CACHE = time.Minute * 10
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

func HSet(client *redis.Client, key string, field string, val interface{}) error {
	pipe := client.Pipeline()
	pipe.HSet(key, field, val).Result()
	pipe.Expire(key, USER_TOKEN_CACHE)
	_, err := pipe.Exec()
	if err != nil {
		return err
	}

	return nil
}

func HGet(client *redis.Client, key string, field string) (string, error) {
	res, err := client.HGet(key, field).Result()
	if err != nil {
		return "", err
	}

	if err == redis.Nil {
		return "", errors.New("redis nil")
	}

	return res, nil
}
