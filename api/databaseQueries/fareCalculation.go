package databaseQueries

import "strconv"

// Declare initial variables to be used during function call
var express bool
var XpressRoutes = []string{"27x", "33d", "33x", "39x", "41x",
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

// CalculateFare is a function designed to populate the Fares property
// of a busRouteJSON object that is returned following an api call to
// match a route to a given pair of bus stops. It takes as parameters the
// busRoute that is being used to calculate the appropriate fares, the
// origin bus stop number as a string and the destination bus stop
// numbers, also as a string. It returns a busFares object containing the
// appropriate fares for each demographic.
func CalculateFare(route busRoute,
	originStop string,
	destinationStop string) busFares {

	// Declare variables to be maintained locally during each call to
	// this function
	var originDist float64
	var destDist float64
	var calculatedFares busFares

	// Boolean condition defaults to false unless determined otherwise
	express = false
	for _, routeNum := range XpressRoutes {
		if route.Id[0] == routeNum {
			express = true
		}
	}

	// If route is express route, immediate assignment and return is possible
	if express {
		calculatedFares.AdultLeap = XpressoAdultLeap
		calculatedFares.AdultCash = XpressoAdultCash
		calculatedFares.StudentLeap = XpressoStudentLeap
		calculatedFares.ChildLeap = XpressoChildLeap
		calculatedFares.ChildCash = XpressoChildCash
		return calculatedFares
	} else {

		// If route was not express, the distance between two stops must be
		// calculated and compared against the short zone limiter

		for _, stopCounter := range route.Stops {
			if stopCounter.StopNumber == originStop {
				originDist, _ = strconv.ParseFloat(stopCounter.DistanceTravelled, 64)
			} else if stopCounter.StopNumber == destinationStop {
				destDist, _ = strconv.ParseFloat(stopCounter.DistanceTravelled, 64)
			}
		}

		totalDist := destDist - originDist

		// Use comparison of distance travelled against the short zone
		// limit to determine the appropriate fares and return

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
