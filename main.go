package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8" // Библиотека для работы с Redis
	"github.com/streadway/amqp"    // Библиотека для работы с RabbitMQ
)

// Интерфейсы для Redis и RabbitMQ
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
}

type RabbitMQClient interface {
	ConsumeMessage(queue string) (string, error)
}

// Реализация клиента RabbitMQ
type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(uri string) (*RabbitClient, error) {
	// Создаем подключение к RabbitMQ
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к RabbitMQ: %v", err)
	}

	// Открываем канал
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("ошибка открытия канала RabbitMQ: %v", err)
	}

	return &RabbitClient{
		conn:    conn,
		channel: channel,
	}, nil
}

func (r *RabbitClient) ConsumeMessage(queue string) (string, error) {
	msgs, err := r.channel.Consume(
		queue, "", true, false, false, false, nil,
	)
	if err != nil {
		return "", fmt.Errorf("ошибка получения сообщения: %v", err)
	}

	for msg := range msgs {
		message := string(msg.Body)
		fmt.Println("Полученное сообщение:", message)
		return message, nil
	}

	return "", fmt.Errorf("сообщение не найдено")
}

func (r *RabbitClient) Close() {
	r.channel.Close()
	r.conn.Close()
}

// Реализация клиента Redis
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

// Основная функция
func main() {
	rabbitHost := "amqp://guest:guest@localhost:5672/"
	redisHost := "localhost:6379"

	rabbitClient, err := NewRabbitMQ(rabbitHost)
	if err != nil {
		fmt.Println("Ошибка подключения к RabbitMQ:", err)
		return
	}
	defer rabbitClient.Close()

	redisClient := NewRedis(redisHost)

	// Вызов CompareData с параметрами
	result, err := CompareData(rabbitClient, redisClient, "abons_log", "name")
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Printf("Результат сравнения: %v\n", result)
	}
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

	fmt.Println("Данные из Redis:", redisData)

	// 3. Сравниваем данные
	if rabbitMsg == redisData {
		return true, nil
	}
	return false, nil
}
