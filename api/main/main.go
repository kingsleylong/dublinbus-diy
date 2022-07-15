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
	router.GET("/stopsList/:stopSearch", databaseQueries.GetStopsList)

	// Bus Route queries
	router.GET("/matchingRoute/:originStopNum/:destStopNum", databaseQueries.FindMatchingRoute)
	router.GET("/matchingRouteTest/:originStopNum", databaseQueries.FindMatchingRouteDemo)

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
}
