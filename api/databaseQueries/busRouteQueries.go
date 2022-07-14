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

// busRoute is a type that is designed to read from the trips_new collection
// from MongoDB. It contains the id fields that combine to form a unique key
// for each entry (i.e. the route id, the shape id, the direction id, the trip id
// and the stop id). It also includes coordinates for the stop associated with this
// object as well as information on the route and the shape string used to
// draw the shape on the map. All fields map to type string from the database
type busRoute struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	RouteId     string             `bson:"route_id" json:"route_id"`
	TripId      string             `bson:"trip_id" json:"trip_id"`
	ShapeId     string             `bson:"shape_id" json:"shape_id"`
	DirectionId string             `bson:"direction_id" json:"direction_id"`
	Route       route              `bson:"route" json:"route"`
	Shapes      []shape            `bson:"shapes" json:"shapes"`
	Stops       []routeStop        `bson:"stops" json:"stops"`
}

// routeStop represents the stop information contained within the trips_n_stops
// collection in MongoDB. The information contains the StopId that can be used
// to identify each stop uniquely, the name of that stop, the stop number used
// by consumers of the Dublin Bus service, the coordinates of
// the stop that can be used to mark the stop on the map and finally a sequence
// number that can be used to sort the stops to ensure that they are in the
// correct order on a given route. All fields are returned as strings from the
// database
type routeStop struct {
	StopId       string `bson:"stop_id" json:"stop_id"`
	StopName     string `bson:"stop_name" json:"stop_name"`
	StopNumber   string `bson:"stop_number" json:"stop_number"`
	StopLat      string `bson:"stop_lat" json:"stop_lat"`
	StopLon      string `bson:"stop_lon" json:"stop_lon"`
	StopSequence string `bson:"stop_sequence" json:"stop_sequence"`
}

// route is a struct that contains a means of matching the route number (referred to
// as RouteShortName in this object) to the route id (i.e. the RouteId). All
// fields map to type string from the database
type route struct {
	RouteId        string `bson:"route_id" json:"route_id"`
	RouteShortName string `bson:"route_short_name" json:"route_short_name"`
}

// shape is struct that contains the coordinates for each turn in a bus
// line as it travels its designated route that combined together allow
// the bus route to be drawn on a map matching the road network of Dublin.
// All fields map to type string from the database
type shape struct {
	ShapeId         string `bson:"shape_id" json:"shape_id"`
	ShapePtLat      string `bson:"shape_pt_lat" json:"shape_pt_lat"`
	ShapePtLon      string `bson:"shape_pt_lon" json:"shape_pt_lon"`
	ShapePtSequence string `bson:"shape_pt_sequence" json:"shape_pt_sequence"`
	ShapeDistTravel string `bson:"shape_dist_traveled" json:"shape_dist_traveled"`
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
	var matchingRoutes []busRoute
	var originRoute busRoute
	var destinationRoute busRoute

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("trips_n_stops")

	// Find documents that have the required origin stop as a stop on the route
	// and store these routes in array
	originBusRoutes, err := collectionPointer.Find(ctx, bson.D{{"stops.stop_number",
		string(originStopNum)}})
	if err != nil {
		log.Print(err)
	}

	for originBusRoutes.Next(ctx) {
		originBusRoutes.Decode(&originRoute)
		originRoutes = append(originRoutes, originRoute)
	}

	destinationBusRoutes, err := collectionPointer.Find(ctx, bson.D{{"stops.stop_number",
		string(destStopNum)}})
	if err != nil {
		log.Print(err)
	}

	for destinationBusRoutes.Next(ctx) {
		destinationBusRoutes.Decode(&destinationRoute)
		destinationRoutes = append(destinationRoutes, destinationRoute)
	}

	for _, origin := range originRoutes {
		for _, destination := range destinationRoutes {
			if destination.RouteId == origin.RouteId {
				matchingRoutes = append(matchingRoutes, destination)
				break
			}
		}
	}

	c.IndentedJSON(http.StatusOK, matchingRoutes)
}

func FindMatchingRouteDemo(c *gin.Context) {

	// Assign values to connection string variables
	mongoHost = os.Getenv("MONGO_INITDB_ROOT_HOST")
	mongoPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongoUsername = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPort = os.Getenv("MONGO_INITDB_ROOT_PORT")

	// Read in route number parameter provided in URL
	//	originStopNum := c.Param("originStopNum")

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
	var matchingRoutes []busRoute
	var originRoute busRoute
	var destinationRoute busRoute

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("trips_n_stops")

	// Find documents that have the required origin stop as a stop on the route
	// and store these routes in array
	originBusRoutes, err := collectionPointer.Find(ctx, bson.D{{"stops.stop_number",
		"2955"}})
	if err != nil {
		log.Print(err)
	}

	for originBusRoutes.Next(ctx) {
		originBusRoutes.Decode(&originRoute)
		originRoutes = append(originRoutes, originRoute)
	}

	destinationBusRoutes, err := collectionPointer.Find(ctx, bson.D{{"stops.stop_number",
		"7067"}})
	if err != nil {
		log.Print(err)
	}

	for destinationBusRoutes.Next(ctx) {
		destinationBusRoutes.Decode(&destinationRoute)
		destinationRoutes = append(destinationRoutes, destinationRoute)
	}

	for _, origin := range originRoutes {
		for _, destination := range destinationRoutes {
			if destination.RouteId == origin.RouteId {
				matchingRoutes = append(matchingRoutes, destination)
				break
			}
		}
	}

	c.IndentedJSON(http.StatusOK, matchingRoutes)
}
