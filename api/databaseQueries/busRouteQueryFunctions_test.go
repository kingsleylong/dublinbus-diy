package databaseQueries

import (
	"reflect"
	"testing"
)

func TestGetTimeString(t *testing.T) {
	timeString := GetTimeString("2022-08-12 00:30:00")
	if timeString != "00:30:00" {
		t.Log("error: timeString should be '00:30:00' but is", timeString)
		t.Fail()
	}
}

func TestCreateStopsSlice(t *testing.T) {

	testOrigin := "1"
	testDestination := "3"
	testId := []byte{}
	testDirection := "1"
	var testStops []BusStop
	var testShapes []Shape
	var testRouteStop RouteStop
	var testRoute busRoute

	stopOne := BusStop{
		StopId:            "random",
		StopNumber:        "1",
		StopName:          "First",
		StopLat:           "53.1",
		StopLon:           "-6.0",
		StopSequence:      "1",
		ArrivalTime:       "07:00:00",
		DepartureTime:     "07:01:00",
		DistanceTravelled: "1000",
	}

	stopTwo := BusStop{
		StopId:            "random2",
		StopNumber:        "2",
		StopName:          "Second",
		StopLat:           "53.2",
		StopLon:           "-6.1",
		StopSequence:      "2",
		ArrivalTime:       "07:11:00",
		DepartureTime:     "07:12:00",
		DistanceTravelled: "2000",
	}

	stopThree := BusStop{
		StopId:            "random3",
		StopNumber:        "3",
		StopName:          "Third",
		StopLat:           "53.3",
		StopLon:           "-6.2",
		StopSequence:      "3",
		ArrivalTime:       "07:22:00",
		DepartureTime:     "07:23:00",
		DistanceTravelled: "3000",
	}

	testStops = append(testStops, stopOne, stopTwo, stopThree)

	shapeOne := Shape{
		ShapePtLat:      "53.1",
		ShapePtLon:      "-6.0",
		ShapePtSequence: "1",
		ShapeDistTravel: "1000",
	}

	shapeTwo := Shape{
		ShapePtLat:      "53.2",
		ShapePtLon:      "-6.1",
		ShapePtSequence: "2",
		ShapeDistTravel: "2000",
	}

	shapeThree := Shape{
		ShapePtLat:      "53.3",
		ShapePtLon:      "-6.2",
		ShapePtSequence: "3",
		ShapeDistTravel: "3000",
	}

	testShapes = append(testShapes, shapeOne, shapeTwo, shapeThree)

	testRoute.Id = testId
	testRoute.Stops = testStops
	testRoute.Shapes = testShapes
	testRoute.Direction = testDirection

	stopsSlice := CreateStopsSlice(testOrigin, testDestination, testRoute, testRouteStop)

	for _, stop := range stopsSlice {
		if reflect.TypeOf(stop.StopLat).Name() != "float64" ||
			reflect.TypeOf(stop.StopLon).Name() != "float64" ||
			reflect.TypeOf(stop.DistanceTravelled).Name() != "float64" {
			t.Log("Data type conversion unsuccessful")
			t.Fail()
		}
	}
	if firstStopArrivalTime != "07:00:00" {
		t.Log("First stop arrival time should be '07:00:00, but is", firstStopArrivalTime)
		t.Fail()
	}
	if originStopSequence != 1 {
		t.Log("Error with origin stop sequence, which should be int of value 1")
		t.Fail()
	}
	if originStopArrivalTime != "07:00:00" {
		t.Log("Origin stop arrival time should be 07:00:00, but is", originStopArrivalTime)
		t.Fail()
	}
	if originDistTravelled != float64(1000) {
		t.Log("Origin distance travelled assignment was incorrect")
		t.Fail()
	}
	if finalStopArrivalTime != "07:22:00" {
		t.Log("Final stop arrival time should be '07:22:00, but is", finalStopArrivalTime)
		t.Fail()
	}
	if destinationStopSequence != 3 {
		t.Log("Error with destination stop sequence, which should be int of value 3")
		t.Fail()
	}
	if destinationStopArrivalTime != "07:22:00" {
		t.Log("Destination stop arrival time should be 07:22:00, but is", destinationStopArrivalTime)
		t.Fail()
	}
	if destinationDistTravelled != float64(3000) {
		t.Log("Destination distance travelled assignment was incorrect")
		t.Fail()
	}
}

