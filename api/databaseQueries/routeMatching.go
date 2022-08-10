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
var route busRouteJSON
var stop RouteStop
var shape ShapeJSON
var originStopNumber string
var destinationStopNumber string

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
		busRoutes := FindMatchingRouteForArrivalV2(origin, destination, dateAndTime)
		c.IndentedJSON(http.StatusOK, busRoutes)
	} else if timeType == "departure" {
		busRoutes := FindMatchingRouteForDepartureV2(destination, origin, dateAndTime)
		c.IndentedJSON(http.StatusOK, busRoutes)
	} else {
		c.IndentedJSON(http.StatusBadRequest, "Invalid time type parameter in request")
	}

}

func FindMatchingRouteForDepartureV2(destination string,
	origin string,
	date string) []busRouteJSON {

	var resultJSON []busRouteJSON

	originCoordinates := TurnParameterToCoordinates(origin)
	destinationCoordinates := TurnParameterToCoordinates(destination)

	stopsNearDestination := FindNearbyStopsV2(destinationCoordinates)
	stopsNearOrigin := FindNearbyStopsV2(originCoordinates)

	originStops := CurateNearbyStops(stopsNearOrigin, originCoordinates)
	destinationStops := CurateNearbyStops(stopsNearDestination, destinationCoordinates)

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
					{"$and",
						bson.A{
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
												{"arrival_time", bson.D{{"$gte", timeString}}},
											},
										},
									},
								},
							},
							bson.D{
								{"stops",
									bson.D{
										{"$elemMatch",
											bson.D{
												{"stop_number",
													bson.D{
														{"$in",
															destinationStopNums,
														},
													},
												},
												{"arrival_time", bson.D{{"$gte", timeString}}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"stops.arrival_time", 1}}}},
		bson.D{
			{"$group",
				bson.D{
					{"_id",
						bson.A{
							"$route.route_short_name",
							"$direction_id",
						},
					},
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

	for _, currentRoute := range routes {

		var routeWithOAndD busRouteV2
		routeWithOAndD.Id = currentRoute.Id[0]
		routeWithOAndD.Stops = currentRoute.Stops
		routeWithOAndD.Shapes = currentRoute.Shapes
		routeWithOAndD.Direction = currentRoute.Direction

		route.RouteNum = currentRoute.Id[0]

		originAndDestinationFound := false
		originFound := false
		for _, allStops := range currentRoute.Stops {
			if originFound == false {
				for _, originStop := range originStops {
					if allStops.StopNumber == originStop.StopNumber {
						routeWithOAndD.OriginStopNumber = originStop.StopNumber
						originFound = true
						break
					}
				}
			}
			if originFound == true && originAndDestinationFound == false {
				for _, destinationStop := range destinationStops {
					if allStops.StopNumber == destinationStop.StopNumber {
						routeWithOAndD.DestinationStopNumber = destinationStop.StopNumber
						originAndDestinationFound = true
						break
					}
				}
			}
			if originAndDestinationFound == true {
				break
			}
		}

		if originAndDestinationFound == false {
			continue
		}
		route.Stops = CreateStopsSlice(routeWithOAndD.OriginStopNumber,
			routeWithOAndD.DestinationStopNumber, currentRoute, stop)

		// An empty slice of shapes is created here for each outer iteration for
		// the same reason as the empty slice for the stops above
		route.Shapes = CreateShapesSlice(currentRoute)

		if originStopSequence > destinationStopSequence {
			continue
		}

		// Use the CalculateFare function from fareCalculation.go to get the fares
		// object for each route
		route.Fares = CalculateFare(currentRoute,
			routeWithOAndD.OriginStopNumber, routeWithOAndD.DestinationStopNumber)

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
			journeyTravelTime.EstimatedArrivalTime = GetTimeStringAsHoursAndMinutes(destinationStopArrivalTime)
			journeyTravelTime.EstimatedArrivalHighTime = GetTimeStringAsHoursAndMinutes(destinationStopArrivalTime)
			journeyTravelTime.EstimatedArrivalLowTime = GetTimeStringAsHoursAndMinutes(destinationStopArrivalTime)
		}
		route.TravelTime = journeyTravelTime

		originStopIndex, destinationStopIndex := CurateStopsSlice(routeWithOAndD.OriginStopNumber,
			routeWithOAndD.DestinationStopNumber)
		route.Stops = route.Stops[originStopIndex : destinationStopIndex+1]
		route.TravelTime.ScheduledDepartureTime = GetTimeStringAsHoursAndMinutes(route.Stops[0].ArrivalTime)

		resultJSON = append(resultJSON, route)
	}

	return resultJSON
}

func FindMatchingRouteForArrivalV2(origin string,
	destination string,
	date string) []busRouteJSON {

	var resultJSON []busRouteJSON

	originCoordinates := TurnParameterToCoordinates(origin)
	destinationCoordinates := TurnParameterToCoordinates(destination)

	stopsNearDestination := FindNearbyStopsV2(destinationCoordinates)
	stopsNearOrigin := FindNearbyStopsV2(originCoordinates)

	originStops := CurateNearbyStops(stopsNearOrigin, originCoordinates)
	destinationStops := CurateNearbyStops(stopsNearDestination, destinationCoordinates)

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
	log.Println(timeString)

	// Aggregation pipeline created in Mongo Compass and then transformed to suit
	// the mongo driver in Go
	collection := client.Database("BusData").Collection("trips_n_stops")

	query, err := collection.Aggregate(ctx, bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"$and",
						bson.A{
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
												{"arrival_time", bson.D{{"$lte", timeString}}},
											},
										},
									},
								},
							},
							bson.D{
								{"stops",
									bson.D{
										{"$elemMatch",
											bson.D{
												{"stop_number",
													bson.D{
														{"$in",
															destinationStopNums,
														},
													},
												},
												{"arrival_time", bson.D{{"$lte", timeString}}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"stops.arrival_time", -1}}}},
		bson.D{
			{"$group",
				bson.D{
					{"_id",
						bson.A{
							"$route.route_short_name",
							"$direction_id",
						},
					},
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

	for _, currentRoute := range routes {

		var routeWithOAndD busRouteV2
		routeWithOAndD.Id = currentRoute.Id[0]
		routeWithOAndD.Stops = currentRoute.Stops
		routeWithOAndD.Shapes = currentRoute.Shapes
		routeWithOAndD.Direction = currentRoute.Direction

		route.RouteNum = currentRoute.Id[0]

		originAndDestinationFound := false
		originFound := false
		for _, allStops := range currentRoute.Stops {
			if originFound == false {
				for _, originStop := range originStops {
					if allStops.StopNumber == originStop.StopNumber {
						routeWithOAndD.OriginStopNumber = originStop.StopNumber
						originFound = true
						break
					}
				}
			}
			if originFound == true && originAndDestinationFound == false {
				for _, destinationStop := range destinationStops {
					if allStops.StopNumber == destinationStop.StopNumber {
						routeWithOAndD.DestinationStopNumber = destinationStop.StopNumber
						originAndDestinationFound = true
						break
					}
				}
			}
			if originAndDestinationFound == true {
				break
			}
		}

		if originAndDestinationFound == false {
			continue
		}
		route.Stops = CreateStopsSlice(routeWithOAndD.OriginStopNumber,
			routeWithOAndD.DestinationStopNumber, currentRoute, stop)

		route.Shapes = CreateShapesSlice(currentRoute)

		if originStopSequence > destinationStopSequence {
			continue
		}

		// Use the CalculateFare function from fareCalculation.go to get the fares
		// object for each route
		route.Fares = CalculateFare(currentRoute,
			routeWithOAndD.OriginStopNumber, routeWithOAndD.DestinationStopNumber)

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
			journeyTravelTime.EstimatedArrivalTime = GetTimeStringAsHoursAndMinutes(destinationStopArrivalTime)
			journeyTravelTime.EstimatedArrivalHighTime = GetTimeStringAsHoursAndMinutes(destinationStopArrivalTime)
			journeyTravelTime.EstimatedArrivalLowTime = GetTimeStringAsHoursAndMinutes(destinationStopArrivalTime)
		}
		route.TravelTime = journeyTravelTime

		originStopIndex, destinationStopIndex := CurateStopsSlice(routeWithOAndD.OriginStopNumber,
			routeWithOAndD.DestinationStopNumber)
		route.Stops = route.Stops[originStopIndex : destinationStopIndex+1]
		route.TravelTime.ScheduledDepartureTime = GetTimeStringAsHoursAndMinutes(route.Stops[0].ArrivalTime)

		resultJSON = append(resultJSON, route)
	}

	resultJSON = CurateReturnedArrivalRoutes(date, resultJSON)

	return resultJSON
}
