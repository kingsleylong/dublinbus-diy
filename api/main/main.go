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
	router.GET("/findStopByName/:stopName", databaseQueries.GetStopByName)

	// Bus Route queries
	router.GET("/matchingRoute/:originStopNum/:destStopNum", databaseQueries.FindMatchingRoute)

	router.Run("0.0.0.0:8080")
}
