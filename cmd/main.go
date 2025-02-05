package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mockgo/internal/rabbit"
	"mockgo/internal/redis"
	"net/http"
	"strconv"
	"time"
)

// Основная функция
func main() {
	rabbitHost := "amqp://guest:guest@localhost:5672/"
	redisHost := "localhost:6379"
	redisClient := redis.NewRedis(redisHost)
	rabbitClient, err := rabbit.NewRabbitMQ(rabbitHost)
	if err != nil {
		fmt.Println("Ошибка подключения к RabbitMQ:", err)
		return
	}
	defer rabbitClient.Close()

	StartHTTPServer()
	fmt.Println("hello")

	cacheKey := strconv.Itoa(111) + "_opts"
	ctx := context.Background()
	b, err := redisClient.Get(ctx, cacheKey)
	if err != nil {
		fmt.Println(err)
	}

	var cacheData struct {
		RoutingKey string         `json:"routing_key"`
		WhiteList  []string       `json:"white_list"`
		Locations  map[string]int `json:"locations"`
	}
	fmt.Println("cachedata", cacheData)

	fmt.Println(b)

	data, err := fetchData()
	if err != nil {
		log.Println("Ошибка при получении данных:", err)

	}
	fmt.Println("response", data)

	var response struct {
		Msisdns []string `json:"msisdns"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		log.Println("Ошибка при разборе JSON:", err)

	}

	fmt.Println(response)
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
	http.HandleFunc("/api/v1/getwhitelist", getWhiteListHandler)

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

type Response struct {
	WhiteList []string `json:"msisdns"`
}

func getWhiteListHandler(w http.ResponseWriter, r *http.Request) {
	whiteList := []string{
		"77014151777", "77052396660", "77774834003", "77025007010",
	}

	response := Response{WhiteList: whiteList}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func mockResponseHandler(w http.ResponseWriter, r *http.Request) {
	mockData := map[string]interface{}{
		"msisdns": []string{
			"77773228655", "77773536688", "77778284808", "77081444648",
			"77715038467", "77715038466", "77715038468", "77717260395",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mockData)
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

func fetchData() ([]byte, error) {
	// Здесь укажи реальный URL вместо заглушки
	resp, err := http.Get("http://localhost:8051/api/v1/mockresponse")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
