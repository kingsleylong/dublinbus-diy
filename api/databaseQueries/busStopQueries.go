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

type busStop struct {
	StopId     string `bson:"stop_id" json:"stop_id"`
	StopName   string `bson:"stop_name" json:"stop_name"`
	StopNumber string `bson:"stop_number" json:"stop_number"`
	StopLat    string `bson:"stop_lat" json:"stop_lat"`
	StopLon    string `bson:"stop_lon" json:"stop_lon"`
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

// GetBusStop returns a single JSON object representing a bus stop from the
// MongoDB instance. This JSON object includes name, number
// and coordinates of the stop. Note that this object is slightly different to
// the bus stop objects present in functions to describe stops on a given route.
func GetBusStop(c *gin.Context) {

	// Assign values to connection string variables
	mongoHost = os.Getenv("MONGO_INITDB_ROOT_HOST")
	mongoPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongoUsername = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPort = os.Getenv("MONGO_INITDB_ROOT_PORT")

	// Read in bus stop number parameter provided in URL
	stopNum := c.Param("stopNum")

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

	var result bson.M

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("stops")

	// Find one document that matches criteria and decode results into result address
	err = collectionPointer.FindOne(ctx, bson.D{{"stop_number", string(stopNum)}}).
		Decode(&result)

	// Check for any errors associated with the query and if there are, log them
	if err != nil {
		// Specific condition for valid query without matching documents
		if err == mongo.ErrNoDocuments {
			log.Print("No matching document found")
			return
		}
		log.Print(err)
	}

	// Return result as JSON along with code 200
	c.IndentedJSON(http.StatusOK, result)
}

// GetAllStops returns an array of JSON objects. Each object is a bus stop object as
// per the GetBusStop function. Each stop object includes the stop name, number and coordinates.
func GetAllStops(c *gin.Context) {

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

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("stops")

	// Return all documents in the collection by leaving filter parameters empty.
	// Filter parameters specified in the bson.D{{}}
	busStops, err := collectionPointer.Find(ctx, bson.D{{}})
	if err != nil {
		log.Print(err)
	}

	var busStopResults []bson.M

	// Check for an error while reading through all the documents found in
	// the query and as they are read, append them to the busStopResults slice
	// using that slice's memory address reference
	if err = busStops.All(ctx, &busStopResults); err != nil {
		log.Print(err)
	}

	// Return result slice as JSON along with code 200
	c.IndentedJSON(http.StatusOK, busStopResults)
}

// GetPrototypeStops is a function to return a curated list of bus stops for
// the purposes of the presentation and live demo on 6th July 22. As such, this
// function is incomplete at this time and is not supported. Due for deletion before
// final build.
func GetPrototypeStops(c *gin.Context) {

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

	var result []bson.M

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("stopsOnRoute")

	busRoutes, err := collectionPointer.Find(ctx, bson.D{{"route_stops.stop_num", "774"}})
	if err != nil {
		log.Print(err)
	}

	if err = busRoutes.All(ctx, &result); err != nil {
		log.Print(err)
	}

	var busStopJSONArray []busRoute

	busStopsBytes, err := bson.Marshal(result)

	if err := bson.Unmarshal(busStopsBytes, &busStopJSONArray); err != nil {
		log.Print(err)
	}

	c.IndentedJSON(http.StatusOK, busStopJSONArray)
}

// GetStopByName takes a string passed into the request URL and then
// searches the database for a bus stop with a name that matches. For all
// the stops with a matching name or similar name, these stops are
// returned as JSON objects from the stops collection in MongoDB.
func GetStopByName(c *gin.Context) {

	mongoHost = os.Getenv("MONGO_INITDB_ROOT_HOST")
	mongoPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongoUsername = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPort = os.Getenv("MONGO_INITDB_ROOT_PORT")

	stopName := c.Param("stopName")

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

	var matchingStops []busStop

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("stops")

	busStops, err := collectionPointer.Find(ctx, bson.D{{"stop_name", bson.M{
		"$regex": primitive.Regex{Pattern: "^" + stopName + "*", Options: "i"}}}})
	if err != nil {
		log.Print(err)
	}

	if err = busStops.All(ctx, &matchingStops); err != nil {
		log.Print(err)
	}

	c.IndentedJSON(http.StatusOK, matchingStops)
}
