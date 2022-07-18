package databaseQueries

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"

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

type BusStop struct {
	StopId     string  `bson:"stop_id" json:"stop_id"`
	StopName   string  `bson:"stop_name" json:"stop_name"`
	StopNumber string  `bson:"stop_number" json:"stop_number"`
	StopLat    float64 `bson:"stop_lat" json:"stop_lat"`
	StopLon    float64 `bson:"stop_lon" json:"stop_lon"`
}

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
func GetStopByName(stopName string) []BusStop {

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

	var matchingStops []BusStop

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
	for busStops.Next(ctx) {
		if err := busStops.Decode(&stop); err != nil {
			log.Print(err)
		}
		matchingStops = append(matchingStops, stop)
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

	var busStops [][]BusStop

	busStops = append(busStops, stopsFromDB, stopsFromGeocoding)

	c.IndentedJSON(http.StatusOK, busStops)
}
