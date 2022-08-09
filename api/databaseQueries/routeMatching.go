package databaseQueries

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

// Global Package variables

// Variables of both busRoute and busRouteJSON need to be initialised as
// some unmarshalling from Mongo cannot be done automatically and
// so must be done manually from one structure to another in the backend
var result []busRoute
var resultJSON []busRouteJSON
var route busRouteJSON
var stop RouteStop
var shape ShapeJSON

//var stops []RouteStop
var shapes []ShapeJSON
var originStopArrivalTime string
var destinationStopArrivalTime string
var finalStopArrivalTime string
var firstStopArrivalTime string
var originStopSequence int64
var destinationStopSequence int64
var originDistTravelled float64
var destinationDistTravelled float64

// FindMatchingRoute is a function that takes in four parameters for its
// api call - the origin bus stop, the destination bus stop, the type of
// time being passed in (either an arrival or departure time) and finally
// the time itself. This function calls the FindMatchingRouteForArrival or
// the FindMatchingRouteForDeparture function depending on the time type
// that is passed in and then returns an array of busRouteJSON type containing
// the routes found that match the query. It may also return a status 400 with
// the appropriate string message if the time type passed in is invalid
func FindMatchingRoute(c *gin.Context) {

	origin := c.Param("origin")
	destination := c.Param("destination")
	timeType := c.Param("timeType")
	dateAndTime := c.Param("time")

	if timeType == "arrival" {
		busRoutes := FindMatchingRouteForArrival(origin, destination, dateAndTime)
		c.IndentedJSON(http.StatusOK, busRoutes)
	} else if timeType == "departure" {
		busRoutes := FindMatchingRouteForDepartureV2(destination, origin, dateAndTime)
		c.IndentedJSON(http.StatusOK, busRoutes)
	} else {
		c.IndentedJSON(http.StatusBadRequest, "Invalid time type parameter in request")
	}

}

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
		route.Shapes = CreateShapesSlice(currentRoute)

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
			journeyTravelTime.EstimatedArrivalTime = destinationStopArrivalTime
			journeyTravelTime.EstimatedArrivalTime = destinationStopArrivalTime
			journeyTravelTime.EstimatedArrivalTime = destinationStopArrivalTime
		}
		route.TravelTime = journeyTravelTime

		originStopIndex, destinationStopIndex := CurateStopsSlice(origin, destination)

		route.Stops = route.Stops[originStopIndex : destinationStopIndex+1]
		route.TravelTime.ScheduledDepartureTime = GetScheduledDepartureTime(route.Stops[0].ArrivalTime)

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
		route.Shapes = CreateShapesSlice(currentRoute)

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
			journeyTravelTime.EstimatedArrivalTime = destinationStopArrivalTime
			journeyTravelTime.EstimatedArrivalTime = destinationStopArrivalTime
			journeyTravelTime.EstimatedArrivalTime = destinationStopArrivalTime
		}
		route.TravelTime = journeyTravelTime

		originStopIndex, destinationStopIndex := CurateStopsSlice(origin, destination)
		route.Stops = route.Stops[originStopIndex : destinationStopIndex+1]
		route.TravelTime.ScheduledDepartureTime = GetScheduledDepartureTime(route.Stops[0].ArrivalTime)

		resultJSON = append(resultJSON, route)
	}

	resultJSON = CurateReturnedArrivalRoutes(date, resultJSON)

	return resultJSON
}

