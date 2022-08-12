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

var stop RouteStop
var shape ShapeJSON
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
// api call - the origin coordinates pair, the destination coordinates pair,
// the type of time being passed in (either an arrival or departure time) and
// finally the time itself. This function calls the FindMatchingRouteForArrival
// or the FindMatchingRouteForDeparture function depending on the time type
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
		busRoutes := FindMatchingRouteForDeparture(destination, origin, dateAndTime)
		c.IndentedJSON(http.StatusOK, busRoutes)
	} else {
		c.IndentedJSON(http.StatusBadRequest, "Invalid time type parameter in request")
	}

}

// FindMatchingRouteForDeparture takes in the destination coordinates,
// the origin coordinates and the date for a bus trip and returns a slice
// of bus routes that match the given parameters. Its three parameters are all
// taken in as strings and the returned bus routes are of type busRouteJSON. It's
// internal query for the MongoDB database distinguishes it from the
// FindMatchingRouteForArrival function by basing its query on the time
// a bus leaves the origin
func FindMatchingRouteForDeparture(destination string,
	origin string,
	date string) []busRouteJSON {

	// resultJSON kept local so that routes from other calls don't persist
	var resultJSON []busRouteJSON
	var route busRouteJSON

	// First step is taking in coordinates, locating the stops near those
	// coordinates and then returning the 10 closest stops to that initial
	// coordinate pair
	originCoordinates := TurnParameterToCoordinates(origin)
	destinationCoordinates := TurnParameterToCoordinates(destination)

	stopsNearDestination := FindNearbyStopsV2(destinationCoordinates)
	stopsNearOrigin := FindNearbyStopsV2(originCoordinates)

	originStops := CurateNearbyStops(stopsNearOrigin, originCoordinates)
	destinationStops := CurateNearbyStops(stopsNearDestination, destinationCoordinates)

	// Stop numbers for the origin and destination are then extracted from the
	// 10 nearest stops
	originStopNums := []string{}
	for _, originStop := range originStops {
		originStopNums = append(originStopNums, originStop.StopNumber)
	}

	destinationStopNums := []string{}
	for _, destinationStop := range destinationStops {
		destinationStopNums = append(destinationStopNums, destinationStop.StopNumber)
	}

	log.Println("OriginStopNums:")
	log.Println(originStopNums)
	log.Println("")
	log.Println("DestinationStopNums:")
	log.Println(destinationStopNums)
	log.Println("")
	client, err := ConnectToMongo()

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Print(err)
	}
	defer client.Disconnect(ctx) // defer has rest of function complete before disconnect

	// Time of day portion of the date entered extracted here
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
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
	}

	// routes object used to decode the results of the query and prepare for
	// transformation
	var routes []MatchedRoute

	if err = query.All(ctx, &routes); err != nil {
		log.Println(err)
	}

	log.Println("Results from first query:")
	for _, value := range routes {
		log.Println(value)
	}
	log.Println("")
	var routesWithOAndD []MatchedRouteWithOAndD
	var routeWithOAndD MatchedRouteWithOAndD
	for _, matchingRoute := range routes {
		routeWithOAndD.Id = matchingRoute.Id
		routeWithOAndD.Stops = matchingRoute.Stops
		routeWithOAndD.OriginStopNumber, _ = FindNearestStop(originStops,
			matchingRoute.Stops, originCoordinates)
		routeWithOAndD.DestinationStopNumber, _ = FindNearestStop(destinationStops,
			matchingRoute.Stops, destinationCoordinates)
		//log.Println("Matching Route with O and D for route:", routeWithOAndD.Id[0])
		//log.Println(routeWithOAndD)
		routesWithOAndD = append(routesWithOAndD, routeWithOAndD)
	}

	var fullRoutes = []busRoute{}
	var allRoutes = []busRoute{}
	for _, routeDocument := range routesWithOAndD {
		log.Println("Route document id field 0 = " + routeDocument.Id[0])
		log.Println("Route document id field 01 = " + routeDocument.Id[1])
		query, err = collection.Aggregate(ctx, bson.A{
			bson.D{{
				"$match", bson.D{
					{"route.route_short_name", routeDocument.Id[0]},
					{"direction_id", routeDocument.Id[1]},
					{"stops", bson.D{
						{"$elemMatch", bson.D{
							{"stop_number", routeDocument.OriginStopNumber},
							{"departure_time", bson.D{
								{"$gt", timeString},
							},
							}}},
					}}},
			}},
			bson.D{{"$sort", bson.D{{"stops.departure_time", 1}}}},
			bson.D{
				{"$group", bson.D{
					{"_id", "$route.route_short_name"},
					{"direction", bson.D{{"$first", "$direction_id"}}},
					{"stops", bson.D{{"$first", "$stops"}}},
					{"shapes", bson.D{{"$first", "$shapes"}}},
				}},
			},
		})
		if err != nil {
			log.Println("Error from Mongo query")
			log.Println(err)
		}

		if err = query.All(ctx, &fullRoutes); err != nil {
			log.Println("Error reading query result into fullRoutes array")
			log.Println(err)
		}
		log.Println("Number of routes found:", len(fullRoutes))
		for index, _ := range fullRoutes {
			fullRoutes[index].Direction = routeDocument.Id[1]
		}

		allRoutes = append(allRoutes, fullRoutes...)
		log.Println("Full Routes array after reading in query result:")
		log.Println(fullRoutes)
		//log.Println(fullRoutes[0].Direction)
		//log.Println(len(fullRoutes[0].Id))
		//log.Println(string(fullRoutes[0].Id[1]))
	}

	log.Println("All routes in busRoute format:")
	log.Println(allRoutes)
	log.Println("")
	// Iterate over the result objects to transform them into suitable return
	// objects while also generating travel time predictions and fare calculations
	for _, currentRoute := range allRoutes {

		log.Println(string(currentRoute.Id))
		log.Println(currentRoute.Direction)
		// Intermediary object to hold route information with fields for
		// origin and destination stop numbers
		var routeWithOAndD busRouteV2
		routeWithOAndD.Id = string(currentRoute.Id)
		routeWithOAndD.Direction = currentRoute.Direction
		routeWithOAndD.Stops = currentRoute.Stops
		routeWithOAndD.Shapes = currentRoute.Shapes
		log.Println(currentRoute.Direction, routeWithOAndD.Direction)
		log.Println("Route with Origin and Destination:")
		log.Println(routeWithOAndD)
		log.Println("")
		route.RouteNum = string(currentRoute.Id)

		// Two flags used within main loop when checking for matching origin and destination
		originAndDestinationFound := false
		originFound := false

		// Main loop iterates over each stop in the route object
		for _, allStops := range currentRoute.Stops {

			// If origin isn't found then check against the array of origin
			// stops for a match here. Otherwise, ignore this inner loop
			if originFound == false {
				for _, originStop := range originStops {
					if allStops.StopNumber == originStop.StopNumber {
						routeWithOAndD.OriginStopNumber = originStop.StopNumber
						originFound = true
						break
					}
				}
			}

			// If origin has been found and destination hasn't yet, then check against
			// destination stops array for a match. Once found, this step is skipped
			if originFound == true && originAndDestinationFound == false {
				for _, destinationStop := range destinationStops {
					if allStops.StopNumber == destinationStop.StopNumber {
						routeWithOAndD.DestinationStopNumber = destinationStop.StopNumber
						originAndDestinationFound = true
						break
					}
				}
			}

			// Finally check at the end of loop if both stops have been found
			// and if they have, exit the loop
			if originAndDestinationFound == true {
				break
			}
		}

		log.Println("Origin and Destination after loop:")
		log.Println(routeWithOAndD.OriginStopNumber, routeWithOAndD.DestinationStopNumber)
		log.Println("")
		// At the end of the iteration over all stops, if a matching origin
		// and destination were never found, skip remaining steps and move on
		// to next route document in result slice
		if originAndDestinationFound == false {
			continue
		}

		// Slice for stops in route created along with time variables used later
		route.Stops = CreateStopsSlice(routeWithOAndD.OriginStopNumber,
			routeWithOAndD.DestinationStopNumber, currentRoute, stop)

		// Shapes slice created
		route.Shapes = CreateShapesSlice(currentRoute)

		// If the origin and destination were somehow found out of order then
		// skip this iteration and move onto the next route document in result
		// slice
		if originStopSequence > destinationStopSequence {
			continue
		}

		// Use the CalculateFare function from fareCalculation.go to get the fares
		// object for each route
		route.Fares = CalculateFare(currentRoute,
			routeWithOAndD.OriginStopNumber, routeWithOAndD.DestinationStopNumber)

		// Set route direction variable so that it matches necessary direction input
		// for travel time prediction
		if currentRoute.Direction == "1" {
			route.Direction = "2"
		} else {
			route.Direction = "1"
		}

		// Get travel time prediction as floating point numbers based on call to external api
		// connecting to flask application
		initialTravelTime, err := GetTravelTimePrediction(route.RouteNum, date, route.Direction)
		if err != nil {
			log.Println(err)
		}

		// Floating point travel time used in conjunction with static timetable time
		// information to generate more user-friendly travel time information
		journeyTravelTime := AdjustTravelTime(initialTravelTime, originStopArrivalTime,
			destinationStopArrivalTime, firstStopArrivalTime, finalStopArrivalTime)

		// If the travel time prediction could not be calculated then source will have
		// been set to static, where now static timetable information is used for the
		// travel time estimation returned to the user
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

		// The stops slice is finally adjusted so that it only contains stops along the route being
		// travelled
		originStopIndex, destinationStopIndex := CurateStopsSlice(routeWithOAndD.OriginStopNumber,
			routeWithOAndD.DestinationStopNumber, route)
		route.Stops = route.Stops[originStopIndex : destinationStopIndex+1]

		// Static timetable departure time is used to provide the user of an estimate
		// for how when a bus will arrive to begin their journey
		route.TravelTime.ScheduledDepartureTime = GetTimeStringAsHoursAndMinutes(route.Stops[0].ArrivalTime)

		resultJSON = append(resultJSON, route)
	}

	return resultJSON
}