func TestCreateShapesSlice(t *testing.T) {

	testId := []byte{}
	testDirection := "1"
	var testStops []BusStop
	var testShapes []Shape
	var testRoute busRoute

	stopOne := BusStop{
		StopId:            "random",
		StopNumber:        "1",
		StopName:          "First",
		StopLat:           "53.1",
		StopLon:           "-6.0",
		StopSequence:      "1",
		ArrivalTime:       "07:00:00",
		DepartureTime:     "07:01:00",
		DistanceTravelled: "1000",
	}

	stopTwo := BusStop{
		StopId:            "random2",
		StopNumber:        "2",
		StopName:          "Second",
		StopLat:           "53.2",
		StopLon:           "-6.1",
		StopSequence:      "2",
		ArrivalTime:       "07:11:00",
		DepartureTime:     "07:12:00",
		DistanceTravelled: "2000",
	}

	stopThree := BusStop{
		StopId:            "random3",
		StopNumber:        "3",
		StopName:          "Third",
		StopLat:           "53.3",
		StopLon:           "-6.2",
		StopSequence:      "3",
		ArrivalTime:       "07:22:00",
		DepartureTime:     "07:23:00",
		DistanceTravelled: "3000",
	}

	testStops = append(testStops, stopOne, stopTwo, stopThree)

	shapeOne := Shape{
		ShapePtLat:      "53.1",
		ShapePtLon:      "-6.0",
		ShapePtSequence: "1",
		ShapeDistTravel: "1000",
	}

	shapeTwo := Shape{
		ShapePtLat:      "53.2",
		ShapePtLon:      "-6.1",
		ShapePtSequence: "2",
		ShapeDistTravel: "2000",
	}

	shapeThree := Shape{
		ShapePtLat:      "53.3",
		ShapePtLon:      "-6.2",
		ShapePtSequence: "3",
		ShapeDistTravel: "3000",
	}

	testShapes = append(testShapes, shapeOne, shapeTwo, shapeThree)

	testRoute.Id = testId
	testRoute.Stops = testStops
	testRoute.Shapes = testShapes
	testRoute.Direction = testDirection

	testShapesJSON := CreateShapesSlice(testRoute)

	for _, shape := range testShapesJSON {

		if reflect.TypeOf(shape.ShapePtLon).Name() != "float64" ||
			reflect.TypeOf(shape.ShapePtLat).Name() != "float64" {
			t.Log("Coordinates were not converted to floating point numbers")
			t.Fail()
		}
	}
}

