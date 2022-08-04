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

// Global variables

// Variables of both busRoute and busRouteJSON need to be initialised as
// some unmarshalling from Mongo cannot be done automatically and
// so must be done manually from one structure to another in the backend
var result []busRoute
var resultJSON []busRouteJSON
var route busRouteJSON
var stop RouteStop
var shape ShapeJSON
var stops []RouteStop
var shapes []ShapeJSON
var originStopArrivalTime string
var destinationStopArrivalTime string
var finalStopArrivalTime string
var firstStopArrivalTime string
var originStopSequence int64
var destinationStopSequence int64
var originDistTravelled float64
var destinationDistTravelled float64

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

func GetTimeString(date string) string {

	dateStringSplit := strings.Split(date, " ")
	timeString := dateStringSplit[1]

	return timeString
}

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
