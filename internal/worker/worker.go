package worker

import (
	"context"
	"fmt"
	"mock_go/internal/rabbit"
	"mock_go/internal/redis"
)

// Реализация функции, которая берет данные из RabbitMQ и Redis и сравнивает их
func CompareData(rabbitClient rabbit.RabbitMQClient, redisClient redis.RedisClient, queueName string, redisKey string) (bool, error) {
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

	fmt.Println("Данные из Redis:", redisData)

	// 3. Сравниваем данные
	if rabbitMsg == redisData {
		return true, nil
	}
	return false, nil
}
