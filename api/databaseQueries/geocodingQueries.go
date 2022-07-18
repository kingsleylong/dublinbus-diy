package databaseQueries

import (
	"context"
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

var DublinMapBoundsNE maps.LatLng
var DublinMapBoundsSW maps.LatLng
var DublinMapBounds maps.LatLngBounds

func GetCoordinates(stopSearch string) (Lat float64, Lon float64) {

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

	geo := &maps.GeocodingRequest{Address: stopSearch, Bounds: &DublinMapBounds, Region: "ie"}

	result, _ := client.Geocode(ctx, geo)

	queryLat := result[0].Geometry.Location.Lat
	queryLon := result[0].Geometry.Location.Lng

	return queryLat, queryLon
}

func FindNearbyStops(stopSearch string) []BusStop {

	queryLat, queryLon := GetCoordinates(stopSearch)

	halfMileAdjustment := 0.008

	minLat := queryLat - halfMileAdjustment
	maxLat := queryLat + halfMileAdjustment
	minLon := queryLon - halfMileAdjustment
	maxLon := queryLon + halfMileAdjustment

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

	var matchingStops []BusStop
	var currentStop BusStop

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("stops")

	stops, err := collectionPointer.Find(ctx, bson.D{{}})
	if err != nil {
		log.Print(err)
	}

	for stops.Next(ctx) {
		stops.Decode(&currentStop)
		currentLat, _ := strconv.ParseFloat(currentStop.StopLat, 64)
		currentLon, _ := strconv.ParseFloat(currentStop.StopLon, 64)
		if currentLon > minLon && currentLat > minLat {
			if currentLat < maxLat && currentLon < maxLon {
				matchingStops = append(matchingStops, currentStop)
			}
		}
	}

	return matchingStops
}

func GetCoordinatesTest(c *gin.Context) {

	searchString := c.Param("searchString")

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

	geo := &maps.GeocodingRequest{Address: searchString, Bounds: &DublinMapBounds, Region: "ie"}

	result, _ := client.Geocode(ctx, geo)

	queryLat := result[0].Geometry.Location.Lat
	queryLon := result[0].Geometry.Location.Lng

	var latLngTest maps.LatLng

	latLngTest.Lat = queryLat
	latLngTest.Lng = queryLon

	c.IndentedJSON(http.StatusOK, latLngTest)
}
