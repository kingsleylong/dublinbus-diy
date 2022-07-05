package main

import (
	"example.com/api/databaseQueries"
	"github.com/gin-gonic/gin"
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
	router.GET("/findStopByName", databaseQueries.GetStopByName)

	// Bus Route queries
	router.GET("/busRoute/:routeNum", databaseQueries.GetBusRoute)
	router.GET("/allRoutes", databaseQueries.GetAllRoutes)
	router.GET("/stopsOnRoute/:routeNum", databaseQueries.GetStopsOnRoute)
	router.GET("/matchingRoute/:originStopNum/:destStopNum", databaseQueries.FindMatchingRoute)

	router.Run("0.0.0.0:8080")
}
