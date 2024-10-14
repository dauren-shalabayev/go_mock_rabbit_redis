package test

import (
	"testing"

	"mockgo/internal/worker" // Импортируй пакет с CompareData
	// Импортируй сгенерированные моки
	"mockgo/test/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCompareData(t *testing.T) {
	// Создаем контроллер gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Используем сгенерированные моки
	mockRabbitMQ := mocks.NewMockRabbitMQClient(ctrl) // Исправление
	mockRedis := mocks.NewMockRedisClient(ctrl)       // Исправление

	// Определяем тестовые данные
	queueName := "test_queue"
	redisKey := "test_key"
	expectedRabbitMsg := "test_message"
	expectedRedisData := "test_message"

	// Настраиваем поведение моков
	mockRabbitMQ.EXPECT().
		ConsumeMessage(queueName).
		Return(expectedRabbitMsg, nil)

	mockRedis.EXPECT().
		Get(gomock.Any(), redisKey).
		Return(expectedRedisData, nil)

	// Тестируем функцию CompareData
	result, err := worker.CompareData(mockRabbitMQ, mockRedis, queueName, redisKey)

	// Проверяем результат
	assert.NoError(t, err)
	assert.True(t, result)
}
