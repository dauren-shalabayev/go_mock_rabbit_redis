package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

// RabbitMQ клиент
type RabbitMQ struct {
	conn *amqp.Connection
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &RabbitMQ{conn: conn}, nil
}

func (r *RabbitMQ) ConsumeMessage(queue string) (string, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return "", err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queue, // имя очереди
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return "", err
	}

	for msg := range msgs {
		return string(msg.Body), nil
	}
	return "", fmt.Errorf("нет сообщений")
}

// Redis клиент
type Redis struct {
	client *redis.Client
}

func NewRedis(addr string) *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("ключ не найден")
	} else if err != nil {
		return "", err
	}
	return val, nil
}
