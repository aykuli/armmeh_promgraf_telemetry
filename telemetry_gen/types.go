package main

type TelemetryPayload struct {
	SchemaVersion int    `json:"schema_version"`
	VehicleID     int    `json:"vehicle_id"`
	VehicleType   string `json:"vehicle_type"`
	FuelType      string `json:"fuel_type"`
	Timestamp     int64  `json:"timestamp"`
}

type TractorTelemetryPayload struct {
	TelemetryPayload
	Metrics TractorMetrics `json:"metrics"`
}

type MetricCommon struct {
	GpsLat       float64 `json:"gps_lat"`
	GpsLon       float64 `json:"gps_lon"`
	GpsAlt       float64 `json:"gps_alt"`
	SpeedKmh     float64 `json:"speed_kmh"`
	EngineStatus string  `json:"engine_status"`
}

type TractorMetrics struct {
	MetricCommon
	EngineRPM      int     `json:"engine_rpm"`
	FuelLevelPct   int     `json:"fuel_level_pct"`
	TempC          float64 `json:"temp_c"`
	OilPressureBar float64 `json:"oil_pressure_bar"`
	EngineHours    float64 `json:"engine_hours"`
}
