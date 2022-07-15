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
	router.GET("/findStopByName/:stopName", databaseQueries.GetStopByName)

	// Bus Route queries
	router.GET("/matchingRoute/:originStopNum/:destStopNum", databaseQueries.FindMatchingRoute)
	router.GET("/matchingRouteTest/:originStopNum", databaseQueries.FindMatchingRouteDemo)

	router.GET("/geocoding/:address", geocoding.GetCoordinates)

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
}
