package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kelseyhightower/envconfig"
)

type MQTTconfigs struct {
	url            string
	user           string
	password       string
	keepAliveSec   int
	pingTimeoutSec int
}

type VehicleInfo struct {
	id           int
	startLat     float64
	startLon     float64
	vehicleType  string
	fuelType     string
	engineStatus string
}

func main() {
	// Базовые координаты центрального гаража
	baseLat := 55.7485
	baseLon := 37.6085
	vehicleCount := 100

	var wg sync.WaitGroup
	log.Printf("Start simulation for %d vehicles", vehicleCount)

	engineOffIndices := map[int]bool{
		8: true, 9: true, 11: true, 12: true,
		22: true, 23: true, 33: true, 34: true,
	}

	for i := 1; i <= vehicleCount; i++ {
		wg.Add(1)

		vehicleType, fuelType := getVehicleSpecs(i)
		if vehicleType == "unknown" || fuelType == "unknown" {
			continue
		}

		engineStatus := "on"
		if engineOffIndices[i] {
			engineStatus = "off"
		}

		vehicleInfo := VehicleInfo{
			id:           i,
			startLat:     baseLat + (float64(i)-1)*0.003,
			startLon:     baseLon + (float64(i)-1)*0.003,
			vehicleType:  vehicleType,
			fuelType:     fuelType,
			engineStatus: engineStatus}

		var mqttConfigs MQTTconfigs
		if err := envconfig.Process("mqtt", &mqttConfigs); err != nil {
			log.Fatal("Provide MQTT broker configuration, dureha")
		}

		go func(vehicleInfo VehicleInfo, mqttConfigs MQTTconfigs) {
			defer wg.Done()

			runVehicle(vehicleInfo, mqttConfigs)
		}(vehicleInfo, mqttConfigs)
	}

	wg.Wait()
}

/*
1..8   tractor
9..20  forklift
21..29 robot
30..45 cart
*/
func getVehicleSpecs(i int) (vehicleType string, fuelType string) {
	switch {
	case i >= 1 && i <= 8:
		vehicleType = "tractor"
		fuelType = "diesel"
	case i >= 9 && i <= 20:
		vehicleType = "forklift"
		fuelType = "electric"
	case i >= 21 && i <= 29:
		vehicleType = "robot"
		fuelType = "electric"
	case i >= 30 && i <= 45:
		vehicleType = "cart"
		fuelType = "electric"
	default:
		vehicleType = "unknown"
		fuelType = "unknown"
	}
	return
}

func runVehicle(vehicleInfo VehicleInfo, mqttConfigs MQTTconfigs) {
	options := mqtt.NewClientOptions().AddBroker(mqttConfigs.url)
	options.SetClientID(fmt.Sprintf("aynur_telemetry_gen_%d", vehicleInfo.id))
	options.SetKeepAlive(time.Duration(mqttConfigs.keepAliveSec) * time.Second)
	options.SetPingTimeout(time.Duration(mqttConfigs.pingTimeoutSec) * time.Second)
	options.SetPassword(mqttConfigs.password)
	options.SetUsername(mqttConfigs.user)

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
	}
	defer client.Disconnect(500) // milliseconds to wait for existing work to be completed

	log.Println("Connection to MQTT Broker was successfully made.")

	path := GenerateCircularPath(&vehicleInfo.startLat, &vehicleInfo.startLon)
	pointIndex := 0

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		currentCoord := path[pointIndex]
		pointIndex = (pointIndex + 1) % len(path)

		payload := TractorTelemetryPayload{
			TelemetryPayload: TelemetryPayload{
				SchemaVersion: 1,
				VehicleID:     vehicleID,
				VehicleType:   vehicleType,
				FuelType:      fuelType,
				Timestamp:     time.Now().UnixMilli(),
			},
			Metrics: TractorMetrics{
				MetricCommon: MetricCommon{
					GpsLat:       currentCoord.Lat,
					GpsLon:       currentCoord.Lon,
					GpsAlt:       17.8,
					SpeedKmh:     rand.Float64() * 15.0, // Speed between 0 and 15 km/h
					EngineStatus: "on",
				},
				EngineRPM:      1200,
				FuelLevelPct:   78 + (rand.Int() * 10.0),
				TempC:          82.5 + (rand.Float64() * 10.0),
				OilPressureBar: 2.1 + (rand.Float64() * 10.0),
				EngineHours:    1234.5 + (rand.Float64() - 0.5),
			},
		}

		jsonBytes, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			log.Fatalf("Marshalling payload error: %v", err)
			continue
		}
		topic := fmt.Sprintf("%s/%s/%d/telemetry", vehicleInfo.vehicleType, vehicleInfo.fuelType, vehicleInfo.id)

		token := client.Publish(topic, 0, false, jsonBytes)
		token.Wait()
		if token.Error() != nil {
			log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
		} else {
			log.Printf("[PUBLISHED] Target topic: %s \n payload: %+v", topic, payload)
		}
	}
}
