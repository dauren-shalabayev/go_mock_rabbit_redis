package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
	SetString(key string, data string) error
	GetString(key string) (string, error)
	Delete(key string) error
	GetAll() ([]string, error)
	HSET(key string, field string, data []byte) error
	HGET(key string, field string) ([]byte, error)
	HDEL(key string, fields string) error
	HGETALL(key string) (map[string]string, error)
	Expire(key string, duration time.Duration) error
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)
	Close()
}
