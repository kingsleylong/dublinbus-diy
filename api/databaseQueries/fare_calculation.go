package databaseQueries

import "strconv"

var express bool
var XpressRoutes = []string{"46a", "27x", "33d", "33x", "39x", "41x",
	"51x", "51d", "51x", "69x", "77x", "84x"}

const (
	XpressoAdult      = float64(3.00)
	ShortZoneAdult    = float64(1.7)
	LongZoneAdult     = float64(2.60)
	ShortZoneDistance = float64(3000.0)
)

func CalculateFare(route busRoute,
	originStop string,
	destinationStop string) float64 {

	var originDist float64
	var destDist float64

	express = false
	for _, routeNum := range XpressRoutes {
		if route.Id == routeNum {
			express = true
		}
	}

	if express {
		return XpressoAdult
	} else {
		for _, stopCounter := range route.Stops {
			if stopCounter.StopNumber == originStop {
				originDist, _ = strconv.ParseFloat(stopCounter.DistanceTravelled, 64)
			} else if stopCounter.StopNumber == destinationStop {
				destDist, _ = strconv.ParseFloat(stopCounter.DistanceTravelled, 64)
			}
		}

		totalDist := destDist - originDist

		if totalDist < ShortZoneDistance {
			return ShortZoneAdult
		} else {
			return LongZoneAdult
		}
	}
}
