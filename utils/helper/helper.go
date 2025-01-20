package helper

import (
	"math"
	"time"
)

// TimestampToDate timeConvert from unix epoch to timestamp
func TimestampToDate(timestamp int64) time.Time {
	location, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		panic(err)
	}
	return time.Unix(timestamp/1e3, 0).In(location)
}

// CalculateDistanceBetween Calculate distance between two points using Haversine formula
func CalculateDistanceBetween(LatitudeA, LongitudeA, LatitudeB, LongitudeB float64) float64 {
	const R = 6371e3
	phi1 := LatitudeA * math.Pi / 180 // φ, λ in radians
	phi2 := LatitudeB * math.Pi / 180
	deltaPhi := (LatitudeB - LatitudeA) * math.Pi / 180
	deltaLambda := (LongitudeB - LongitudeA) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := (R * c) / 1000
	return d
}
