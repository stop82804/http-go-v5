package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	logFilePath = "network_devices.log"
	serverAddr  = "localhost:8080"
)

// DeviceInfo структура для інформації про пристрій
type DeviceInfo struct {
	DeviceName   string `json:"device_name"`
	DeviceType   string `json:"device_type"`
	IPAddress    string `json:"ip_address"`
	RoutingType  string `json:"routing_type"`
	Timestamp    string `json:"timestamp,omitempty"`
}

var (
	logMutex sync.Mutex
)

func main() {
	// Створюємо файл логів якщо не існує
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		file, err := os.Create(logFilePath)
		if err != nil {
			log.Fatal("Помилка створення лог-файлу:", err)
		}
		file.Close()
	}

	http.HandleFunc("/", handleRequest)

	fmt.Printf("Сервер запущено на http://%s/\n", serverAddr)
	fmt.Println("Доступні методи:")
	fmt.Println("  GET    / - отримати вміст лог-файлу")
	fmt.Println("  POST   / - додати новий запис до лог-файлу")
	fmt.Println("  DELETE / - очистити лог-файл")
	fmt.Println("  PUT    / - повне оновлення всіх записів")
	fmt.Println("  PATCH  / - часткове оновлення запису")

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal("Помилка запуску сервера:", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodDelete:
		handleDelete(w, r)
	case http.MethodPut:
		handlePut(w, r)
	case http.MethodPatch:
		handlePatch(w, r)
	default:
		http.Error(w, "Метод не підтримується", http.StatusMethodNotAllowed)
	}
}

// GET - повертає вміст лог-файлу
func handleGet(w http.ResponseWriter, r *http.Request) {
	logMutex.Lock()
	defer logMutex.Unlock()

	content, err := os.ReadFile(logFilePath)
	if err != nil {
		http.Error(w, "Помилка читання лог-файлу", http.StatusInternalServerError)
		log.Println("Помилка читання лог-файлу:", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if len(content) == 0 {
		w.Write([]byte("Лог-файл порожній\n"))
		return
	}
	w.Write(content)
	log.Println("GET запит: лог-файл прочитано")
}

// POST - додає новий запис до лог-файлу
func handlePost(w http.ResponseWriter, r *http.Request) {
	var device DeviceInfo

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Помилка читання тіла запиту", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &device); err != nil {
		http.Error(w, "Невірний формат JSON", http.StatusBadRequest)
		return
	}

	// Перевірка обов'язкових полів
	if device.DeviceName == "" || device.DeviceType == "" ||
	   device.IPAddress == "" || device.RoutingType == "" {
		http.Error(w, "Всі поля обов'язкові: device_name, device_type, ip_address, routing_type",
			http.StatusBadRequest)
		return
	}

	// Додаємо timestamp
	device.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	logMutex.Lock()
	defer logMutex.Unlock()

	// Записуємо в лог-файл
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Помилка відкриття лог-файлу", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	logEntry := fmt.Sprintf("[%s] Пристрій: %s, Тип: %s, IP: %s, Маршрутизація: %s\n",
		device.Timestamp, device.DeviceName, device.DeviceType,
		device.IPAddress, device.RoutingType)

	if _, err := file.WriteString(logEntry); err != nil {
		http.Error(w, "Помилка запису в лог-файл", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Запис додано до лог-файлу",
	})
	log.Printf("POST запит: додано запис - %s\n", logEntry)
}

// DELETE - очищає лог-файл
func handleDelete(w http.ResponseWriter, r *http.Request) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if err := os.Truncate(logFilePath, 0); err != nil {
		http.Error(w, "Помилка очищення лог-файлу", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Лог-файл очищено",
	})
	log.Println("DELETE запит: лог-файл очищено")
}

// PUT - повне оновлення/заміна всього вмісту лог-файлу
func handlePut(w http.ResponseWriter, r *http.Request) {
	var devices []DeviceInfo

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Помилка читання тіла запиту", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &devices); err != nil {
		http.Error(w, "Невірний формат JSON (очікується масив об'єктів)", http.StatusBadRequest)
		return
	}

	logMutex.Lock()
	defer logMutex.Unlock()

	// Очищаємо файл
	file, err := os.Create(logFilePath)
	if err != nil {
		http.Error(w, "Помилка створення лог-файлу", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Записуємо всі нові записи
	for _, device := range devices {
		device.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		logEntry := fmt.Sprintf("[%s] Пристрій: %s, Тип: %s, IP: %s, Маршрутизація: %s\n",
			device.Timestamp, device.DeviceName, device.DeviceType,
			device.IPAddress, device.RoutingType)
		file.WriteString(logEntry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Лог-файл повністю оновлено, додано %d записів", len(devices)),
	})
	log.Printf("PUT запит: лог-файл повністю оновлено (%d записів)\n", len(devices))
}

// PATCH - часткове оновлення (додаємо записи без очищення)
func handlePatch(w http.ResponseWriter, r *http.Request) {
	var devices []DeviceInfo

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Помилка читання тіла запиту", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &devices); err != nil {
		http.Error(w, "Невірний формат JSON (очікується масив об'єктів)", http.StatusBadRequest)
		return
	}

	logMutex.Lock()
	defer logMutex.Unlock()

	// Відкриваємо файл для додавання
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Помилка відкриття лог-файлу", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Додаємо нові записи
	for _, device := range devices {
		device.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		logEntry := fmt.Sprintf("[%s] Пристрій: %s, Тип: %s, IP: %s, Маршрутизація: %s\n",
			device.Timestamp, device.DeviceName, device.DeviceType,
			device.IPAddress, device.RoutingType)
		file.WriteString(logEntry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Додано %d нових записів", len(devices)),
	})
	log.Printf("PATCH запит: додано %d нових записів\n", len(devices))
}
