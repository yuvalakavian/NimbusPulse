package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// SensorData holds environmental data from a sensor
type SensorData struct {
	ID        string  `json:"id"`
	Timestamp int64   `json:"timestamp"`
	Temp      float64 `json:"temperature"`
	Humidity  float64 `json:"humidity"`
}

var (
	sensorDataStore []SensorData
	storeMutex      sync.Mutex
	clients         = make(map[*websocket.Conn]bool)
	upgrader        = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	
	// Initialize the alerting system
	InitAlerting()
	// HTTP endpoint to receive sensor data
	http.HandleFunc("/sensor", receiveSensorData)

	// WebSocket endpoint for real-time data streaming
	http.HandleFunc("/ws", handleWebSocketConnections)

	// Run a goroutine to broadcast updates to clients every second
	go broadcastSensorData()

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// receiveSensorData handles incoming sensor data via HTTP
func receiveSensorData(w http.ResponseWriter, r *http.Request) {
	var data SensorData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid data format", http.StatusBadRequest)
		return
	}
	data.Timestamp = time.Now().Unix()

	// Store data
	storeMutex.Lock()
	sensorDataStore = append(sensorDataStore, data)
	storeMutex.Unlock()

	fmt.Printf("Received data: %+v\n", data)
	w.WriteHeader(http.StatusNoContent)
}

// handleWebSocketConnections upgrades HTTP to WebSocket and registers clients
func handleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			delete(clients, ws)
			break
		}
	}
}

// broadcastSensorData sends the latest sensor data to all connected clients
func broadcastSensorData() {
	for {
		time.Sleep(1 * time.Second)

		storeMutex.Lock()
		if len(sensorDataStore) == 0 {
			storeMutex.Unlock()
			continue
		}

		latestData := sensorDataStore[len(sensorDataStore)-1]
		storeMutex.Unlock()

		dataJSON, err := json.Marshal(latestData)
		if err != nil {
			log.Println("Error marshaling data:", err)
			continue
		}

		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, dataJSON); err != nil {
				log.Printf("WebSocket error: %v\n", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
