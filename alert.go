package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Define threshold values for temperature and humidity
const (
	TempThreshold     = 30.0 // Celsius
	HumidityThreshold = 70.0 // Percentage
)

// AlertConfig holds configuration for sending alerts
type AlertConfig struct {
	Email      string
	SMTPServer string
	SMTPPort   string
	Username   string
	Password   string
}

// InitAlerting initializes the alerting system
func InitAlerting() {
	var alertConfig AlertConfig
	yamlFile, err := os.ReadFile("conf.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &alertConfig)
	if err != nil {
		panic(err)
	}

	go monitorSensorData(alertConfig)
}

// monitorSensorData continuously checks sensor data for threshold breaches
func monitorSensorData(config AlertConfig) {
	for {
		time.Sleep(1 * time.Second)

		storeMutex.Lock()
		if len(sensorDataStore) == 0 {
			storeMutex.Unlock()
			continue
		}

		latestData := sensorDataStore[len(sensorDataStore)-1]
		storeMutex.Unlock()

		if latestData.Temp > TempThreshold {
			msg := fmt.Sprintf("High Temperature Alert! Sensor ID: %s, Temperature: %.2f°C", latestData.ID, latestData.Temp)
			sendAlert(config, "Temperature Alert", msg)
		}

		if latestData.Humidity > HumidityThreshold {
			msg := fmt.Sprintf("High Humidity Alert! Sensor ID: %s, Humidity: %.2f%%", latestData.ID, latestData.Humidity)
			sendAlert(config, "Humidity Alert", msg)
		}
	}
}

// sendAlert sends an email notification using SMTP
func sendAlert(config AlertConfig, subject, body string) {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.SMTPServer)
	msg := []byte("To: " + config.Email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	err := smtp.SendMail(config.SMTPServer+":"+config.SMTPPort, auth, config.Username, []string{config.Email}, msg)
	if err != nil {
		log.Printf("Failed to send alert: %v\n", err)
	} else {
		log.Printf("Alert sent: %s\n", body)
	}
}
