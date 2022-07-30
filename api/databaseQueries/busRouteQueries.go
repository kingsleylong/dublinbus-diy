package databaseQueries

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// FindMatchingRouteForDeparture takes in three parameters - the destination
// bus stop, the origin bus stop and then the departure time all as strings.
// This function then queries the mongo collection for trips documents that
// match these filters before mapping the documents to the correct structure
// and returning them within a slice of type busRouteJSON.
func FindMatchingRouteForDeparture(destination string,
	origin string,
	date string) []busRouteJSON {

	// Assign values to connection string variables
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

	// Create context variable and assign time for timeout
	// Log any resulting error here also
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Print(err)
	}
	defer client.Disconnect(ctx) // defer has rest of function done before disconnect

	dateStringSplit := strings.Split(date, " ")
	//dateString := dateStringSplit[0]
	timeString := dateStringSplit[1]

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
		stops = []RouteStop{}
		for _, currentStop := range currentRoute.Stops {
			stop.StopId = currentStop.StopId
			stop.StopName = currentStop.StopName
			stop.StopNumber = currentStop.StopNumber
			stop.StopLat, _ = strconv.ParseFloat(currentStop.StopLat, 64)
			stop.StopLon, _ = strconv.ParseFloat(currentStop.StopLon, 64)
			stop.StopSequence = currentStop.StopSequence
			stop.ArrivalTime = currentStop.ArrivalTime
			stop.DepartureTime = currentStop.DepartureTime
			stop.DistanceTravelled, _ =
				strconv.ParseFloat(currentStop.DistanceTravelled, 64)
			if currentStop.StopNumber == origin {
				originStopArrivalTime = currentStop.ArrivalTime
				log.Println("Found the origin stop: " + currentStop.StopNumber)
				log.Println("Origin stop was in the sequence at stop " + currentStop.StopSequence)
				log.Println("Bus Arrived here at " + originStopArrivalTime)
				log.Println("- - - - - - - - - - - - - - - - - - - - - - -")
			}
			if currentStop.StopNumber == destination {
				destinationStopArrivalTime = currentStop.ArrivalTime
				log.Println("Found the destination stop: " + currentStop.StopNumber)
				log.Println("Destination stop was in the sequence at stop " + currentStop.StopSequence)
				log.Println("Bus Arrived here at " + destinationStopArrivalTime)
				log.Println("- - - - - - - - - - - - - - - - - - - - - - -")
			}
			finalStopArrivalTime = currentStop.ArrivalTime
			stops = append(stops, stop)
		}
		route.Stops = stops

		// An empty slice of shapes is created here for each outer iteration for
		// the same reason as the empty slice for the stops above
		shapes = []ShapeJSON{}
		for _, currentShape := range currentRoute.Shapes {
			shape.ShapePtLat, _ = strconv.ParseFloat(currentShape.ShapePtLat, 64)
			shape.ShapePtLon, _ = strconv.ParseFloat(currentShape.ShapePtLon, 64)
			shape.ShapePtSequence = currentShape.ShapePtSequence
			shape.ShapeDistTravel = currentShape.ShapeDistTravel
			shapes = append(shapes, shape)
		}
		route.Shapes = shapes

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
			destinationStopArrivalTime, finalStopArrivalTime)

		route.TravelTime = journeyTravelTime

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

	// Assign values to connection string variables
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

	dateStringSplit := strings.Split(date, " ")
	//dateString := dateStringSplit[0]
	timeString := dateStringSplit[1]

	// Create context variable and assign time for timeout
	// Log any resulting error here also
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Print(err)
	}
	defer client.Disconnect(ctx) // defer has rest of function done before disconnect

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
		stops = []RouteStop{}
		for _, currentStop := range currentRoute.Stops {
			stop.StopId = currentStop.StopId
			stop.StopName = currentStop.StopName
			stop.StopNumber = currentStop.StopNumber
			stop.StopLat, _ = strconv.ParseFloat(currentStop.StopLat, 64)
			stop.StopLon, _ = strconv.ParseFloat(currentStop.StopLon, 64)
			stop.StopSequence = currentStop.StopSequence
			stop.ArrivalTime = currentStop.ArrivalTime
			stop.DepartureTime = currentStop.DepartureTime
			stop.DistanceTravelled, _ =
				strconv.ParseFloat(currentStop.DistanceTravelled, 64)
			if currentStop.StopNumber == origin {
				originStopArrivalTime = currentStop.ArrivalTime
				log.Println("Found the origin stop: " + currentStop.StopNumber)
				log.Println("Origin stop was in the sequence at stop " + currentStop.StopSequence)
				log.Println("Bus Arrived here at " + originStopArrivalTime)
				log.Println("- - - - - - - - - - - - - - - - - - - - - - -")
			}
			if currentStop.StopNumber == destination {
				destinationStopArrivalTime = currentStop.ArrivalTime
				log.Println("Found the destination stop: " + currentStop.StopNumber)
				log.Println("Destination stop was in the sequence at stop " + currentStop.StopSequence)
				log.Println("Bus Arrived here at " + destinationStopArrivalTime)
				log.Println("- - - - - - - - - - - - - - - - - - - - - - -")
			}
			finalStopArrivalTime = currentStop.ArrivalTime
			stops = append(stops, stop)
		}
		route.Stops = stops

		// An empty slice of shapes is created here for each outer iteration for
		// the same reason as the empty slice for the stops above
		shapes = []ShapeJSON{}
		for _, currentShape := range currentRoute.Shapes {
			shape.ShapePtLat, _ = strconv.ParseFloat(currentShape.ShapePtLat, 64)
			shape.ShapePtLon, _ = strconv.ParseFloat(currentShape.ShapePtLon, 64)
			shape.ShapePtSequence = currentShape.ShapePtSequence
			shape.ShapeDistTravel = currentShape.ShapeDistTravel
			shapes = append(shapes, shape)
		}
		route.Shapes = shapes

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
			destinationStopArrivalTime, finalStopArrivalTime)

		route.TravelTime = journeyTravelTime

		resultJSON = append(resultJSON, route)
	}

	return resultJSON
}

