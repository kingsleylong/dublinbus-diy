package geocoding

import (
	"googlemaps.github.io/maps"
	"log"
)

func getCoordinates() {

	c, err := maps.NewClient(maps.WithAPIKey("Insert-API-Key-Here"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}
