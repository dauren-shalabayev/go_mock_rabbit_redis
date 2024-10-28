package main

import (
	"encoding/json"
	"fmt"
	"mockgo/internal/rabbit"
	"mockgo/internal/redis"
	"net/http"
	"time"
)

// Основная функция
func main() {
	rabbitHost := "amqp://guest:guest@localhost:5672/"

	rabbitClient, err := rabbit.NewRabbitMQ(rabbitHost)
	if err != nil {
		fmt.Println("Ошибка подключения к RabbitMQ:", err)
		return
	}
	defer rabbitClient.Close()

	StartHTTPServer()

	select {}
}

// Запуск HTTP-сервера
func StartHTTPServer() {
	httpServer := &http.Server{
		Addr:        ":8051",
		ReadTimeout: 120 * time.Second,
	}

	http.HandleFunc("/api/v1/task/", checkTaskHandler)
	http.HandleFunc("/api/v1/check_phone/", checkMsisdnHandler)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			fmt.Println("Ошибка при запуске сервера:", err)
		}
	}()
}

func checkTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Check request to start new task")
	w.Write([]byte("Success"))
}

// Структуры ответа
type CacheValue struct {
	Imsi     string `json:"imsi"`
	LacCell  string `json:"lac_cell"`
	SectorID int    `json:"sector_id"`
}

type CheckNumberResponse struct {
	Msisdn     string      `json:"msisdn"`
	Found      bool        `json:"found"`
	CacheValue *CacheValue `json:"cache_value,omitempty"`
}

// Обработчик checkMsisdnHandler
func checkMsisdnHandler(w http.ResponseWriter, r *http.Request) {
	msisdn := r.URL.Query().Get("msisdn")
	fmt.Println("Получен msisdn:", msisdn)
	if msisdn == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}

	redisHost := "localhost:6379"
	redisClient := redis.NewRedis(redisHost)
	msisdnData, err := redis.CheckNumber(redisClient, msisdn)
	if err != nil {
		http.Error(w, "Ошибка при проверке номера", http.StatusInternalServerError)
		fmt.Println("Ошибка:", err)
		return
	}

	isFound := msisdnData != (redis.CacheValue{}) // Проверка на нулевое значение
	response := CheckNumberResponse{
		Msisdn: msisdn,
		Found:  isFound,
	}

	// Если данные найдены, заполняем CacheValue
	if isFound {
		response.CacheValue = &CacheValue{
			Imsi:     msisdnData.Imsi,
			LacCell:  msisdnData.LacCell,
			SectorID: msisdnData.SectorID,
		}
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
	}
}
