package main

// Store mutex is in main.go. Here we keep sensor data logic for modularity.
func AddSensorData(data SensorData) {
	storeMutex.Lock()
	sensorDataStore = append(sensorDataStore, data)
	storeMutex.Unlock()
}
