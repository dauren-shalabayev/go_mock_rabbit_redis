package main

import (
	"context"
	"fmt"
)

// Интерфейсы для Redis и RabbitMQ
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
}

type RabbitMQClient interface {
	ConsumeMessage(queue string) (string, error)
}

func main() {
	fmt.Println("Hello")
}

// Реализация функции, которая берет данные из RabbitMQ и Redis и сравнивает их
func CompareData(rabbitClient RabbitMQClient, redisClient RedisClient, queueName string, redisKey string) (bool, error) {
	ctx := context.Background()

	// 1. Получаем сообщение из RabbitMQ
	rabbitMsg, err := rabbitClient.ConsumeMessage(queueName)
	if err != nil {
		return false, fmt.Errorf("ошибка при получении сообщения из RabbitMQ: %v", err)
	}

	// 2. Получаем данные из Redis
	redisData, err := redisClient.Get(ctx, redisKey)
	if err != nil {
		return false, fmt.Errorf("ошибка при получении данных из Redis: %v", err)
	}

	// 3. Сравниваем данные
	if rabbitMsg == redisData {
		return true, nil
	}
	return false, nil
}
