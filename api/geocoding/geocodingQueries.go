package geocoding

import (
	"context"
	"github.com/gin-gonic/gin"
	"googlemaps.github.io/maps"
	"log"
	"net/http"
	"os"
	"time"
)

var DublinMapBoundsNE maps.LatLng
var DublinMapBoundsSW maps.LatLng
var DublinMapBounds maps.LatLngBounds

func GetCoordinates(c *gin.Context) {

	address := c.Param("address")

	DublinMapBoundsNE.Lat = 53.49337
	DublinMapBoundsNE.Lng = -6.05788

	DublinMapBoundsSW.Lng = -6.56495
	DublinMapBoundsSW.Lat = 53.14860

	DublinMapBounds.SouthWest = DublinMapBoundsSW
	DublinMapBounds.NorthEast = DublinMapBoundsNE

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	client, err := maps.NewClient(maps.WithAPIKey(os.Getenv("MAPS_API_KEY")))
	if err != nil {
		log.Print(err)
	}

	geo := &maps.GeocodingRequest{Address: address, Bounds: &DublinMapBounds, Region: "ie"}

	result, _ := client.Geocode(ctx, geo)

	c.IndentedJSON(http.StatusOK, result)
}
