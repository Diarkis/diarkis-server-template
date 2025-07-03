package common

import "math"

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

const (
	equatorRadius        = 6378137.0
	polesRadius          = 6356752.314245
	equatorRadiusSquared = equatorRadius * equatorRadius
	polesRadiusSquared   = polesRadius * polesRadius
)

var excentricitySquared float64 = math.Pow(math.Sqrt((equatorRadiusSquared-polesRadiusSquared)/equatorRadiusSquared), 2)

func ComputeDistanceHubeny(a, b Coordinates) float64 {
	coordinate1Latitude := degreesToRadian(a.Latitude)
	coordinate2Latitude := degreesToRadian(b.Latitude)

	coordinate1Longitude := degreesToRadian(a.Longitude)
	coordinate2Longitude := degreesToRadian(b.Longitude)

	deltaPhi := coordinate2Latitude - coordinate1Latitude
	deltaLambda := coordinate2Longitude - coordinate1Longitude

	phi := (coordinate1Latitude + coordinate2Latitude) / 2.0

	W := 1.0 - excentricitySquared*math.Pow(math.Sin(phi), 2.0)

	M := equatorRadius * (1.0 - excentricitySquared) / math.Sqrt(math.Pow(W, 3.0))
	N := equatorRadius / math.Sqrt(W)

	distance := math.Sqrt(math.Pow(M*deltaPhi, 2.0) + math.Pow(N*math.Cos(phi)*deltaLambda, 2.0))

	return distance / 1000.0 // Return distance in kilometers.
}

func degreesToRadian(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}
