package geocoding

import (
	"context"
	"example.com/api/databaseQueries"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"googlemaps.github.io/maps"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var mongoUsername string
var mongoPassword string
var mongoHost string
var mongoPort string

var DublinMapBoundsNE maps.LatLng
var DublinMapBoundsSW maps.LatLng
var DublinMapBounds maps.LatLngBounds

func GetCoordinates(c *gin.Context) {

	address := c.Param("address")

	DublinMapBoundsNE.Lat = 53.49337
	DublinMapBoundsNE.Lng = -6.05788

	DublinMapBoundsSW.Lng = -6.56495
	DublinMapBoundsSW.Lat = 53.14860

	DublinMapBounds.SouthWest = DublinMapBoundsSW
	DublinMapBounds.NorthEast = DublinMapBoundsNE

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	client, err := maps.NewClient(maps.WithAPIKey(os.Getenv("MAPS_API_KEY")))
	if err != nil {
		log.Print(err)
	}

	geo := &maps.GeocodingRequest{Address: address, Bounds: &DublinMapBounds, Region: "ie"}

	result, _ := client.Geocode(ctx, geo)

	c.IndentedJSON(http.StatusOK, result)
}

func FindNearbyStops(stopLat float64, stopLon float64) []databaseQueries.BusStop {

	halfMileAdjustment := 0.008

	minLat := stopLat - halfMileAdjustment
	maxLat := stopLat + halfMileAdjustment
	minLon := stopLon - halfMileAdjustment
	maxLon := stopLon + halfMileAdjustment

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

	var matchingStops []databaseQueries.BusStop
	var currentStop databaseQueries.BusStop

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("stops")

	stops, err := collectionPointer.Find(ctx, bson.D{{}})
	if err != nil {
		log.Print(err)
	}

	for stops.Next(ctx) {
		stops.Decode(&currentStop)
		queryLat, _ := strconv.ParseFloat(currentStop.StopLat, 64)
		queryLon, _ := strconv.ParseFloat(currentStop.StopLon, 64)
		if queryLon > minLon && queryLat > minLat {
			if queryLat < maxLat && queryLon < maxLon {
				matchingStops = append(matchingStops, currentStop)
			}
		}
	}

	return matchingStops
}