func TestCurateStopsSlice(t *testing.T) {

	testOrigin := "2"
	testDestination := "4"

	testId := "route1"
	testDirection := "1"
	var testStops []RouteStop
	var testShapes []ShapeJSON
	var testRoute busRouteJSON
	var testTravelTime TravelTimePrediction
	var testFares busFares

	stopOne := RouteStop{
		StopId:            "random",
		StopNumber:        "1",
		StopName:          "First",
		StopLat:           53.1,
		StopLon:           -6.0,
		StopSequence:      "1",
		ArrivalTime:       "07:00:00",
		DepartureTime:     "07:01:00",
		DistanceTravelled: 1000,
	}

	stopTwo := RouteStop{
		StopId:            "random2",
		StopNumber:        "2",
		StopName:          "Second",
		StopLat:           53.2,
		StopLon:           -6.1,
		StopSequence:      "2",
		ArrivalTime:       "07:11:00",
		DepartureTime:     "07:12:00",
		DistanceTravelled: 2000,
	}

	stopThree := RouteStop{
		StopId:            "random3",
		StopNumber:        "3",
		StopName:          "Third",
		StopLat:           53.3,
		StopLon:           -6.2,
		StopSequence:      "3",
		ArrivalTime:       "07:22:00",
		DepartureTime:     "07:23:00",
		DistanceTravelled: 3000,
	}

	stopFour := RouteStop{
		StopId:            "random4",
		StopNumber:        "4",
		StopName:          "Fourth",
		StopLat:           53.4,
		StopLon:           -6.3,
		StopSequence:      "4",
		ArrivalTime:       "07:33:00",
		DepartureTime:     "07:34:00",
		DistanceTravelled: 4000,
	}

	stopFive := RouteStop{
		StopId:            "random5",
		StopNumber:        "5",
		StopName:          "Fifth",
		StopLat:           53.4,
		StopLon:           -6.3,
		StopSequence:      "5",
		ArrivalTime:       "07:44:00",
		DepartureTime:     "07:45:00",
		DistanceTravelled: 5000,
	}

	testStops = append(testStops, stopOne, stopTwo, stopThree, stopFour, stopFive)

	shapeOne := ShapeJSON{
		ShapePtLat:      53.1,
		ShapePtLon:      -6.0,
		ShapePtSequence: "1",
		ShapeDistTravel: "1000",
	}

	shapeTwo := ShapeJSON{
		ShapePtLat:      53.2,
		ShapePtLon:      -6.1,
		ShapePtSequence: "2",
		ShapeDistTravel: "2000",
	}

	shapeThree := ShapeJSON{
		ShapePtLat:      53.3,
		ShapePtLon:      -6.2,
		ShapePtSequence: "3",
		ShapeDistTravel: "3000",
	}

	testShapes = append(testShapes, shapeOne, shapeTwo, shapeThree)

	testRoute.RouteNum = testId
	testRoute.Stops = testStops
	testRoute.Shapes = testShapes
	testRoute.Direction = testDirection
	testRoute.TravelTime = testTravelTime
	testRoute.Fares = testFares

	originNum, destinationNum := CurateStopsSlice(testOrigin, testDestination, testRoute)

	if originNum != 1 {
		t.Log("Origin should have been int of value 2 but instead was", originNum)
		t.Fail()
	}
	if destinationNum != 3 {
		t.Log("Destination should have been int of value 4 but instead was", destinationNum)
	}
}

