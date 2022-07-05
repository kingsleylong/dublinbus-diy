package databaseQueries

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

// busRoute is a type that is designed to read from the stopsOnRoute collection
// from MongoDB. It contains an id field, a string that specifies the route number
// as a Dublin Bus user would recognise it and finally an array of busStop structs.
type busRoute struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	RouteNum   string             `bson:"route_num" json:"route_num"`
	RouteStops []busStop          `bson:"route_stops" json:"route_stops"`
}

// busStop is a struct containing information about each of the bus stop objects
// nested within the stopsOnRoute collection in MongoDB. These objects include the number
// of the stop (the number of the stop and not its technical id value), the address
// and location of the stop and finally the stop's coordinates.
type busStop struct {
	StopNum      string `bson:"stop_num" json:"stop_num"`
	StopAddress  string `bson:"stop_address" json:"stop_address"`
	StopLocation string `bson:"stop_location" json:"stop_location"`
	StopLat      string `bson:"stop_lat" json:"stop_lat"`
	StopLon      string `bson:"stop_lon" json:"stop_lon"`
}

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

	// Leave filter values empty to retrieve all documents in this collection
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

	// Map resulting information to busRoute struct
	var result busRoute

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("stopsOnRoute")

	// Find one document that matches criteria and decode results into result address
	err = collectionPointer.FindOne(ctx, bson.D{{"route_num", string(routeNum)}}).
		Decode(&result)

	// Return result as JSON along with code 200
	c.IndentedJSON(http.StatusOK, result)
}

// FindMatchingRoute takes in two parameters (the origin and destination bus stop number)
// and then this function attempts to find the bus route objects(s) that contain both the
// origin and destination stop and then returns these specific routes as JSON.
func FindMatchingRoute(c *gin.Context) {

	// Assign values to connection string variables
	mongoHost = os.Getenv("MONGO_INITDB_ROOT_HOST")
	mongoPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongoUsername = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPort = os.Getenv("MONGO_INITDB_ROOT_PORT")

	// Read in route number parameter provided in URL
	originStopNum := c.Param("originStopNum")
	destStopNum := c.Param("destStopNum")

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

	// Arrays to hold routes for the origin and destination stops
	var originResult []busRoute
	var destinationResult []busRoute

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("stopsOnRoute")

	// Find documents that have the required origin stop as a stop on the route
	// and store these routes in array
	originBusRoutes, err := collectionPointer.Find(ctx, bson.D{{"route_stops.stop_num", originStopNum}})
	if err != nil {
		log.Print(err)
	}

	if err = originBusRoutes.All(ctx, &originResult); err != nil {
		log.Print(err)
	}

	// Repeat above procedure but for the destination stop
	destBusRoutes, err := collectionPointer.Find(ctx, bson.D{{"route_stops.stop_num", destStopNum}})
	if err != nil {
		log.Print(err)
	}

	if err = destBusRoutes.All(ctx, &destinationResult); err != nil {
		log.Print(err)
	}

	// Loop through the origin route objects and then within that loop examine the destination
	// route objects and check for matching route numbers. If there is a match, store the matched
	// objects in a final array which is returned as JSON.
	var matchedRoutes []busRoute

	for _, originRoute := range originResult {
		for _, destRoute := range destinationResult {
			if destRoute.RouteNum == originRoute.RouteNum {
				matchedRoutes = append(matchedRoutes, destRoute)
			}
		}
	}
	c.IndentedJSON(http.StatusOK, matchedRoutes)
}
