package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/go-redis/redis/v8" // Библиотека для работы с Redis
	"github.com/valyala/gozstd"
	// Библиотека для работы с RabbitMQ
)

type RedisClient interface {
	Set(ctx context.Context, key string, data []byte) error
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

func (r *RedisClientImpl) Set(ctx context.Context, key string, data []byte) (string, error) {
	// Получаем данные по ключу из Redis
	val, err := r.client.Set(ctx, key, data, 0).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("ключ не найден в Redis: %s", key)
	} else if err != nil {
		return "", fmt.Errorf("ошибка при получении данных из Redis: %v", err)
	}

	// Возвращаем полученное значение
	return val, nil
}

type CacheValue struct {
	Imsi     string
	LacCell  string
	SectorID int
}

type Cache struct {
	LacCell     string
	SectorID    int
	Subscribers map[string]string
}

func CacheData(taskID string, r *RedisClientImpl) error {

	newRes := make(map[string]CacheValue)
	newRes["77052395560"] = CacheValue{
		Imsi:     "77052395560900",
		LacCell:  "100-300",
		SectorID: 1020,
	}
	b, err := json.Marshal(newRes)
	if err != nil {
		return err
	}
	compB := gozstd.CompressLevel(nil, b, 1)
	ctx := context.Background()
	val, err := r.Set(ctx, taskID, compB)
	fmt.Println(val)
	return nil
}

func CacheData2(taskID string, r *RedisClientImpl) error {

	newRes := make(map[string]CacheValue)
	newRes["77014151777"] = CacheValue{
		Imsi:     "77014151777999",
		LacCell:  "100-200",
		SectorID: 1021,
	}
	b, err := json.Marshal(newRes)
	if err != nil {
		return err
	}
	compB := gozstd.CompressLevel(nil, b, 1)
	ctx := context.Background()
	val, err := r.Set(ctx, taskID, compB)
	fmt.Println(val)
	return nil
}

func GetPreviousResult(key int, r *RedisClientImpl) (map[string]CacheValue, error) {
	ctx := context.Background()
	b, err := r.Get(ctx, strconv.Itoa(key))
	if err != nil {
		return nil, err
	}

	// Convert the string to a byte slice
	bBytes := []byte(b)

	decompB, err := gozstd.Decompress(nil, bBytes)
	if err != nil {
		return nil, err
	}

	var prevRes map[string]CacheValue
	if err = json.Unmarshal(decompB, &prevRes); err != nil {
		return nil, err
	}
	return prevRes, nil
}

func CheckNumber(r *RedisClientImpl, msisdn string) (CacheValue, error) {
	ctx := context.Background() // Создаем контекст
	var cursor uint64
	numberPattern := regexp.MustCompile(`^\d+$`) // Регулярное выражение для поиска ключей, содержащих только числа

	for {
		keys, newCursor, err := r.client.Scan(ctx, cursor, "*", 100).Result() // Используем "100" для ограничения числа возвращаемых ключей
		if err != nil {
			return CacheValue{}, fmt.Errorf("ошибка при сканировании ключей: %v", err)
		}

		for _, key := range keys {
			if numberPattern.MatchString(key) { // Проверяем, соответствует ли ключ шаблону
				workerKey, err := strconv.Atoi(key)
				if err != nil {
					fmt.Println("Ошибка при преобразовании ключа в формат int:", err)
					continue // Переход к следующему ключу
				}

				cacheData, err := GetPreviousResult(workerKey, r)
				if err != nil {
					fmt.Println("Ошибка при получении данных из кеша:", err)
					continue // Переход к следующему ключу
				}

				if val, ok := cacheData[msisdn]; ok {
					fmt.Printf("Телефонный номер найден в кеше: %s, значение: %+v\n", key, val)
					return val, nil // Возвращаем найденное значение
				}
			}
		}

		cursor = newCursor // Обновляем курсор для следующей итерации

		if cursor == 0 { // Если курсор равен 0, значит, все ключи обработаны
			break
		}
	}

	return CacheValue{}, nil // Возвращаем сообщение о том, что номер не найден
}
