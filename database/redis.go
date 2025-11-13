package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	ctx    context.Context
)

func InitRedis(addr, password string) error {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	ctx = context.Background()
	_, err := client.Ping(ctx).Result()
	return err
}

func Set(key string, value interface{}, expire time.Time) error {
	return client.Set(ctx, key, value, time.Since(expire)).Err()
}

func Get(key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func Exists(key string) (bool, error) {
	exists, err := client.Exists(ctx, key).Result()
	return exists == 1, err
}

func RPush(key string, value interface{}) error {
	_, err := client.RPush(ctx, key, value).Result()
	return err
}

func LRange(key string, first int64, second int64) ([]string, error) {
	return client.LRange(ctx, key, first, second).Result()
}

func Del(key []string) error {
	_, err := client.Del(ctx, key...).Result()
	return err
}
