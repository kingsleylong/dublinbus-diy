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

	router.Run("localhost:8080")
}
