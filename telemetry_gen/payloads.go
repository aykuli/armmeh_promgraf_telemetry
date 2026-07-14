package main

import (
	"math/rand"
	"time"
)

var RobotModes = []string{"idle", "human", "teleop", "supervis", "autonom"}
var RobotMissionStatuses = []string{"none", "pause", "run", "complete", "autonom"}
var RobotRTKStatuses = []string{"fix", "float", "none"}

func createPayload(vehicleInfo VehicleInfo, coord Coordinate) any {
	commonPayload := TelemetryPayload{
		SchemaVersion: 1,
		VehicleID:     vehicleInfo.id,
		VehicleType:   vehicleInfo.vehicleType,
		FuelType:      vehicleInfo.fuelType,
		Timestamp:     time.Now().UnixMilli(),
	}
	metricCommon := MetricCommon{
		GpsLat:       coord.Lat,
		GpsLon:       coord.Lon,
		GpsAlt:       rand.Float64() * 5,
		SpeedKmh:     rand.Float32() * 15.0, // Speed between 0 and 15 km/h
		EngineStatus: vehicleInfo.engineStatus,
	}

	if vehicleInfo.fuelType == "diesel" {
		engineRPM := 0
		if vehicleInfo.engineStatus == "on" {
			engineRPM = rand.Int() * 4000
		}

		return DieselTelemetryPayload{
			TelemetryPayload: commonPayload,
			Metrics: DieselMetrics{
				MetricCommon:   metricCommon,
				EngineRPM:      engineRPM,
				FuelLevelPct:   rand.Int() * 100,
				TempC:          82.5 + (rand.Float32() * 10.0),
				OilPressureBar: 2.1 + (rand.Float32() * 10.0),
				EngineHours:    1234.5 + (rand.Float32() - 0.5),
			},
		}
	}
	currentA := float32(0)
	voltageV := float32(0)
	if vehicleInfo.engineStatus == "on" {
		currentA = 120 + rand.Float32()
		voltageV = 50 - rand.Float32()*5
	}

	electricMetrics := ElectricMetrics{
		MetricCommon:  metricCommon,
		BatterySocPct: rand.Int() * 100,
		BatteryTempC:  rand.Float32() * 40,
		CurrentA:      currentA,
		VoltageV:      voltageV,
	}

	if vehicleInfo.vehicleType == "robot" {
		return RobotTelemetryPayload{
			TelemetryPayload: commonPayload,
			Metrics: RobotMetrics{
				ElectricMetrics:  electricMetrics,
				Mode:             RobotModes[vehicleInfo.id%len(RobotModes)],
				MissionStatus:    RobotMissionStatuses[vehicleInfo.id%len(RobotMissionStatuses)],
				MissionID:        string(rand.Int() % 4000),
				EstopStatus:      vehicleInfo.engineStatus,
				RTKStatus:        RobotRTKStatuses[vehicleInfo.id%len(RobotRTKStatuses)],
				SteeringAngleDeg: rand.Float64() * 10,
				TempCpuC:         rand.Float32() * 50,
				LteRssi:          rand.Float32() * 50,
			},
		}
	}

	return ElectricTelemetryPayload{
		TelemetryPayload: commonPayload,
		Metrics:          electricMetrics,
	}
}
