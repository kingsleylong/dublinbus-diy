package databaseQueries

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// FindMatchingRoute is a function that takes in four parameters for its
// api call - the origin bus stop, the destination bus stop, the type of
// time being passed in (either an arrival or departure time) and finally
// the time itself. This function calls the FindMatchingRouteForArrival or
// the FindMatchingRouteForDeparture function depending on the time type
// that is passed in and then returns an array of busRouteJSON type containing
// the routes found that match the query. It may also return a status 400 with
// the appropriate string message if the time type passed in is invalid
func FindMatchingRoute(c *gin.Context) {

	origin := c.Param("origin")
	destination := c.Param("destination")
	timeType := c.Param("timeType")
	time := c.Param("time")

	if timeType == "arrival" {
		busRoutes := FindMatchingRouteForArrival(origin, destination, time)
		c.IndentedJSON(http.StatusOK, busRoutes)
	} else if timeType == "departure" {
		busRoutes := FindMatchingRouteForDeparture(destination, origin, time)
		c.IndentedJSON(http.StatusOK, busRoutes)
	} else {
		c.IndentedJSON(http.StatusBadRequest, "Invalid time type parameter in request")
	}

}
