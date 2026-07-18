package payload

import (
	"fleet-app-gr/internal/vechpath"
	"fmt"
	"math/rand"
	"time"
)

var RobotModes = []string{"idle", "human", "teleop", "supervis", "autonom"}
var RobotMissionStatuses = []string{"none", "pause", "run", "complete", "abort"}
var baseRoomTmp = float32(22.2)

const charset = "abcdefghijklmnopqrstuvwxyz " + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateString(length int) string {
	return stringWithCharset(length, charset)
}

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
	event := Event{
		EventType:   "unknown",
		Severity:    "any",
		Timestamp:   commonPayload.Timestamp,
		Description: generateString(10),
		Code:        fmt.Sprintf("%04d", rand.Intn(10000)),
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
			event.Severity = Critical
			event.EventType = "oil_pressure_critical"
		case 1:
			metricCommon.SpeedKmh = 65
			event.Severity = Critical
			event.EventType = "high_speed"
		case 2:
			tempC = 102
			event.Severity = Warning
			event.EventType = "high_temp_c"
		case 3:
			fuelLevelPct = 14
			event.Severity = Warning
			event.EventType = "low_fuel_level_pc"
		}
		commonPayload.Events = []Event{event}

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
		event.Severity = Critical
		event.EventType = "low_battery_soc_pct"
	case 1:
		batteryTempC = 65
		event.Severity = Warning
		event.EventType = "high_battery_temp_c"
	case 2:
		currentA = 160
		event.Severity = Warning
		event.EventType = "anonam_current_a"
	case 3:
		metricCommon.SpeedKmh = 20
		event.Severity = Warning
		event.EventType = "high_speed"
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
			event.Severity = Critical
			event.EventType = "high_temp_cpu_c"
		case 1:
			lteRssi = -90
			event.Severity = Warning
			event.EventType = "low_lte_rssi"
		case 2:
			estopStatus = "on"
			event.Severity = Critical
			event.EventType = "active_estop_status"
		case 3:
			rtkStatus = "float"
			event.Severity = Warning
			event.EventType = "float_rtk_status"
		case 4:
			rtkStatus = "none"
			event.Severity = Critical
			event.EventType = "none_rtk_statu"
		}
		commonPayload.Events = []Event{event}

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

	commonPayload.Events = []Event{event}
	return ElectricTelemetryPayload{
		TelemetryPayload: commonPayload,
		Metrics:          electricMetrics,
	}
	// ELECTRIC END

}
