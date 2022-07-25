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

// FindMatchingRouteForDeparture takes in three parameters - the destination
// bus stop, the origin bus stop and then the departure time all as strings.
// This function then queries the mongo collection for trips documents that
// match these filters before mapping the documents to the correct structure
// and returning them within a slice of type busRouteJSON.
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

	// Variables of both busRoute and busRouteJSON need to be initialised as
	// some unmarshalling from Mongo cannot be done automatically and
	// so must be done manually from one structure to another in the backend
	var result []busRoute
	var resultJSON []busRouteJSON
	var route busRouteJSON
	var stop RouteStop
	var shape ShapeJSON
	if err = cursor.All(ctx, &result); err != nil {
		log.Print(err)
	}

	for _, currentRoute := range result {
		route.RouteNum = currentRoute.Id
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
			route.Stops = append(route.Stops, stop)
		}
		for _, currentShape := range currentRoute.Shapes {
			shape.ShapePtLat, _ = strconv.ParseFloat(currentShape.ShapePtLat, 64)
			shape.ShapePtLon, _ = strconv.ParseFloat(currentShape.ShapePtLon, 64)
			shape.ShapePtSequence = currentShape.ShapePtSequence
			shape.ShapeDistTravel = currentShape.ShapeDistTravel
			route.Shapes = append(route.Shapes, shape)
		}
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

	// Variables of both busRoute and busRouteJSON need to be initialised as
	// some unmarshalling from Mongo cannot be done automatically and
	// so must be done manually from one structure to another in the backend
	var result []busRoute
	var resultJSON []busRouteJSON
	var route busRouteJSON
	var stop RouteStop
	var shape ShapeJSON
	if err = cursor.All(ctx, &result); err != nil {
		log.Print(err)
	}

	for _, currentRoute := range result {
		route.RouteNum = currentRoute.Id
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
			route.Stops = append(route.Stops, stop)
		}
		for _, currentShape := range currentRoute.Shapes {
			shape.ShapePtLat, _ = strconv.ParseFloat(currentShape.ShapePtLat, 64)
			shape.ShapePtLon, _ = strconv.ParseFloat(currentShape.ShapePtLon, 64)
			shape.ShapePtSequence = currentShape.ShapePtSequence
			shape.ShapeDistTravel = currentShape.ShapeDistTravel
			route.Shapes = append(route.Shapes, shape)
		}
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
	var shape ShapeJSON
	if err = cursor.All(ctx, &result); err != nil {
		log.Print(err)
	}

	for _, currentRoute := range result {
		route.RouteNum = currentRoute.Id
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
			route.Stops = append(route.Stops, stop)
		}
		for _, currentShape := range currentRoute.Shapes {
			shape.ShapePtLat, _ = strconv.ParseFloat(currentShape.ShapePtLat, 64)
			shape.ShapePtLon, _ = strconv.ParseFloat(currentShape.ShapePtLon, 64)
			shape.ShapePtSequence = currentShape.ShapePtSequence
			shape.ShapeDistTravel = currentShape.ShapeDistTravel
			route.Shapes = append(route.Shapes, shape)
		}
		resultJSON = append(resultJSON, route)
	}

	// IF DESTINATION (SHAPE_DIST_TRAVELLED) - ORIGIN (SHAPE_DIST_TRAVELLED) > 3000
	// PRINT THE LONGZONE FEE, ELSE PRINT THE SHORTZONE FEE
	//IF ROUTE_SHORT_NAME CONTAINS THE XPRESSO ROUTES, PRINT THE XPRESSO FEE

	var originDist float64
	var destDist float64
	//var fee busRouteJSON

	const (
		XpressoAdult   = float32(3.00)
		ShortZoneAdult = float32(1.7)
		LongZoneAdult  = float32(2.60)
	)
	//	xpress routes names variable - 46A added just for testing
	XpressRoutes := []string{"46a", "27x", "33d", "33x", "39x", "41x", "51x", "51d", "51x", "69x", "77x", "84x"}

	for _, calcRoute := range resultJSON {
		for i := range XpressRoutes {
			if calcRoute.ID == XpressRoutes[i] {
				xpressFee := Fare{XpressoAdult}
				route.FareCalculation = append(route.FareCalculation, xpressFee)
				resultJSON = append(resultJSON, route)
				break
			}

			for _, calcStop := range calcRoute.Stops {
				if calcStop.StopNumber == "4727" {
					originDist = calcStop.DistanceTravelled
					continue
				} else if calcStop.StopNumber == "2070" {
					destDist = calcStop.DistanceTravelled
					continue
				}
				totalDist := destDist - originDist

				ShortZoneDistance := 3000.0

				if totalDist < ShortZoneDistance {
					distFee := Fare{ShortZoneAdult}
					route.FareCalculation = append(route.FareCalculation, distFee)
					resultJSON = append(resultJSON, route)
					continue
					//c.IndentedJSON(http.StatusOK, ShortZoneAdult)
				} else if totalDist >= ShortZoneDistance {
					distFee := Fare{LongZoneAdult}
					route.FareCalculation = append(route.FareCalculation, distFee)
					resultJSON = append(resultJSON, route)
					break
				}
			}
		}
	}

	c.IndentedJSON(http.StatusOK, resultJSON)
}
