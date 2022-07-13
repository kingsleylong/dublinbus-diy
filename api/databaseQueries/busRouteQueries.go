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
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	RouteId     string             `bson:"route_id" json:"route_id"`
	TripId      string             `bson:"trip_id" json:"trip_id"`
	ShapeId     string             `bson:"shape_id" json:"shape_id"`
	DirectionId string             `bson:"direction_id" json:"direction_id"`
	Route       []route            `bson:"route" json:"route"`
	Shapes      []shapes           `bson:"shapes" json:"shapes"`
	RouteStops  []routeBusStop     `bson:"stops" json:"stops"`
}

// routeBusStop is a struct containing information about each of the bus stop objects
// nested within the stopsOnRoute collection in MongoDB. These objects include the number
// of the stop (the number of the stop and not its technical id value), the address
// and location of the stop and finally the stop's coordinates.
type routeBusStop struct {
	StopId       string `bson:"stop_id" json:"stop_id"`
	StopSequence string `bson:"stop_sequence" json:"stop_sequence"`
	StopHeadsign string `bson:"stop_headsign" json:"stop_headsign"`
}

type route struct {
	RouteId        string `bson:"route_id" json:"route_id"`
	RouteShortName string `bson:"route_short_name" json:"route_short_name"`
}

type shapes struct {
	ShapeId         string `bson:"shape_id" json:"shape_id"`
	ShapePtLat      string `bson:"shape_pt_lat" json:"shape_pt_lat"`
	ShapePtLon      string `bson:"shape_pt_lon" json:"shape_pt_lon"`
	ShapePtSequence string `bson:"shape_pt_sequence" json:"shape_pt_sequence"`
	ShapeDistTravel string `bson:"shape_dist_traveled" json:"shape_dist_traveled"`
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
	var route busRoute
	var stops []busStop
	var stop busStop

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("trips_new")

	// Find one document that matches criteria and decode results into result address
	err = collectionPointer.FindOne(ctx, bson.D{{"route.route_short_name", routeNum}}).
		Decode(&route)

	collectionPointer = dbPointer.Collection("stops")
	busStops, err := collectionPointer.Find(ctx, bson.D{{}})
	if err != nil {
		log.Print(err)
	}

	for busStops.Next(ctx) {
		if err := busStops.Decode(&stop); err != nil {
			log.Print(err)
		}
		for _, stopValue := range route.RouteStops {
			if stop.StopId == stopValue.StopId {
				stops = append(stops, stop)
			}
		}
	}

	// Return result as JSON along with code 200
	c.IndentedJSON(http.StatusOK, stops)
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
	var originRoutes []busRoute
	var destinationRoutes []busRoute
	var originStop busStop
	var destinationStop busStop

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("stops")

	// Find documents that have the required origin stop as a stop on the route
	// and store these routes in array
	err = collectionPointer.FindOne(ctx, bson.D{{"stop_number", originStopNum}}).
		Decode(&originStop)
	if err != nil {
		log.Print(err)
	}

	err = collectionPointer.FindOne(ctx, bson.D{{"stop_number", destStopNum}}).
		Decode(&destinationStop)
	if err != nil {
		log.Print(err)
	}

	collectionPointer = dbPointer.Collection("trips_new")

	originBusRoutes, err := collectionPointer.Find(ctx, bson.D{{"stops.stop_id", originStop.StopId}})
	if err != nil {
		log.Print(err)
	}

	if err := originBusRoutes.All(ctx, &originRoutes); err != nil {
		log.Print(err)
	}

	// Repeat above procedure but for the destination stop
	destinationBusRoutes, err := collectionPointer.Find(ctx, bson.D{{"stops.stop_id",
		destinationStop.StopId}})
	if err != nil {
		log.Print(err)
	}

	if err := destinationBusRoutes.All(ctx, &destinationRoutes); err != nil {
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
