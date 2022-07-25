package databaseQueries

import "strconv"

var express bool
var XpressRoutes = []string{"46a", "27x", "33d", "33x", "39x", "41x",
	"51x", "51d", "51x", "69x", "77x", "84x"}

const (
	XpressoAdultCash   = float64(3.00)
	ShortZoneAdultCash = float64(1.70)
	LongZoneAdultCash  = float64(2.60)
	XpressoAdultLeap   = float64(2.40)
	ShortZoneAdultLeap = float64(1.30)
	LongZoneAdultLeap  = float64(2.00)

	XpressoStudentLeap   = float64(1.20)
	ShortZoneStudentLeap = float64(0.65)
	LongZoneStudentLeap  = float64(1.00)

	XpressoChildLeap  = float64(1.00)
	LongZoneChildLeap = float64(0.65)
	XpressoChildCash  = float64(1.30)
	LongZoneChildCash = float64(0.90)
	ShortZoneDistance = float64(3000.0)
)

func CalculateFare(route busRoute,
	originStop string,
	destinationStop string) busFares {

	var originDist float64
	var destDist float64
	var calculatedFares busFares

	express = false
	for _, routeNum := range XpressRoutes {
		if route.Id == routeNum {
			express = true
		}
	}

	if express {
		calculatedFares.AdultLeap = XpressoAdultLeap
		calculatedFares.AdultCash = XpressoAdultCash
		calculatedFares.StudentLeap = XpressoStudentLeap
		calculatedFares.ChildLeap = XpressoChildLeap
		calculatedFares.ChildCash = XpressoChildCash
		return calculatedFares
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
			calculatedFares.AdultLeap = ShortZoneAdultLeap
			calculatedFares.AdultCash = ShortZoneAdultCash
			calculatedFares.StudentLeap = ShortZoneStudentLeap
			calculatedFares.ChildLeap = LongZoneChildLeap
			calculatedFares.ChildCash = LongZoneChildCash
			return calculatedFares
		} else {
			calculatedFares.AdultLeap = LongZoneAdultLeap
			calculatedFares.AdultCash = LongZoneAdultCash
			calculatedFares.StudentLeap = LongZoneStudentLeap
			calculatedFares.ChildLeap = LongZoneChildLeap
			calculatedFares.ChildCash = LongZoneChildCash
			return calculatedFares
		}
	}
}
