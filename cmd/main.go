package main

import (
	"fmt"
	"mockgo/internal/rabbit"
	"mockgo/internal/redis"
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
	type CacheValue struct {
		Imsi     string
		LacCell  string
		SectorID int
	}
	redisClient := redis.NewRedis(redisHost)

	// Call CacheData function
	//redis.CacheData("1111", redisClient)
	//redis.CacheData2("1112", redisClient)

	redis.CheckNumber(redisClient, "77014151777")
	// Вызов CompareData с параметрами
	// result, err := worker.CompareData(rabbitClient, redisClient, "abons_log", "name")
	// if err != nil {
	// 	fmt.Printf("Ошибка: %v\n", err)
	// } else {
	// 	fmt.Printf("Результат сравнения: %v\n", result)
	// }
}
