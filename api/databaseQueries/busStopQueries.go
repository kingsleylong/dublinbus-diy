package databaseQueries

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"googlemaps.github.io/maps"
	"math"
	"sort"
	"strconv"

	//"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

// Variables to hold connection string values
var mongoUsername string
var mongoPassword string
var mongoHost string
var mongoPort string

// GetDatabases returns the databases present in the MongoDB connection.
// Useful as a debugging query.
func GetDatabases(c *gin.Context) {

	client, err := ConnectToMongo()

	// Create context variable and assign time for timeout
	// Log any resulting error here also
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Print(err)
	}
	defer client.Disconnect(ctx) // defer has rest of function complete before this disconnect

	// Create list of databases and return as JSON
	databases, err := client.ListDatabases(ctx, bson.D{})

	c.IndentedJSON(http.StatusOK, databases)
}

// GetStopByName takes a string passed into the request URL and then
// searches the database for a bus stop with a name that matches. For all
// the stops with a matching name or similar name, these stops are
// returned as JSON objects from the stops collection in MongoDB.
func GetStopByName(stopName string) []StopWithCoordinates {

	client, err := ConnectToMongo()

	// Create context variable and assign time for timeout
	// Log any resulting error here also
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Print(err)
	}
	defer client.Disconnect(ctx) // defer has rest of function done before disconnect

	var matchingStops []StopWithCoordinates

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("stops")

	// Use regex to search for a pattern in the bus stop names to locate stops
	// with similar names to help users find stops by their name
	busStops, err := collectionPointer.Find(ctx, bson.D{{"stop_name", bson.M{
		"$regex": primitive.Regex{Pattern: stopName, Options: "i"}}}})
	if err != nil {
		log.Print(err)
	}

	// Iteratively go through returned options and add them to slice to return
	// until slice length hits the limit and then stop the loop
	var stop BusStop
	var stopWithCoordinates StopWithCoordinates
	for busStops.Next(ctx) {
		if err := busStops.Decode(&stop); err != nil {
			log.Print(err)
		}
		stopWithCoordinates.StopID = stop.StopId
		stopWithCoordinates.StopNumber = stop.StopNumber
		stopWithCoordinates.StopName = stop.StopName
		stopWithCoordinates.StopLat, _ = strconv.ParseFloat(stop.StopLat, 64)
		stopWithCoordinates.StopLon, _ = strconv.ParseFloat(stop.StopLon, 64)
		matchingStops = append(matchingStops, stopWithCoordinates)
		if len(matchingStops) > 4 {
			break
		}
	}

	return matchingStops
}

// GetStopsList is a function that maps to the findByAddress API which
// returns a pair of arrays that contain stops found by regex patterns
// with stop names in our database and stops found that are nearby using
// geolocation respectively
func GetStopsList(c *gin.Context) {

	stopSearch := c.Param("stopSearch")
	stopsFromDB := GetStopByName(stopSearch)
	stopsFromGeocoding := FindNearbyStops(stopSearch)

	var busStops findByAddressResponse

	busStops.Matched = stopsFromDB
	busStops.Nearby = stopsFromGeocoding

	c.IndentedJSON(http.StatusOK, busStops)
}

// CurateNearbyStops is a function that is used to sort the stops that are
// determined to be close to the coordinates provided in the variable 'location'
// in order of distance from that point. Once this has been done, the ten
// nearest stops are returned as a result of this function from using a subslice
// of the then sorted input list of stops
func CurateNearbyStops(stopsList []StopWithCoordinates, location maps.LatLng) []StopWithCoordinates {

	closestStops := []StopWithCoordinates{}

	// Distance from one point to another on a 2d plane is the root of
	// (x2-x1)^2 + (y2-y1)^2. This is given to the sort function that is built
	// in to then determine the order in which to sort the stops
	sort.Slice(stopsList, func(i, j int) bool {
		distanceForPointI := math.Sqrt(math.Pow(stopsList[i].StopLon-location.Lng, 2) +
			math.Pow(stopsList[i].StopLat-location.Lat, 2))
		distanceForPointJ := math.Sqrt(math.Pow(stopsList[j].StopLon-location.Lng, 2) +
			math.Pow(stopsList[j].StopLat-location.Lat, 2))
		return distanceForPointI < distanceForPointJ
	})

	if len(stopsList) < 10 {
		closestStops = stopsList
	} else {
		closestStops = append(closestStops, stopsList[:10]...)
	}

	return closestStops
}
