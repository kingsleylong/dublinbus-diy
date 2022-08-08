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

	// Bus Stop specific queries
	router.GET("/databases", databaseQueries.GetDatabases)
	router.GET("/stop/findByAddress/:stopSearch", databaseQueries.GetStopsList)

	// Bus Route queries
	router.GET("route/matchingRoute/:origin/:destination/:timeType/:time",
		databaseQueries.FindMatchingRoute)

	router.GET("route/matchingRouteV2Test/:origin/:destination/:timeType/:time")

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
}
