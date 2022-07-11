package main

import (
	"example.com/api/databaseQueries"
	"example.com/api/geocoding"
	"github.com/gin-gonic/gin"
	"log"
)

// Main function contains the routed URIs mapped to functions and starts
// the server engine
func main() {

	router := gin.Default()

	// Bus Stop specific queries
	router.GET("/databases", databaseQueries.GetDatabases)
	router.GET("/busStop/:stopNum",
		databaseQueries.GetBusStop)
	router.GET("/allStops", databaseQueries.GetAllStops)
	router.GET("/prototypeStops", databaseQueries.GetPrototypeStops)
	router.GET("/findStopByName/:stopName", databaseQueries.GetStopByName)

	// Bus Route queries
	router.GET("/busRoute/:routeNum", databaseQueries.GetBusRoute)
	router.GET("/allRoutes", databaseQueries.GetAllRoutes)
	router.GET("/stopsOnRoute/:routeNum", databaseQueries.GetStopsOnRoute)
	router.GET("/matchingRoute/:originStopNum/:destStopNum", databaseQueries.FindMatchingRoute)

	router.GET("/geocoding/:address", geocoding.GetCoordinates)

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
}
