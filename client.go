package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const serverURL = "http://localhost:8080/"

type DeviceInfo struct {
	DeviceName  string `json:"device_name"`
	DeviceType  string `json:"device_type"`
	IPAddress   string `json:"ip_address"`
	RoutingType string `json:"routing_type"`
}

func main() {
	fmt.Println("========================================")
	fmt.Println("HTTP Клієнт для тестування сервера")
	fmt.Println("========================================\n")

	// Даємо час серверу запуститися якщо потрібно
	time.Sleep(1 * time.Second)

	// 1. DELETE - Очищаємо лог перед початком
	fmt.Println("1. DELETE запит - Очищення лог-файлу")
	makeRequest("DELETE", nil)
	time.Sleep(500 * time.Millisecond)

	// 2-6. POST запити - Додаємо 5 пристроїв
	devices := []DeviceInfo{
		{
			DeviceName:  "Router-Main",
			DeviceType:  "Router",
			IPAddress:   "192.168.1.1",
			RoutingType: "Static",
		},
		{
			DeviceName:  "Switch-Core-01",
			DeviceType:  "Switch",
			IPAddress:   "192.168.1.2",
			RoutingType: "Dynamic",
		},
		{
			DeviceName:  "Server-DB",
			DeviceType:  "Server",
			IPAddress:   "192.168.1.100",
			RoutingType: "Static",
		},
		{
			DeviceName:  "Firewall-Edge",
			DeviceType:  "Firewall",
			IPAddress:   "10.0.0.1",
			RoutingType: "BGP",
		},
		{
			DeviceName:  "AP-Office-01",
			DeviceType:  "Access Point",
			IPAddress:   "192.168.2.10",
			RoutingType: "Static",
		},
	}

	for i, device := range devices {
		fmt.Printf("\n%d. POST запит - Додавання пристрою: %s\n", i+2, device.DeviceName)
		makeRequest("POST", device)
		time.Sleep(500 * time.Millisecond)
	}

	// 7. GET запит - Читаємо лог
	fmt.Println("\n7. GET запит - Отримання вмісту лог-файлу")
	makeRequest("GET", nil)
	time.Sleep(500 * time.Millisecond)

	// 8. PATCH запит - Часткове оновлення (додаємо ще 2 пристрої)
	fmt.Println("\n8. PATCH запит - Часткове оновлення (додавання пристроїв)")
	patchDevices := []DeviceInfo{
		{
			DeviceName:  "Laptop-Admin",
			DeviceType:  "Client",
			IPAddress:   "192.168.1.150",
			RoutingType: "DHCP",
		},
		{
			DeviceName:  "Camera-Entrance",
			DeviceType:  "IP Camera",
			IPAddress:   "192.168.3.20",
			RoutingType: "Static",
		},
	}
	makeRequest("PATCH", patchDevices)
	time.Sleep(500 * time.Millisecond)

	// 9. PUT запит - Повне оновлення (заміна всіх записів)
	fmt.Println("\n9. PUT запит - Повне оновлення лог-файлу")
	putDevices := []DeviceInfo{
		{
			DeviceName:  "New-Router",
			DeviceType:  "Router",
			IPAddress:   "10.10.10.1",
			RoutingType: "OSPF",
		},
		{
			DeviceName:  "New-Switch",
			DeviceType:  "Switch",
			IPAddress:   "10.10.10.2",
			RoutingType: "Static",
		},
	}
	makeRequest("PUT", putDevices)
	time.Sleep(500 * time.Millisecond)

	// 10. GET запит - Читаємо оновлений лог
	fmt.Println("\n10. GET запит - Перевірка оновленого вмісту")
	makeRequest("GET", nil)

	fmt.Println("\n========================================")
	fmt.Println("Тестування завершено успішно!")
	fmt.Println("========================================")
}

func makeRequest(method string, data interface{}) {
	var req *http.Request
	var err error

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Помилка JSON: %v\n", err)
			return
		}
		req, err = http.NewRequest(method, serverURL, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Помилка створення запиту: %v\n", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, serverURL, nil)
		if err != nil {
			fmt.Printf("Помилка створення запиту: %v\n", err)
			return
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Помилка виконання запиту: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Помилка читання відповіді: %v\n", err)
		return
	}

	fmt.Printf("Статус: %s\n", resp.Status)
	if method == "GET" {
		fmt.Printf("Відповідь:\n%s\n", string(body))
	} else {
		fmt.Printf("Відповідь: %s\n", string(body))
	}
}
