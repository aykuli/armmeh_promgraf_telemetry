package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	options := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	options.SetClientID("aynur_telemetry_generator")
	options.SetKeepAlive(2 * time.Second)
	options.SetPingTimeout(1 * time.Second)
	options.SetPassword("qwertyAynur")
	options.SetUsername("aynur")

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
	}
	defer client.Disconnect(500) // milliseconds to wait for existing work to be completed

	log.Println("Connection to MQTT Broker was successfully made.")

	lat := 55.7489
	lon := 37.6087
	path := GenerateCircularPath(&lat, &lon)
	pointIndex := 0

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		currentCoord := path[pointIndex]
		pointIndex = (pointIndex + 1) % len(path)

		vehicleType := "tractor"
		fuelType := "diesel"
		vehicleId := 12
		payload := TractorTelemetryPayload{
			TelemetryPayload: TelemetryPayload{
				SchemaVersion: 1,
				VehicleID:     vehicleId,
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
		topic := fmt.Sprintf("%s/%s/%d/telemetry", vehicleType, fuelType, vehicleId)

		token := client.Publish(topic, 0, false, jsonBytes)
		token.Wait()
		if token.Error() != nil {
			log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
		} else {
			log.Printf("[PUBLISHED] Target topic: %s \n payload: %+v", topic, payload)
		}
	}
}
