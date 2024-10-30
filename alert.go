package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// Define threshold values for temperature and humidity
const (
	TempThreshold     = 30.0 // Celsius
	HumidityThreshold = 70.0 // Percentage
)

// AlertConfig holds configuration for sending SMS alerts
type AlertConfig struct {
	TwilioClient   *twilio.RestClient
	FromNumber     string
	ToNumber       string
	CooldownPeriod time.Duration
}

// Map to track the last alert time for each sensor to prevent duplicate alerts
var lastAlertTime = make(map[string]time.Time)

// InitAlerting initializes the alerting system
func InitAlerting() {
	// Twilio configuration
	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: viper.GetString("TWILIO_ACCOUNT_SID"),
		Password: viper.GetString("TWILIO_AUTH_TOKEN"),
	})

	alertConfig := AlertConfig{
		TwilioClient:   twilioClient,
		FromNumber:     viper.GetString("TWILIO_FROM_NUMBER"),
		ToNumber:       viper.GetString("TWILIO_TO_NUMBER"),
		CooldownPeriod: 10 * time.Minute, // Set to 10 minutes
	}

	// Initialize the alerting system
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

		// Check temperature threshold
		if latestData.Temp > TempThreshold {
			msg := fmt.Sprintf("High Temperature Alert! Sensor ID: %s, Temperature: %.2f°C", latestData.ID, latestData.Temp)
			checkAndSendAlert(config, latestData.ID+"_temperature", msg)
		}

		// Check humidity threshold
		if latestData.Humidity > HumidityThreshold {
			msg := fmt.Sprintf("High Humidity Alert! Sensor ID: %s, Humidity: %.2f%%", latestData.ID, latestData.Humidity)
			checkAndSendAlert(config, latestData.ID+"_humidity", msg)
		}
	}
}

// checkAndSendAlert checks the cooldown period before sending an SMS alert
func checkAndSendAlert(config AlertConfig, alertKey, message string) {
	lastSent, exists := lastAlertTime[alertKey]
	if !exists || time.Since(lastSent) > config.CooldownPeriod {
		sendAlert(config, message)
		lastAlertTime[alertKey] = time.Now()
	} else {
		log.Printf("Skipped alert for %s: cooldown period active.\n", alertKey)
	}
}

// sendAlert sends an SMS notification using Twilio
func sendAlert(config AlertConfig, body string) {
	_, err := config.TwilioClient.Api.CreateMessage(&twilioApi.CreateMessageParams{
		From: &config.FromNumber,
		To:   &config.ToNumber,
		Body: &body,
	})

	if err != nil {
		log.Printf("Failed to send SMS alert: %v\n", err)
	} else {
		log.Printf("SMS Alert sent: %s\n", body)
	}
}