// FindMatchingRouteForArrival takes in the destination coordinates,
// the origin coordinates and the date for a bus trip and returns a slice
// of bus routes that match the given parameters. Its three parameters are all
// taken in as strings and the returned bus routes are of type busRouteJSON. It
// is distinct from the FindMatchingRouteForDeparture function as it bases
// its MongoDb query on the arrival time at the destination stop rather
// than the time to leave the origin stop
func FindMatchingRouteForArrival(origin string,
	destination string,
	date string) []busRouteJSON {

	// resultJSON kept local so that routes from other calls don't persist
	var resultJSON []busRouteJSON
	var route busRouteJSON
	// First step is taking in coordinates, locating the stops near those
	// coordinates and then returning the 10 closest stops to that initial
	// coordinate pair
	originCoordinates := TurnParameterToCoordinates(origin)
	destinationCoordinates := TurnParameterToCoordinates(destination)

	stopsNearDestination := FindNearbyStopsV2(destinationCoordinates)
	stopsNearOrigin := FindNearbyStopsV2(originCoordinates)

	originStops := CurateNearbyStops(stopsNearOrigin, originCoordinates)
	destinationStops := CurateNearbyStops(stopsNearDestination, destinationCoordinates)

	// Stop numbers for the origin and destination are then extracted from the
	// 10 nearest stops
	originStopNums := []string{}
	for _, originStop := range originStops {
		originStopNums = append(originStopNums, originStop.StopNumber)
	}

	destinationStopNums := []string{}
	for _, destinationStop := range destinationStops {
		destinationStopNums = append(destinationStopNums, destinationStop.StopNumber)
	}

	log.Println("OriginStopNums:")
	log.Println(originStopNums)
	log.Println("")
	log.Println("DestinationStopNums:")
	log.Println(destinationStopNums)
	log.Println("")
	client, err := ConnectToMongo()

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Print(err)
	}
	defer client.Disconnect(ctx) // defer has rest of function complete before disconnect

	// Time of day portion of the date entered extracted here
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
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
	}

	// routes object used to decode the results of the query and prepare for
	// transformation
	var routes []MatchedRoute

	if err = query.All(ctx, &routes); err != nil {
		log.Println(err)
	}

	log.Println("Results from first query:")
	for _, value := range routes {
		log.Println(value)
	}
	log.Println("")
	var routesWithOAndD []MatchedRouteWithOAndD
	var routeWithOAndD MatchedRouteWithOAndD
	for _, matchingRoute := range routes {
		routeWithOAndD.Id = matchingRoute.Id
		routeWithOAndD.Stops = matchingRoute.Stops
		routeWithOAndD.OriginStopNumber, _ = FindNearestStop(originStops,
			matchingRoute.Stops, originCoordinates)
		routeWithOAndD.DestinationStopNumber, _ = FindNearestStop(destinationStops,
			matchingRoute.Stops, destinationCoordinates)
		routesWithOAndD = append(routesWithOAndD, routeWithOAndD)
	}

	var fullRoutes []busRoute
	var allRoutes []busRoute
	for _, routeDocument := range routesWithOAndD {
		query, err = collection.Aggregate(ctx, bson.A{
			bson.D{{
				"$match", bson.D{
					{"route.route_short_name", routeDocument.Id[0]},
					{"direction_id", routeDocument.Id[1]},
					{"stops", bson.D{
						{"$elemMatch", bson.D{
							{"stop_number", routeDocument.DestinationStopNumber},
							{"arrival_time", bson.D{
								{"$lte", timeString},
							},
							}}},
					}}},
			}},
			bson.D{{"$sort", bson.D{{"stops.arrival_time", -1}}}},
			bson.D{
				{"$group", bson.D{
					{"_id", "$route.route_short_name"},
					{"stops", bson.D{{"$first", "$stops"}}},
					{"shapes", bson.D{{"$first", "$shapes"}}},
					{"direction", bson.D{{"$first", "$direction_id"}}},
				}},
			},
		})
		if err != nil {
			log.Println("Error from Mongo query")
			log.Println(err)
		}

		if err = query.All(ctx, &fullRoutes); err != nil {
			log.Println("Error reading query result into fullRoutes array")
			log.Println(err)
		}

		allRoutes = append(allRoutes, fullRoutes...)

	}

	log.Println("All routes in busRoute format:")
	log.Println(allRoutes)
	log.Println("")
	// Iterate over the result objects to transform them into suitable return
	// objects while also generating travel time predictions and fare calculations
	for _, currentRoute := range allRoutes {

		log.Println(string(currentRoute.Id))
		// Intermediary object to hold route information with fields for
		// origin and destination stop numbers
		var routeWithOAndD busRouteV2
		routeWithOAndD.Id = string(currentRoute.Id)
		routeWithOAndD.Stops = currentRoute.Stops
		routeWithOAndD.Shapes = currentRoute.Shapes
		routeWithOAndD.Direction = currentRoute.Direction
		log.Println(currentRoute.Direction, routeWithOAndD.Direction)
		log.Println("Route with Origin and Destination:")
		log.Println(routeWithOAndD)
		log.Println("")
		route.RouteNum = string(currentRoute.Id)

		// Two flags used within main loop when checking for matching origin and destination
		originAndDestinationFound := false
		originFound := false

		// Main loop iterates over each stop in the route object
		for _, allStops := range currentRoute.Stops {

			// If origin isn't found then check against the array of origin
			// stops for a match here. Otherwise, ignore this inner loop
			if originFound == false {
				for _, originStop := range originStops {
					if allStops.StopNumber == originStop.StopNumber {
						routeWithOAndD.OriginStopNumber = originStop.StopNumber
						originFound = true
						break
					}
				}
			}

			// If origin has been found and destination hasn't yet, then check against
			// destination stops array for a match. Once found, this step is skipped
			if originFound == true && originAndDestinationFound == false {
				for _, destinationStop := range destinationStops {
					if allStops.StopNumber == destinationStop.StopNumber {
						routeWithOAndD.DestinationStopNumber = destinationStop.StopNumber
						originAndDestinationFound = true
						break
					}
				}
			}

			// Finally check at the end of loop if both stops have been found
			// and if they have, exit the loop
			if originAndDestinationFound == true {
				break
			}
		}

		log.Println("Route with O and D after loop to check for origin and destination:")
		log.Println(routeWithOAndD)
		log.Println("")
		// At the end of the iteration over all stops, if a matching origin
		// and destination were never found, skip remaining steps and move on
		// to next route document in result slice
		if originAndDestinationFound == false {
			continue
		}

		// Slice for stops in route created along with time variables used later
		route.Stops = CreateStopsSlice(routeWithOAndD.OriginStopNumber,
			routeWithOAndD.DestinationStopNumber, currentRoute, stop)

		// Shapes slice created
		route.Shapes = CreateShapesSlice(currentRoute)

		// If the origin and destination were somehow found out of order then
		// skip this iteration and move onto the next route document in result
		// slice
		if originStopSequence > destinationStopSequence {
			continue
		}

		// Use the CalculateFare function from fareCalculation.go to get the fares
		// object for each route
		route.Fares = CalculateFare(currentRoute,
			routeWithOAndD.OriginStopNumber, routeWithOAndD.DestinationStopNumber)

		// Set route direction variable so that it matches necessary direction input
		// for travel time prediction
		if currentRoute.Direction == "1" {
			route.Direction = "2"
		} else {
			route.Direction = "1"
		}

		// Get travel time prediction as floating point numbers based on call to external api
		// connecting to flask application
		initialTravelTime, err := GetTravelTimePrediction(route.RouteNum, date, route.Direction)
		if err != nil {
			log.Println(err)
		}

		// Floating point travel time used in conjunction with static timetable time
		// information to generate more user-friendly travel time information
		journeyTravelTime := AdjustTravelTime(initialTravelTime, originStopArrivalTime,
			destinationStopArrivalTime, firstStopArrivalTime, finalStopArrivalTime)

		// If the travel time prediction could not be calculated then source will have
		// been set to static, where now static timetable information is used for the
		// travel time estimation returned to the user
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

		// The stops slice is finally adjusted so that it only contains stops along the route being
		// travelled
		originStopIndex, destinationStopIndex := CurateStopsSlice(routeWithOAndD.OriginStopNumber,
			routeWithOAndD.DestinationStopNumber, route)
		route.Stops = route.Stops[originStopIndex : destinationStopIndex+1]

		// Static timetable departure time is used to provide the user of an estimate
		// for how when a bus will arrive to begin their journey
		route.TravelTime.ScheduledDepartureTime = GetTimeStringAsHoursAndMinutes(route.Stops[0].ArrivalTime)

		resultJSON = append(resultJSON, route)
	}

	resultJSON = CurateReturnedArrivalRoutes(date, resultJSON)
	return resultJSON
}
