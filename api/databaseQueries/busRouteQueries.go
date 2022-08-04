package databaseQueries

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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

// FindMatchingRouteForDeparture takes in three parameters - the destination
// bus stop, the origin bus stop and then the departure time all as strings.
// This function then queries the mongo collection for trips documents that
// match these filters before mapping the documents to the correct structure
// and returning them within a slice of type busRouteJSON.
func FindMatchingRouteForDeparture(destination string,
	origin string,
	date string) []busRouteJSON {

	client, err := ConnectToMongo()

	// Create context variable and assign time for timeout
	// Log any resulting error here also
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Print(err)
	}
	defer client.Disconnect(ctx) // defer has rest of function complete before disconnect

	timeString := GetTimeString(date)

	// Aggregation pipeline created in Mongo Compass and then transformed to suit
	// the mongo driver in Go
	coll := client.Database("BusData").Collection("trips_n_stops")
	cursor, err := coll.Aggregate(ctx, bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"stops.stop_number", destination},
					{"stops",
						bson.D{
							{"$elemMatch",
								bson.D{
									{"stop_number", origin},
									{"departure_time",
										bson.D{{"$gt", timeString}}},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$sort",
				bson.D{
					{"route.route_short_name", 1},
					{"stops.departure_time", 1},
					{"stops.stop_sequence", 1},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$route.route_short_name"},
					{"stops", bson.D{{"$first", "$stops"}}},
					{"shapes", bson.D{{"$first", "$shapes"}}},
				},
			},
		},
	})
	if err != nil {
		log.Print(err)
	}

	if err = cursor.All(ctx, &result); err != nil {
		log.Print(err)
	}

	// Loop through the stops that are in the result slice and start manually
	// converting them to the RouteStop type to be added to a busRouteJSON
	// object that is part of the returned slice. This is necessary as some
	// data types need to be changed and this has to be done manually
	for _, currentRoute := range result {
		route.RouteNum = currentRoute.Id

		route.Stops = CreateStopsSlice(origin, destination, currentRoute, stop)

		// An empty slice of shapes is created here for each outer iteration for
		// the same reason as the empty slice for the stops above
		shapes = []ShapeJSON{}
		for _, currentShape := range currentRoute.Shapes {
			currentDistTravelled, _ := strconv.ParseFloat(currentShape.ShapeDistTravel, 64)
			if currentDistTravelled >= originDistTravelled && currentDistTravelled <= destinationDistTravelled {
				shape.ShapePtLat, _ = strconv.ParseFloat(currentShape.ShapePtLat, 64)
				shape.ShapePtLon, _ = strconv.ParseFloat(currentShape.ShapePtLon, 64)
				shape.ShapePtSequence = currentShape.ShapePtSequence
				shape.ShapeDistTravel = currentShape.ShapeDistTravel
				shapes = append(shapes, shape)
			}
		}
		route.Shapes = shapes

		if originStopSequence > destinationStopSequence {
			continue
		}

		// Use the CalculateFare function from fareCalculation.go to get the fares
		// object for each route
		route.Fares = CalculateFare(currentRoute, origin, destination)

		if currentRoute.Direction == "1" {
			route.Direction = "2"
		} else {
			route.Direction = "1"
		}

		initialTravelTime, err := GetTravelTimePrediction(route.RouteNum, date, route.Direction)
		if err != nil {
			log.Println(err)
		}

		journeyTravelTime := AdjustTravelTime(initialTravelTime, originStopArrivalTime,
			destinationStopArrivalTime, firstStopArrivalTime, finalStopArrivalTime)
		if journeyTravelTime.Source == "static" {
			staticTravelTime := GetStaticTime(originStopArrivalTime, destinationStopArrivalTime)
			journeyTravelTime.TransitTime = staticTravelTime
			journeyTravelTime.TransitTimeMinusMAE = staticTravelTime
			journeyTravelTime.TransitTimePlusMAE = staticTravelTime
		}
		route.TravelTime = journeyTravelTime

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

		route.Stops = route.Stops[originStopIndex : destinationStopIndex+1]

		resultJSON = append(resultJSON, route)
	}

	return resultJSON
}

