package databaseQueries

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"googlemaps.github.io/maps"
	"log"
	"os"
	"strconv"
	"time"
)

// Initialise some variables for setting the geocoding boundary in order
// to prevent the service from returning areas outside of Dublin

var DublinMapBoundsNE maps.LatLng
var DublinMapBoundsSW maps.LatLng
var DublinMapBounds maps.LatLngBounds

// GetCoordinates is a function used in the geocoding service that
// is designed to take in a string representation of an address, be
// it a single keyword or multiple words together (hyphenated) and
// then return floating point numbers to represent the latitude and
// longitude of the address respectively. It depends upon the Google
// Maps geocoding service as an external dependency.
func GetCoordinates(stopSearch string) (Lat float64, Lon float64) {

	// Set the values for the boundaries which will help to create a search grid
	// for the geocoding service

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

	// The geocoding process is done entirely externally and only providing the
	// necessary variables is of concern with this function call

	geo := &maps.GeocodingRequest{Address: stopSearch, Bounds: &DublinMapBounds, Region: "ie"}

	result, _ := client.Geocode(ctx, geo)

	// This is a very simple check to make sure that some kind of response was
	// returned. The geocoding request always returns an array with either nothing
	// in it (meaning that an address couldn't be geocoded) or an array with one
	// element with the geocoding information contained within
	if len(result) < 1 {
		return 0, 0
	}

	queryLat := result[0].Geometry.Location.Lat
	queryLon := result[0].Geometry.Location.Lng

	return queryLat, queryLon
}

// FindNearbyStops is function that takes the coordinates returned from
// the GetCoordinates function and then uses that to search within a square
// (approximately) for a bus stop within the database. This function takes
// in a string representing the address being searched, which is passed
// to GetCoordinates and then returns a slice of the structure
// StopWithCoordinates that contains all the identifying information
// about a stop as well as its coordinates
func FindNearbyStops(stopSearch string) []StopWithCoordinates {

	queryLat, queryLon := GetCoordinates(stopSearch)

	halfMileAdjustment := 0.008

	minLat := queryLat - halfMileAdjustment
	maxLat := queryLat + halfMileAdjustment
	minLon := queryLon - halfMileAdjustment
	maxLon := queryLon + halfMileAdjustment

	client, err := ConnectToMongo()

	// Create context variable and assign time for timeout
	// Log any resulting error here also
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Print(err)
	}
	defer client.Disconnect(ctx) // defer has rest of function complete before this disconnect

	var matchingStops []StopWithCoordinates
	var currentStop BusStop
	var currentStopWithCoordinates StopWithCoordinates

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("stops")

	stops, err := collectionPointer.Find(ctx, bson.D{{}})
	if err != nil {
		log.Print(err)
	}

	// The coordinates from the database are read in a string
	// representation and so can't be automatically unmarshalled
	// into floats, so they have to be read in as a BusStop before
	// being read in as a StopWithCoordinates
	for stops.Next(ctx) {
		stops.Decode(&currentStop)
		currentLat, _ := strconv.ParseFloat(currentStop.StopLat, 64)
		currentLon, _ := strconv.ParseFloat(currentStop.StopLon, 64)
		if currentLon > minLon && currentLat > minLat {
			if currentLat < maxLat && currentLon < maxLon {
				currentStopWithCoordinates.StopID = currentStop.StopId
				currentStopWithCoordinates.StopNumber = currentStop.StopNumber
				currentStopWithCoordinates.StopName = currentStop.StopName
				currentStopWithCoordinates.StopLat = currentLat
				currentStopWithCoordinates.StopLon = currentLon
				matchingStops = append(matchingStops, currentStopWithCoordinates)
			}
		}
		if len(matchingStops) >= 5 {
			break
		}
	}

	return matchingStops
}