func FindMatchingRouteForDepartureV2(destination string,
	origin string,
	date string) []busRouteJSON {

	originCoordinates := TurnParameterToCoordinates(origin)
	destinationCoordinates := TurnParameterToCoordinates(destination)
	log.Println("Origin Coordinates:")
	log.Print(originCoordinates)
	log.Println("Destination Coordinates:")
	log.Println(destinationCoordinates)
	log.Println("---- ---- ---- ---- ---- ---- ----")

	stopsNearDestination := FindNearbyStopsV2(destinationCoordinates)
	stopsNearOrigin := FindNearbyStopsV2(originCoordinates)
	log.Println("stopsNearDestination:")
	log.Println(stopsNearDestination)
	log.Println("stopsNearOrigin")
	log.Println(stopsNearOrigin)
	log.Println("---- ---- ---- ---- ---- ---- ----")

	originStops := CurateNearbyStops(stopsNearOrigin, originCoordinates)
	log.Println("Origin Stops:")
	log.Println(originStops)
	log.Println("")
	destinationStops := CurateNearbyStops(stopsNearDestination, destinationCoordinates)
	log.Println("Destination Stops:")
	log.Println(destinationStops)
	log.Println("")

	originStopNums := []string{}
	for _, originStop := range originStops {
		originStopNums = append(originStopNums, originStop.StopNumber)
	}

	destinationStopNums := []string{}
	for _, destinationStop := range destinationStops {
		destinationStopNums = append(destinationStopNums, destinationStop.StopNumber)
	}
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
	collection := client.Database("BusData").Collection("trips_n_stops")

	query, err := collection.Aggregate(ctx, bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"stops",
						bson.D{
							{"$elemMatch",
								bson.D{
									{"stop_number",
										bson.D{
											{"$in",
												originStopNums,
											},
										},
									},
									{"departure_time", bson.D{{"$gt", timeString}}},
								},
							},
						},
					},
					{"stops.stop_number",
						bson.D{
							{"$in",
								destinationStopNums,
							},
						},
					},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"stops.departure_time", 1}}}},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$route.route_short_name"},
					{"stops", bson.D{{"$first", "$stops"}}},
					{"shapes", bson.D{{"$first", "$shapes"}}},
					{"direction", bson.D{{"$first", "$direction_id"}}},
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
	}

	var routes []busRoute

	if err = query.All(ctx, &routes); err != nil {
		log.Println(err)
	}

	var routesFound string

	for _, routeNum := range routes {
		routesFound += routeNum.Id + " "
	}

	log.Println("Routes found: " + routesFound)
	for _, currentRoute := range routes {

		var originStopNumber string
		var destinationStopNumber string
		var routeCursor busRoute
		routeCursor.Id = currentRoute.Id
		routeCursor.Stops = currentRoute.Stops
		routeCursor.Shapes = currentRoute.Shapes
		routeCursor.Direction = currentRoute.Direction

		route.RouteNum = currentRoute.Id

		originAndDestinationFound := false
		for _, checkOriginStop := range route.Stops {
			for _, originStop := range originStops {
				if checkOriginStop.StopNumber == originStop.StopNumber {
					log.Println("Found origin stop: " + originStop.StopNumber)
					log.Println("")
					for _, checkDestinationStop := range route.Stops {
						for _, destinationStop := range destinationStops {
							if checkDestinationStop.StopNumber == destinationStop.StopNumber {
								log.Println("Found destination stop: " + destinationStop.StopNumber)
								originStopNumber = originStop.StopNumber
								destinationStopNumber = destinationStop.StopNumber
								originAndDestinationFound = true
								break
							}
						}
						if originAndDestinationFound {
							break
						}
					}
				}
				if originAndDestinationFound {
					break
				}
			}
			if originAndDestinationFound {
				break
			}
		}
		// An empty slice of stops is created with each new outer iteration so
		// that duplicates aren't added to later routes in their stop arrays
		route.Stops = CreateStopsSlice(originStopNumber,
			destinationStopNumber, routeCursor, stop)

		// An empty slice of shapes is created here for each outer iteration for
		// the same reason as the empty slice for the stops above
		route.Shapes = CreateShapesSlice(routeCursor)

		if originStopSequence > destinationStopSequence {
			continue
		}

		// Use the CalculateFare function from fareCalculation.go to get the fares
		// object for each route
		route.Fares = CalculateFare(routeCursor,
			originStopNumber, destinationStopNumber)

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
			journeyTravelTime.EstimatedArrivalTime = destinationStopArrivalTime
			journeyTravelTime.EstimatedArrivalTime = destinationStopArrivalTime
			journeyTravelTime.EstimatedArrivalTime = destinationStopArrivalTime
		}
		route.TravelTime = journeyTravelTime

		originStopIndex, destinationStopIndex := CurateStopsSlice(origin, destination)
		route.Stops = route.Stops[originStopIndex : destinationStopIndex+1]
		route.TravelTime.ScheduledDepartureTime = GetScheduledDepartureTime(route.Stops[0].ArrivalTime)

		resultJSON = append(resultJSON, route)
	}

	return resultJSON
}
