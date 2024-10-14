package rabbit

import (
	"fmt"

	"github.com/streadway/amqp" // Библиотека для работы с RabbitMQ
)

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
