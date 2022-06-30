package main

import (
	"example.com/api/databaseQueries"
	"github.com/gin-gonic/gin"
)

// Main function contains the routed URIs mapped to functions and starts
// the server engine
func main() {

	router := gin.Default()

	router.GET("/databases", databaseQueries.GetDatabases)
	router.GET("/busStop/:stopNum",
		databaseQueries.GetBusStop)
	router.GET("/allStops", databaseQueries.GetAllStops)
	router.GET("/busRoute/:routeNum", databaseQueries.GetBusRoute)
	router.GET("/allRoutes", databaseQueries.GetAllRoutes)

	router.Run("0.0.0.0:8080")
}
