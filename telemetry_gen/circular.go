package main

import (
	"math"
)

type Coordinate struct {
	Lat float64
	Lon float64
}

func GenerateCircularPath(centerLat, centerLon *float64) []Coordinate {
	const totalPoints = 100
	fbLat := 55.7489
	fbLon := 37.6087
	const radius = 0.001 // 100 m

	if centerLat == nil {
		centerLat = &fbLat
	}
	if centerLon == nil {
		centerLon = &fbLon
	}

	path := make([]Coordinate, totalPoints)

	for i := range totalPoints {
		angle := (float64(i) / float64(totalPoints)) * 2.0 * math.Pi

		path[i] = Coordinate{
			Lat: *centerLat + (radius * math.Sin(angle)),
			Lon: *centerLon + (radius * math.Cos(angle)),
		}
	}

	return path
}
