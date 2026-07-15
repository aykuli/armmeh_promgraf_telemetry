package payload

import (
	"fleet-monitor/fleet-monitor/internal/vechpath"
	"math/rand"
	"time"
)

var RobotModes = []string{"idle", "human", "teleop", "supervis", "autonom"}
var RobotMissionStatuses = []string{"none", "pause", "run", "complete", "abort"}
var baseRoomTmp = float32(22.2)

func CreatePayload(vehicleInfo VehicleInfo, coord vechpath.Coordinate) any {
	commonPayload := TelemetryPayload{
		SchemaVersion: 1,
		VehicleID:     vehicleInfo.ID,
		VehicleType:   vehicleInfo.VehicleType,
		FuelType:      vehicleInfo.FuelType,
		Timestamp:     time.Now().UnixMilli(),
	}
	metricCommon := MetricCommon{
		GpsLat:       coord.Lat,
		GpsLon:       coord.Lon,
		GpsAlt:       rand.Float64() * 5,
		SpeedKmh:     rand.Float32() * 60.0, // Speed between 0 and 15 km/h
		EngineStatus: vehicleInfo.EngineStatus,
	}

	// DIESEL
	if vehicleInfo.FuelType == "diesel" {
		oilPressureBar := (rand.Float32() * 5.0)
		fuelLevelPct := rand.Int() % 100
		engineRPM := 0
		engineHours := float32(0)
		tempC := baseRoomTmp
		if vehicleInfo.EngineStatus == "on" {
			engineRPM = rand.Int() % 4000
			engineHours = 2 + (rand.Float32() * 12)
			tempC = baseRoomTmp + (rand.Float32() * 40.0)
		} else {
			metricCommon.SpeedKmh = 0
		}

		alertnum := vehicleInfo.ID % 4
		switch alertnum {
		case 0:
			oilPressureBar = 0.1
		case 1:
			metricCommon.SpeedKmh = 65
		case 2:
			tempC = 102
		case 3:
			fuelLevelPct = 14
		}

		return DieselTelemetryPayload{
			TelemetryPayload: commonPayload,
			Metrics: DieselMetrics{
				MetricCommon:   metricCommon,
				EngineRPM:      engineRPM,
				FuelLevelPct:   fuelLevelPct,
				TempC:          float32(tempC),
				OilPressureBar: oilPressureBar,
				EngineHours:    engineHours,
			},
		}
	}
	// DIESEL END

	// ELECTRIC
	batterySocPct := rand.Int() % 100
	batteryTempC := rand.Float32() * 40
	currentA := float32(0)
	voltageV := float32(0)

	alertnum := vehicleInfo.ID % 4
	switch alertnum {
	case 0:
		batterySocPct = 9
	case 1:
		batteryTempC = 65
	case 2:
		currentA = 160
	case 3:
		metricCommon.SpeedKmh = 20
	}

	if vehicleInfo.EngineStatus == "on" {
		currentA = 120 + rand.Float32()
		voltageV = 50 - rand.Float32()*5
	}

	electricMetrics := ElectricMetrics{
		MetricCommon:  metricCommon,
		BatterySocPct: batterySocPct,
		BatteryTempC:  batteryTempC,
		CurrentA:      currentA,
		VoltageV:      voltageV,
	}

	// ROBOT

	if vehicleInfo.VehicleType == "robot" {
		tempCpuC := rand.Float32() * 50
		lteRssi := rand.Float32() * 50
		estopStatus := "off"
		rtkStatus := "fix"
		alertnum := vehicleInfo.ID % 5
		switch alertnum {
		case 0:
			tempCpuC = 90
		case 1:
			lteRssi = -90
		case 2:
			estopStatus = "on"
		case 3:
			rtkStatus = "float"
		case 4:
			rtkStatus = "none"
		}

		return RobotTelemetryPayload{
			TelemetryPayload: commonPayload,
			Metrics: RobotMetrics{
				ElectricMetrics:  electricMetrics,
				Mode:             RobotModes[vehicleInfo.ID%len(RobotModes)],
				MissionStatus:    RobotMissionStatuses[vehicleInfo.ID%len(RobotMissionStatuses)],
				MissionID:        string(rune(rand.Int() % 4000)),
				EstopStatus:      estopStatus,
				RTKStatus:        rtkStatus,
				SteeringAngleDeg: rand.Float64() * 10,
				TempCpuC:         tempCpuC,
				LteRssi:          lteRssi,
			},
		}
	}
	// ROBOT END

	return ElectricTelemetryPayload{
		TelemetryPayload: commonPayload,
		Metrics:          electricMetrics,
	}
	// ELECTRIC END

}