// FindMatchingRouteForArrival takes in three parameters - the origin
// bus stop, the destination bus stop and then the arrival time all as strings.
// This function then queries the mongo collection for trips documents that
// match these filters before mapping the documents to the correct structure
// and returning them within a slice of type busRouteJSON.
func FindMatchingRouteForArrival(origin string,
	destination string,
	date string) []busRouteJSON {

	client, err := ConnectToMongo()

	// Create context variable and assign time for timeout
	// Log any resulting error here also
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Print(err)
	}
	defer client.Disconnect(ctx) // defer has rest of function complete before disconnect

	timeString := GetTimeString(date)

	// Aggregation pipeline created in Mongo Compass and then transformed to suit
	// the mongo driver in Go
	coll := client.Database("BusData").Collection("trips_n_stops")
	cursor, err := coll.Aggregate(ctx, bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"stops.stop_number", origin},
					{"stops",
						bson.D{
							{"$elemMatch",
								bson.D{
									{"stop_number", destination},
									{"arrival_time",
										bson.D{{"$lte", timeString}}},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$sort",
				bson.D{
					{"route.route_short_name", 1},
					{"stops.arrival_time", -1},
					{"stops.stop_sequence", 1},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$route.route_short_name"},
					{"stops", bson.D{{"$first", "$stops"}}},
					{"shapes", bson.D{{"$first", "$shapes"}}},
				},
			},
		},
	})
	if err != nil {
		log.Print(err)
	}

	if err = cursor.All(ctx, &result); err != nil {
		log.Print(err)
	}

	// Loop through the stops that are in the result slice and start manually
	// converting them to the RouteStop type to be added to a busRouteJSON
	// object that is part of the returned slice. This is necessary as some
	// data types need to be changed and this has to be done manually
	for _, currentRoute := range result {

		route.RouteNum = currentRoute.Id

		// An empty slice of stops is created with each new outer iteration so
		// that duplicates aren't added to later routes in their stop arrays
		route.Stops = CreateStopsSlice(origin, destination, currentRoute, stop)

		// An empty slice of shapes is created here for each outer iteration for
		// the same reason as the empty slice for the stops above
		shapes = []ShapeJSON{}
		for _, currentShape := range currentRoute.Shapes {
			currentDistTravelled, _ := strconv.ParseFloat(currentShape.ShapeDistTravel, 64)
			if currentDistTravelled >= originDistTravelled && currentDistTravelled <= destinationDistTravelled {
				shape.ShapePtLat, _ = strconv.ParseFloat(currentShape.ShapePtLat, 64)
				shape.ShapePtLon, _ = strconv.ParseFloat(currentShape.ShapePtLon, 64)
				shape.ShapePtSequence = currentShape.ShapePtSequence
				shape.ShapeDistTravel = currentShape.ShapeDistTravel
				shapes = append(shapes, shape)
			}
		}
		route.Shapes = shapes

		if originStopSequence > destinationStopSequence {
			continue
		}

		// Use the CalculateFare function from fareCalculation.go to get the fares
		// object for each route
		route.Fares = CalculateFare(currentRoute, origin, destination)

		if currentRoute.Direction == "1" {
			route.Direction = "2"
		} else {
			route.Direction = "1"
		}

		initialTravelTime, err := GetTravelTimePrediction(route.RouteNum, date, route.Direction)
		if err != nil {
			log.Println(err)
		}

		journeyTravelTime := AdjustTravelTime(initialTravelTime, originStopArrivalTime,
			destinationStopArrivalTime, firstStopArrivalTime, finalStopArrivalTime)

		if journeyTravelTime.Source == "static" {
			staticTravelTime := GetStaticTime(originStopArrivalTime, destinationStopArrivalTime)
			journeyTravelTime.TransitTime = staticTravelTime
			journeyTravelTime.TransitTimeMinusMAE = staticTravelTime
			journeyTravelTime.TransitTimePlusMAE = staticTravelTime
		}
		route.TravelTime = journeyTravelTime

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

		route.Stops = route.Stops[originStopIndex : destinationStopIndex+1]

		resultJSON = append(resultJSON, route)
	}

	return resultJSON
}

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
