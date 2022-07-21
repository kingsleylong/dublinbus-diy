package databaseQueries

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"

	//"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
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

func GetStopsList(c *gin.Context) {

	stopSearch := c.Param("stopSearch")
	stopsFromDB := GetStopByName(stopSearch)
	stopsFromGeocoding := FindNearbyStops(stopSearch)

	var busStops findByAddressResponse

	busStops.Matched = stopsFromDB
	busStops.Nearby = stopsFromGeocoding

	c.IndentedJSON(http.StatusOK, busStops)
}