func TestCurateReturnedArrivalRoutes(t *testing.T) {

	testArrivalTime := "2022-08-12 10:00:00"
	var testBusRoutes []busRouteJSON

	var testStopsOne []RouteStop
	var testStopsTwo []RouteStop
	var testShapes []ShapeJSON
	var testRouteOne busRouteJSON
	var testRouteTwo busRouteJSON
	var testTravelTime TravelTimePrediction
	var testFares busFares

	stopOne := RouteStop{
		StopId:            "random",
		StopNumber:        "1",
		StopName:          "First",
		StopLat:           53.1,
		StopLon:           -6.0,
		StopSequence:      "1",
		ArrivalTime:       "07:00:00",
		DepartureTime:     "07:01:00",
		DistanceTravelled: 1000,
	}

	stopTwo := RouteStop{
		StopId:            "random2",
		StopNumber:        "2",
		StopName:          "Second",
		StopLat:           53.2,
		StopLon:           -6.1,
		StopSequence:      "2",
		ArrivalTime:       "07:11:00",
		DepartureTime:     "07:12:00",
		DistanceTravelled: 2000,
	}

	stopThree := RouteStop{
		StopId:            "random3",
		StopNumber:        "3",
		StopName:          "Third",
		StopLat:           53.3,
		StopLon:           -6.2,
		StopSequence:      "3",
		ArrivalTime:       "07:22:00",
		DepartureTime:     "07:23:00",
		DistanceTravelled: 3000,
	}

	stopFour := RouteStop{
		StopId:            "random4",
		StopNumber:        "4",
		StopName:          "Fourth",
		StopLat:           53.4,
		StopLon:           -6.3,
		StopSequence:      "4",
		ArrivalTime:       "07:33:00",
		DepartureTime:     "07:34:00",
		DistanceTravelled: 4000,
	}

	stopFive := RouteStop{
		StopId:            "random5",
		StopNumber:        "5",
		StopName:          "Fifth",
		StopLat:           53.4,
		StopLon:           -6.3,
		StopSequence:      "5",
		ArrivalTime:       "07:44:00",
		DepartureTime:     "07:45:00",
		DistanceTravelled: 5000,
	}

	testStopsOne = append(testStopsOne, stopOne, stopTwo, stopThree, stopFour, stopFive)

	shapeOne := ShapeJSON{
		ShapePtLat:      53.1,
		ShapePtLon:      -6.0,
		ShapePtSequence: "1",
		ShapeDistTravel: "1000",
	}

	shapeTwo := ShapeJSON{
		ShapePtLat:      53.2,
		ShapePtLon:      -6.1,
		ShapePtSequence: "2",
		ShapeDistTravel: "2000",
	}

	shapeThree := ShapeJSON{
		ShapePtLat:      53.3,
		ShapePtLon:      -6.2,
		ShapePtSequence: "3",
		ShapeDistTravel: "3000",
	}

	testShapes = append(testShapes, shapeOne, shapeTwo, shapeThree)

	testRouteOne.RouteNum = "route1"
	testRouteOne.Stops = testStopsOne
	testRouteOne.Shapes = testShapes
	testRouteOne.Direction = "1"
	testRouteOne.TravelTime = testTravelTime
	testRouteOne.Fares = testFares

	stopSix := RouteStop{
		StopId:            "random",
		StopNumber:        "1",
		StopName:          "First",
		StopLat:           53.1,
		StopLon:           -6.0,
		StopSequence:      "1",
		ArrivalTime:       "09:00:00",
		DepartureTime:     "09:01:00",
		DistanceTravelled: 1000,
	}

	stopSeven := RouteStop{
		StopId:            "random7",
		StopNumber:        "2",
		StopName:          "Second",
		StopLat:           53.2,
		StopLon:           -6.1,
		StopSequence:      "2",
		ArrivalTime:       "09:11:00",
		DepartureTime:     "09:12:00",
		DistanceTravelled: 2000,
	}

	stopEight := RouteStop{
		StopId:            "random3",
		StopNumber:        "3",
		StopName:          "Third",
		StopLat:           53.3,
		StopLon:           -6.2,
		StopSequence:      "3",
		ArrivalTime:       "09:22:00",
		DepartureTime:     "09:23:00",
		DistanceTravelled: 3000,
	}

	stopNine := RouteStop{
		StopId:            "random4",
		StopNumber:        "4",
		StopName:          "Fourth",
		StopLat:           53.4,
		StopLon:           -6.3,
		StopSequence:      "4",
		ArrivalTime:       "09:33:00",
		DepartureTime:     "09:34:00",
		DistanceTravelled: 4000,
	}

	stopTen := RouteStop{
		StopId:            "random5",
		StopNumber:        "5",
		StopName:          "Fifth",
		StopLat:           53.4,
		StopLon:           -6.3,
		StopSequence:      "5",
		ArrivalTime:       "09:44:00",
		DepartureTime:     "09:45:00",
		DistanceTravelled: 5000,
	}

	testStopsTwo = append(testStopsTwo, stopSix, stopSeven, stopEight, stopNine, stopTen)

	testRouteTwo.RouteNum = "route2"
	testRouteTwo.Stops = testStopsTwo
	testRouteTwo.Shapes = testShapes
	testRouteTwo.Direction = "1"
	testRouteTwo.TravelTime = testTravelTime
	testRouteTwo.Fares = testFares

	testBusRoutes = append(testBusRoutes, testRouteOne, testRouteTwo)
	testCuratedArrivalRoutes := CurateReturnedArrivalRoutes(testArrivalTime, testBusRoutes)

	if len(testCuratedArrivalRoutes) != 1 {
		t.Log("Only one route should have been returned as valid")
		t.Fail()
	}
	if testCuratedArrivalRoutes[0].RouteNum != "route2" {
		t.Log("Route with id 'route2' should have been returned but was not")
		t.Fail()
	}
}