// FindMatchingRouteDemo is a demo function designed for development
// purposes to test the functionality of different elements of the route
// matching service without affecting the "stable" functionality present.
// This function is for development only and will be removed before the
// final product is created
func FindMatchingRouteDemo(c *gin.Context) {

	// Assign values to connection string variables
	mongoHost = os.Getenv("MONGO_INITDB_ROOT_HOST")
	mongoPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongoUsername = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPort = os.Getenv("MONGO_INITDB_ROOT_PORT")

	//date := "2022-07-30 14:00:00"

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

	// Create context variable and assign time for timeout
	// Log any resulting error here also
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Print(err)
	}
	defer client.Disconnect(ctx) // defer has rest of function done before disconnect

	coll := client.Database("BusData").Collection("trips_n_stops")
	cursor, err := coll.Aggregate(ctx, bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"stops.stop_number", "7067"},
					{"stops",
						bson.D{
							{"$elemMatch",
								bson.D{
									{"stop_number", "2955"},
									{"departure_time", bson.D{{"$gt", "19:55:00"}}},
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
					{"direction_id", bson.D{{"$first", "$direction_id"}}},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	var result []busRoute
	var resultJSON []busRouteJSON
	var route busRouteJSON
	var stop RouteStop
	var shape ShapeJSON
	var stops []RouteStop
	var shapes []ShapeJSON
	var originStopArrivalTime string
	var destinationStopArrivalTime string

	if err = cursor.All(ctx, &result); err != nil {
		log.Print(err)
	}

	for _, currentRoute := range result {
		route.RouteNum = currentRoute.Id
		stops = []RouteStop{}
		for _, currentStop := range currentRoute.Stops {
			stop.StopId = currentStop.StopId
			stop.StopName = currentStop.StopName
			stop.StopNumber = currentStop.StopNumber
			stop.StopLat, _ = strconv.ParseFloat(currentStop.StopLat, 64)
			stop.StopLon, _ = strconv.ParseFloat(currentStop.StopLon, 64)
			stop.StopSequence = currentStop.StopSequence
			stop.ArrivalTime = currentStop.ArrivalTime
			stop.DepartureTime = currentStop.DepartureTime
			stop.DistanceTravelled, _ =
				strconv.ParseFloat(currentStop.DistanceTravelled, 64)
			stops = append(stops, stop)
			if currentStop.StopNumber == "2955" {
				originStopArrivalTime = currentStop.ArrivalTime
			} else if currentStop.StopNumber == "7067" {
				destinationStopArrivalTime = currentStop.ArrivalTime
			}
		}
		route.Stops = stops
		shapes = []ShapeJSON{}
		for _, currentShape := range currentRoute.Shapes {
			shape.ShapePtLat, _ = strconv.ParseFloat(currentShape.ShapePtLat, 64)
			shape.ShapePtLon, _ = strconv.ParseFloat(currentShape.ShapePtLon, 64)
			shape.ShapePtSequence = currentShape.ShapePtSequence
			shape.ShapeDistTravel = currentShape.ShapeDistTravel
			shapes = append(shapes, shape)
		}
		route.Shapes = shapes

		route.Fares = CalculateFare(currentRoute, "2955", "7067")

		if currentRoute.Direction == "1" {
			route.Direction = "2"
		} else {
			route.Direction = "1"
		}
		log.Println("Origin Arrival Time: " + originStopArrivalTime)
		log.Println("Destination Arrival Time: " + destinationStopArrivalTime)

		//initialTravelTime := GetTravelTimePrediction(route.RouteNum, date, route.Direction)
		//
		//journeyTravelTime := AdjustTravelTime(initialTravelTime, originStopArrivalTime, destinationStopArrivalTime)
		//
		//route.TravelTime = journeyTravelTime

		resultJSON = append(resultJSON, route)
	}

	c.IndentedJSON(http.StatusOK, resultJSON)
}
