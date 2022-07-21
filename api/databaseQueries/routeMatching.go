package databaseQueries

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func FindMatchingRoute(c *gin.Context) {

	origin := c.Param("origin")
	destination := c.Param("destination")
	timeType := c.Param("timeType")
	time := c.Param("time")

	if timeType == "arrival" {
		busRoutes := FindMatchingRouteForArrival(origin, destination, time)
		c.IndentedJSON(http.StatusOK, busRoutes)
	} else if timeType == "destination" {
		busRoutes := FindMatchingRouteForDeparture(destination, origin, time)
		c.IndentedJSON(http.StatusOK, busRoutes)
	} else {
		c.IndentedJSON(http.StatusBadRequest, "Invalid time type parameter in request")
	}

}
