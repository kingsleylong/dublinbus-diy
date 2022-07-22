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
	"time"
)

// FindMatchingRouteForDeparture takes in two parameters (the origin and destination bus stop number)
// and then this function attempts to find the bus route objects(s) that contain both the
// origin and destination stop and then returns these specific routes as JSON.
func FindMatchingRouteForDeparture(destination string,
	origin string,
	departureTime string) []busRouteJSON {

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
										bson.D{{"$gt", departureTime}}},
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

	var result []busRoute
	var resultJSON []busRouteJSON
	var route busRouteJSON
	var stop RouteStop
	if err = cursor.All(ctx, &result); err != nil {
		log.Print(err)
	}

	for _, currentRoute := range result {
		route.ID = currentRoute.Id
		for _, currentStop := range currentRoute.Stops {
			stop.StopId = currentStop.StopId
			stop.StopName = currentStop.StopName
			stop.StopNumber = currentStop.StopNumber
			stop.StopLat, _ = strconv.ParseFloat(currentStop.StopLat, 64)
			stop.StopLon, _ = strconv.ParseFloat(currentStop.StopLon, 64)
			stop.StopSequence = currentStop.StopSequence
			stop.ArrivalTime = currentStop.ArrivalTime
			stop.DepartureTime = currentStop.DepartureTime
			route.Stops = append(route.Stops, stop)
		}
		route.Shapes = currentRoute.Shapes
		resultJSON = append(resultJSON, route)
	}

	return resultJSON
}

func FindMatchingRouteForArrival(origin string,
	destination string,
	arrivalTime string) []busRouteJSON {

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
										bson.D{{"$lte", arrivalTime}}},
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

	var result []busRoute
	var resultJSON []busRouteJSON
	var route busRouteJSON
	var stop RouteStop
	if err = cursor.All(ctx, &result); err != nil {
		log.Print(err)
	}

	for _, currentRoute := range result {
		route.ID = currentRoute.Id
		for _, currentStop := range currentRoute.Stops {
			stop.StopId = currentStop.StopId
			stop.StopName = currentStop.StopName
			stop.StopNumber = currentStop.StopNumber
			stop.StopLat, _ = strconv.ParseFloat(currentStop.StopLat, 64)
			stop.StopLon, _ = strconv.ParseFloat(currentStop.StopLon, 64)
			stop.StopSequence = currentStop.StopSequence
			stop.ArrivalTime = currentStop.ArrivalTime
			stop.DepartureTime = currentStop.DepartureTime
			route.Stops = append(route.Stops, stop)
		}
		route.Shapes = currentRoute.Shapes
		resultJSON = append(resultJSON, route)
	}

	return resultJSON
}

func FindMatchingRouteDemo(c *gin.Context) {

	// Assign values to connection string variables
	mongoHost = os.Getenv("MONGO_INITDB_ROOT_HOST")
	mongoPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongoUsername = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPort = os.Getenv("MONGO_INITDB_ROOT_PORT")

	// Read in route number parameter provided in URL
	//	originStopNum := c.Param("originStopNum")

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
					{"stops.stop_number", "4727"},
					{"stops",
						bson.D{
							{"$elemMatch",
								bson.D{
									{"stop_number", "2070"},
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
	if err = cursor.All(ctx, &result); err != nil {
		log.Print(err)
	}

	for _, currentRoute := range result {
		route.ID = currentRoute.Id
		for _, currentStop := range currentRoute.Stops {
			stop.StopId = currentStop.StopId
			stop.StopName = currentStop.StopName
			stop.StopNumber = currentStop.StopNumber
			stop.StopLat, _ = strconv.ParseFloat(currentStop.StopLat, 64)
			stop.StopLon, _ = strconv.ParseFloat(currentStop.StopLon, 64)
			stop.StopSequence = currentStop.StopSequence
			stop.ArrivalTime = currentStop.ArrivalTime
			stop.DepartureTime = currentStop.DepartureTime
			route.Stops = append(route.Stops, stop)
		}
		route.Shapes = currentRoute.Shapes
		resultJSON = append(resultJSON, route)
	}
	
	c.IndentedJSON(http.StatusOK, resultJSON)
}
