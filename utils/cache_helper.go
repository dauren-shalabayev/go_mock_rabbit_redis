package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"mockgo/cache"

	"regexp"
	"strconv"

	"mockgo/models"

	"github.com/valyala/gozstd"
	"go.uber.org/zap"
)

func CacheData(taskID string, subsCache map[string]models.CacheValue, cache cache.Cache) error {
	b, err := json.Marshal(subsCache)
	if err != nil {
		return err
	}
	compB := gozstd.CompressLevel(nil, b, 1)

	return cache.Set(taskID, compB)
}

func GetPreviousResult(key int, cache cache.Cache, log *zap.Logger) (map[string]models.CacheValue, error) {
	b, err := cache.Get(strconv.Itoa(key))
	if err != nil {
		log.Error("getting prev result", zap.Error(err))
		return nil, err
	}
	decompB, err := gozstd.Decompress(nil, b)
	if err != nil {
		log.Error("decompressing prev result", zap.Error(err))
		return nil, err
	}
	var prevRes map[string]models.CacheValue
	if err = json.Unmarshal(decompB, &prevRes); err != nil {
		log.Error("unmarshalling prev result", zap.Error(err))
		return nil, err
	}
	return prevRes, nil
}

func GetRetryBlockResult(key int, cache cache.Cache, log *zap.Logger) (map[string]models.CacheValue, error) {
	b, err := cache.Get(strconv.Itoa(key) + "_retry_block")
	if err != nil {
		return nil, err
	} else {
		log.Error("Data detected at retry block cache", zap.Int("TaskId", key))
	}
	decompB, err := gozstd.Decompress(nil, b)
	if err != nil {
		log.Error("decompressing retry block result", zap.Error(err))
		return nil, err
	}
	var retryBlockRes map[string]models.CacheValue
	if err = json.Unmarshal(decompB, &retryBlockRes); err != nil {
		log.Error("unmarshalling retry block result", zap.Error(err))
		return nil, err
	}
	return retryBlockRes, nil
}

func GetRetryUnblockResult(key int, cache cache.Cache, log *zap.Logger) (map[string]models.CacheValue, error) {
	b, err := cache.Get(strconv.Itoa(key) + "_retry_unblock")
	if err != nil {
		return nil, err
	} else {
		log.Error("Data detected at retry unblock cache", zap.Int("TaskId", key))
	}
	decompB, err := gozstd.Decompress(nil, b)
	if err != nil {
		log.Error("decompressing retry unblock result", zap.Error(err))
		return nil, err
	}
	var retryUnblockRes map[string]models.CacheValue
	if err = json.Unmarshal(decompB, &retryUnblockRes); err != nil {
		log.Error("unmarshalling retry unblock result", zap.Error(err))
		return nil, err
	}

	// Выводим номера (ключи карты)
	for number := range retryUnblockRes {
		log.Info("Detected retry unblock number", zap.String("number", number))
	}

	return retryUnblockRes, nil
}

func Fill(cache cache.Cache) {
	taskID := "1000"
	taskID2 := "1001"
	unblock_retry := "1000_retry_unblock"
	block_retry := "1000_retry_block"

	// Create sample subsCache data
	subsCache := map[string]models.CacheValue{
		"77014151777": {
			Imsi:     "123456789012345",
			LacCell:  "12345-67890",
			SectorID: 1,
		},
		"77052395560": {
			Imsi:     "987654321098765",
			LacCell:  "54321-09876",
			SectorID: 2,
		},
		"77012502020": {
			Imsi:     "987654321098765",
			LacCell:  "54321-09876",
			SectorID: 2,
		},
	}

	subsCache2 := map[string]models.CacheValue{
		"77015031581": {
			Imsi:     "123456789012345",
			LacCell:  "12345-67890",
			SectorID: 1,
		},
		"77025007070": {
			Imsi:     "987654321098765",
			LacCell:  "54321-09876",
			SectorID: 2,
		},
	}

	subsCache_retry_unblock := map[string]models.CacheValue{
		"77014119100": {
			Imsi:     "123456789012345",
			LacCell:  "12345-67890",
			SectorID: 1,
		},
	}

	subsCache_retry_block := map[string]models.CacheValue{
		"77012502020": {
			Imsi:     "123456789012345",
			LacCell:  "12345-67890",
			SectorID: 1,
		},
	}

	CacheData(taskID, subsCache, cache)
	CacheData(taskID2, subsCache2, cache)
	CacheData(unblock_retry, subsCache_retry_unblock, cache)
	CacheData(block_retry, subsCache_retry_block, cache)
}

func CheckNumber(cache cache.Cache, msisdn string, logger *zap.Logger) (models.CacheValue, error) {
	ctx := context.Background()
	var cursor uint64
	numberPattern := regexp.MustCompile(`^\d+$`)

	for {
		keys, newCursor, err := cache.Scan(ctx, cursor, "*", 100)
		if err != nil {
			return models.CacheValue{}, fmt.Errorf("ошибка при сканировании ключей: %v", err)
		}

		for _, key := range keys {
			if numberPattern.MatchString(key) {
				workerKey, err := strconv.Atoi(key)
				if err != nil {
					fmt.Println("Ошибка при преобразовании ключа в формат int:", err)
					continue
				}

				cacheData, err := GetPreviousResult(workerKey, cache, logger)
				if err != nil {
					fmt.Println("Ошибка при получении данных из кеша:", err)
					continue
				}

				if val, ok := cacheData[msisdn]; ok {
					fmt.Printf("Телефонный номер найден в кеше: %s, значение: %+v\n", key, val)
					return val, nil
				}
			}
		}

		cursor = newCursor // Обновляем курсор для следующей итерации

		if cursor == 0 { // Если курсор равен 0, значит, все ключи обработаны
			break
		}
	}

	return models.CacheValue{}, nil
}

func GetBlockedNumbersForTask(cache cache.Cache, taskID int, logger *zap.Logger) (map[string]models.CacheValue, error) {

	cacheData, err := GetPreviousResult(taskID, cache, logger)
	if err != nil {
		logger.Warn("Не удалось получить данные из retry_unblock кэша, продолжение выполнения с пустым кэшем", zap.Error(err))
		cacheData = make(map[string]models.CacheValue)
	}

	cacheRetryBlockData, err := GetRetryBlockResult(taskID, cache, logger)
	if err != nil {
		logger.Warn("Не удалось получить данные из retry_block кэша, продолжение выполнения с пустым кэшем", zap.Error(err))
		cacheRetryBlockData = make(map[string]models.CacheValue)
	}

	cacheRetryUnblockData, err := GetRetryUnblockResult(taskID, cache, logger)
	if err != nil {
		logger.Warn("Не удалось получить данные из retry_unblock кэша, продолжение выполнения с пустым кэшем", zap.Error(err))
		cacheRetryUnblockData = make(map[string]models.CacheValue)
	}

	filteredData := make(map[string]models.CacheValue)
	for msisdn, val := range cacheData {
		if _, ok := cacheRetryBlockData[msisdn]; ok {
			continue
		}
		filteredData[msisdn] = val
	}

	for msisdn, val := range cacheRetryUnblockData {
		filteredData[msisdn] = val
	}

	return filteredData, nil
}
