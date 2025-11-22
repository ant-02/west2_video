package database

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisInstance struct {
	client *redis.Client
}

var (
	instance redisInstance
	once     sync.Once
)

func InitRedis(ctx context.Context, addr, password string) error {
	var err error
	once.Do(func() {
		instance.client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
		})

		_, err = instance.client.Ping(ctx).Result()
	})
	return err
}

func GetRedisInstance() redisInstance {
	return instance
}

func (ri *redisInstance) Set(ctx context.Context, key string, value interface{}, expire time.Time) error {
	return ri.client.Set(ctx, key, value, time.Since(expire)).Err()
}

func (ri *redisInstance) Get(ctx context.Context, key string) (string, error) {
	return ri.client.Get(ctx, key).Result()
}

func (ri *redisInstance) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := ri.client.Exists(ctx, key).Result()
	return exists == 1, err
}

func (ri *redisInstance) RPush(ctx context.Context, key string, value interface{}) error {
	_, err := ri.client.RPush(ctx, key, value).Result()
	return err
}

func (ri *redisInstance) LRange(ctx context.Context, key string, first int64, second int64) ([]string, error) {
	return ri.client.LRange(ctx, key, first, second).Result()
}

func (ri *redisInstance) Del(ctx context.Context, key []string) error {
	_, err := ri.client.Del(ctx, key...).Result()
	return err
}

func (ri *redisInstance) LLen(ctx context.Context, key string) (int64, error) {
	return ri.client.LLen(ctx, key).Result()
}

func (ri *redisInstance) ZAdd(ctx context.Context, key string, value redis.Z) error {
	return ri.client.ZAdd(ctx, key, value).Err()
}

func (ri *redisInstance) HSet(ctx context.Context, key string, value interface{}) error {
	return ri.client.HSet(ctx, key, value).Err()
}

func (ri *redisInstance) Eval(ctx context.Context, script string, keys []string, args []interface{}) error {
	return ri.client.Eval(ctx, script, keys, args...).Err()
}
