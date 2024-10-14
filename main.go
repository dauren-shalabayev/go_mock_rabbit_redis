package main

import (
	"fmt"
	"mock_go/internal/rabbit"
	"mock_go/internal/redis"
	"mock_go/internal/worker"
)

// Реализация клиента Redis

// Основная функция
func main() {
	rabbitHost := "amqp://guest:guest@localhost:5672/"
	redisHost := "localhost:6379"

	rabbitClient, err := rabbit.NewRabbitMQ(rabbitHost)
	if err != nil {
		fmt.Println("Ошибка подключения к RabbitMQ:", err)
		return
	}
	defer rabbitClient.Close()

	redisClient := redis.NewRedis(redisHost)

	// Вызов CompareData с параметрами
	result, err := worker.CompareData(rabbitClient, redisClient, "abons_log", "name")
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Printf("Результат сравнения: %v\n", result)
	}
}
