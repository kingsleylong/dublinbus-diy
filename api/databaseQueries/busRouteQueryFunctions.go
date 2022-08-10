package databaseQueries

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
	"strings"
)

// ConnectToMongo is a function used specifically to handle the database connection
// within the backend and remove some boilerplate code from individual function calls
// elsewhere. It requires no parameters but returns a pointer to a Mongo client as well
// as an error
func ConnectToMongo() (*mongo.Client, error) {

	mongoHost = os.Getenv("MONGO_INITDB_ROOT_HOST")
	mongoPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongoUsername = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPort = os.Getenv("MONGO_INITDB_ROOT_PORT")

	// Create connection to mongo server and log any resulting error
	client, err := mongo.NewClient(options.Client().
		ApplyURI(
			fmt.Sprintf(
				"mongodb://%s:%s@%s:%s/?retryWrites=true&w=majority",
				mongoUsername,
				mongoPassword,
				mongoHost,
				mongoPort)))
	if err != nil {
		log.Print(err)
	}

	return client, err
}

// GetTimeString is used to take a string that has a date and time represented in it
// and then return in string format the time part specifically of that initial parameter.
// This function takes a string with the date in the format "yyyy-mm-dd hh:mm:ss", with the
// whitespace between calendar representation and time representation an important and
// necessary element of this parameter. This function returns a string of the time that
// was inputted into the function in the format "hh:mm:ss"
func GetTimeString(date string) string {

	dateStringSplit := strings.Split(date, " ")
	timeString := dateStringSplit[1]

	return timeString
}

// CreateStopsSlice takes in the origin and destination bus stop numbers along
// a route as strings, as well as the busRoute object that these stops belong to
// and a RouteStop object (with empty fields or otherwise) and returns a slice
// of RouteStop objects. This function is designed to handle the conversion of
// the stops as taken from the document in the Mongo collection and transform them
// into the format necessary for the return value of the function for the api call.
// It also assigns values to other global variables that are necessary later for other
// functions in the parent route matching function.
func CreateStopsSlice(origin string, destination string,
	route busRoute, stop RouteStop) []RouteStop {

	transformedStops := []RouteStop{}

	for _, initialStopDescription := range route.Stops {
		stop.StopId = initialStopDescription.StopId
		stop.StopName = initialStopDescription.StopName
		stop.StopNumber = initialStopDescription.StopNumber
		stop.StopLat, _ = strconv.ParseFloat(initialStopDescription.StopLat, 64)
		stop.StopLon, _ = strconv.ParseFloat(initialStopDescription.StopLon, 64)
		stop.StopSequence = initialStopDescription.StopSequence
		stop.ArrivalTime = initialStopDescription.ArrivalTime
		stop.DepartureTime = initialStopDescription.DepartureTime
		stop.DistanceTravelled, _ =
			strconv.ParseFloat(initialStopDescription.DistanceTravelled, 64)
		if initialStopDescription.StopSequence == "1" {
			firstStopArrivalTime = initialStopDescription.ArrivalTime
		}
		if initialStopDescription.StopNumber == origin {
			originStopSequence, _ = strconv.ParseInt(initialStopDescription.StopSequence, 10, 64)
			originStopArrivalTime = initialStopDescription.ArrivalTime
			originDistTravelled, _ = strconv.ParseFloat(initialStopDescription.DistanceTravelled, 64)
		}
		if initialStopDescription.StopNumber == destination {
			destinationStopSequence, _ = strconv.ParseInt(initialStopDescription.StopSequence, 10, 64)
			destinationStopArrivalTime = initialStopDescription.ArrivalTime
			destinationDistTravelled, _ = strconv.ParseFloat(initialStopDescription.DistanceTravelled, 64)
		}
		finalStopArrivalTime = initialStopDescription.ArrivalTime
		transformedStops = append(transformedStops, stop)
	}

	return transformedStops
}

// CreateShapesSlice is a function that takes in a busRoute object and then returns
// a slice of ShapeJSON objects that are then used for the final creation of the
// busRouteJSON objects that are returned to the frontend following a successful route
// finding operation.
func CreateShapesSlice(route busRoute) []ShapeJSON {

	shapes = []ShapeJSON{}
	for _, currentShape := range route.Shapes {
		currentDistTravelled, _ := strconv.ParseFloat(currentShape.ShapeDistTravel, 64)
		if currentDistTravelled >= originDistTravelled && currentDistTravelled <= destinationDistTravelled {
			shape.ShapePtLat, _ = strconv.ParseFloat(currentShape.ShapePtLat, 64)
			shape.ShapePtLon, _ = strconv.ParseFloat(currentShape.ShapePtLon, 64)
			shape.ShapePtSequence = currentShape.ShapePtSequence
			shape.ShapeDistTravel = currentShape.ShapeDistTravel
			shapes = append(shapes, shape)
		}
	}

	return shapes
}

// CurateStopsSlice is a function that takes in the origin and destination
// bus stop numbers on a journey as strings and then returns integers for their
// respective indexes in the route object that is to be added to the resultJSON object
// at the end of a route finding operation
func CurateStopsSlice(origin string, destination string) (int, int) {

	var originStopIndex int
	var destinationStopIndex int

	for index, stopToAdd := range route.Stops {
		if stopToAdd.StopNumber == origin {
			originStopIndex = index
		}
		if stopToAdd.StopNumber == destination {
			destinationStopIndex = index
		}
	}

	return originStopIndex, destinationStopIndex
}

// CurateReturnedArrivalRoutes takes in as a string the original time that was queried
// for the arrival time based query and the slice of busRouteJSON objects about to be
// returned to the front end at the end of route finding sequence. This function selects
// for routes that have an arrival time within one hour before of the specified time only
// to limit the number of routes being returned and potentially improve the UX of the application
func CurateReturnedArrivalRoutes(arrivalQueryTime string, routes []busRouteJSON) []busRouteJSON {

	returnedRoutes := []busRouteJSON{}
	dateTimeSplit := strings.Split(arrivalQueryTime, " ")
	querySeconds := convertStringTimeToTotalSeconds(dateTimeSplit[1])
	var arrivalSeconds float64
	var stopsLength int

	for _, route := range routes {
		stopsLength = len(route.Stops)
		arrivalSeconds = convertStringTimeToTotalSeconds(route.Stops[stopsLength-1].ArrivalTime)
		if querySeconds > arrivalSeconds+float64(60*60) {
			continue
		}
		returnedRoutes = append(returnedRoutes, route)
	}

	return returnedRoutes
}

// GetTimeStringAsHoursAndMinutes is a function designed to take in a string representing
// time of day in the format "hh:mm:ss" and return a string approximating this time by
// removing the seconds component and just displaying "hh:mm"
func GetTimeStringAsHoursAndMinutes(timeString string) string {

	timeSplit := strings.Split(timeString, ":")
	timeAdjusted := timeSplit[0] + ":" + timeSplit[1]

	return timeAdjusted
}
