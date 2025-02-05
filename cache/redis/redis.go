package redis

import (
	"context"
	"log"
	"time"

	"mockgo/cache"

	"github.com/go-redis/redis"
)

type Redis struct {
	Client *redis.Client
}

func (r *Redis) Delete(key string) error {
	return r.Client.Del(key).Err()
}

func (r *Redis) Set(key string, data []byte) error {
	return r.Client.Set(key, data, 0).Err()
}

func (r *Redis) Get(key string) ([]byte, error) {
	return r.Client.Get(key).Bytes()
}

func (r *Redis) SetString(key string, data string) error {
	return r.Client.Set(key, data, 0).Err()
}

func (r *Redis) GetString(key string) (string, error) {
	return r.Client.Get(key).Result()
}

func (r *Redis) GetAll() ([]string, error) {
	var keys []string
	iter := r.Client.Scan(0, "*", 0).Iterator()
	if iter.Err() != nil {
		return keys, iter.Err()
	}
	for iter.Next() {
		keys = append(keys, iter.Val())
	}
	return keys, iter.Err()
}

func (r *Redis) HSET(key string, field string, data []byte) error {
	return r.Client.HSet(key, field, data).Err()
}

func (r *Redis) HGET(key string, field string) ([]byte, error) {
	return r.Client.HGet(key, field).Bytes()
}

func (r *Redis) HDEL(key string, field string) error {
	_, err := r.Client.HDel(key, field).Result()
	return err
}

func (r *Redis) HGETALL(key string) (map[string]string, error) {
	return r.Client.HGetAll(key).Result()
}

func (r *Redis) Expire(key string, duration time.Duration) error {
	return r.Client.Expire(key, duration).Err()
}

func (r *Redis) Close() {
	err := r.Client.Close()
	if err != nil {
		log.Printf("error on redis close: %s", err)
	}
}

func (r *Redis) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return r.Client.Scan(cursor, match, count).Result()
}

func NewRedis() (cache.Cache, error) {
	// connect to redis
	redisCli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := redisCli.Ping().Result()
	return &Redis{redisCli}, err
}
