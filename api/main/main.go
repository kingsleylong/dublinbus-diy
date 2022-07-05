package main

import (
	"example.com/api/databaseQueries"
	"github.com/gin-gonic/gin"
	"log"
)

// Main function contains the routed URIs mapped to functions and starts
// the server engine
func main() {

	router := gin.Default()

	router.GET("/databases", databaseQueries.GetDatabases)
	router.GET("/busStop/:stopNum",
		databaseQueries.GetBusStop)
	router.GET("/allStops", databaseQueries.GetAllStops)
	router.GET("/prototypeStops", databaseQueries.GetPrototypeStops)
	router.GET("/busRoute/:routeNum", databaseQueries.GetBusRoute)
	router.GET("/allRoutes", databaseQueries.GetAllRoutes)
	router.GET("/stopsOnRoute/:routeNum", databaseQueries.GetStopsOnRoute)
	router.GET("/matchingRoute/:originStopNum/:destStopNum", databaseQueries.FindMatchingRoute)

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
}
