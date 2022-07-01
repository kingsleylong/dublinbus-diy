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
	"time"
)

// GetBusRoute queries the database for a single bus route and returns
// a JSON object representing that route. Includes the route number used
// by bus services as well as route id value for historical GTFS-R data
func GetBusRoute(c *gin.Context) {

	// Assign values to connection string variables
	mongoHost = os.Getenv("MONGO_INITDB_ROOT_HOST")
	mongoPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongoUsername = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPort = os.Getenv("MONGO_INITDB_ROOT_PORT")

	// Read in route number parameter provided in URL
	routeNum := c.Param("routeNum")

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
	collectionPointer := dbPointer.Collection("routes")

	// Find one document that matches criteria and decode results into result address
	err = collectionPointer.FindOne(ctx, bson.D{{"route_short_name", string(routeNum)}}).
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

// GetAllRoutes returns an array of bus route objects as per those
// returned above. These JSON objects include the bus route short name
// used by service users as well as the route id to match with historical data
func GetAllRoutes(c *gin.Context) {

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

	var busRoutesResult []bson.M

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("routes")

	// Find one document that matches criteria and decode results into result address
	busRoutes, err := collectionPointer.Find(ctx, bson.D{{}})
	if err != nil {
		log.Print(err)
	}

	// Check for any errors associated with the query and if there are, log them
	if err = busRoutes.All(ctx, &busRoutesResult); err != nil {
		log.Print(err)
	}

	// Return result as JSON along with code 200
	c.IndentedJSON(http.StatusOK, busRoutesResult)
}

// GetStopsOnRoute returns an array of JSON objects representing bus stops
// that lie along a particular bus route. These objects include bus stop names,
// numbers and coordinates.
func GetStopsOnRoute(c *gin.Context) {

	// Assign values to connection string variables
	mongoHost = os.Getenv("MONGO_INITDB_ROOT_HOST")
	mongoPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongoUsername = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPort = os.Getenv("MONGO_INITDB_ROOT_PORT")

	// Read in route number parameter provided in URL
	routeNum := c.Param("routeNum")

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

	var result bson.M

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("routes")

	// Find one document that matches criteria and decode results into result address
	err = collectionPointer.FindOne(ctx, bson.D{{"route_short_name", string(routeNum)}}).
		Decode(&result)

	var busStopTimesResults []bson.M

	dbPointer = client.Database("BusData")
	collectionPointer = dbPointer.Collection("storeGtfrs")

	// Find one document that matches criteria and decode results into result address
	busStopTimes, err := collectionPointer.Find(ctx, bson.D{
		{"Entity.TripUpdate.Trip.RouteId", result["route_id"]}})
	if err != nil {
		log.Print(err)
	}

	if err = busStopTimes.All(ctx, &busStopTimesResults); err != nil {
		log.Print(err)
	}

	// Return result as JSON along with code 200
	c.IndentedJSON(http.StatusOK, busStopTimesResults)
}
