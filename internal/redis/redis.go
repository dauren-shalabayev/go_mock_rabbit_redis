package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8" // Библиотека для работы с Redis
	// Библиотека для работы с RabbitMQ
)

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
}

type RedisClientImpl struct {
	client *redis.Client
}

// Создаем нового клиента Redis
func NewRedis(host string) *RedisClientImpl {
	// Подключаемся к Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     host, // Адрес Redis-сервера (например, "localhost:6379")
		Password: "",   // Пароль, если требуется
		DB:       0,    // Используемая база данных (0 по умолчанию)
	})

	return &RedisClientImpl{
		client: rdb,
	}
}

// Реализация метода Get для получения данных по ключу из Redis
func (r *RedisClientImpl) Get(ctx context.Context, key string) (string, error) {
	// Получаем данные по ключу из Redis
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("ключ не найден в Redis: %s", key)
	} else if err != nil {
		return "", fmt.Errorf("ошибка при получении данных из Redis: %v", err)
	}

	// Возвращаем полученное значение
	return val, nil
}
