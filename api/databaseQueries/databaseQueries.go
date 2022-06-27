package databaseQueries

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type busStop struct {
	stopId     string `json:"stop_id"`
	stopName   string `json:"stop_name"`
	stopNumber string `json:"stop_number"`
	stopLat    string `json:"stop_lat"`
	stopLon    string `json:"stop_lon"`
}

// GetDatabases returns one collection from the DB that has all routes with
// their ids and service route names
func GetDatabases(c *gin.Context) {

	// Create connection to mongo server and log any resulting error
	client, err := mongo.NewClient(options.Client().
		ApplyURI("mongodb://<username>:<password>@<host>:<port>/?retryWrites=true&w=majority"))
	if err != nil {
		log.Print(err)
	}

	// Create context variable and assign time for timeout
	// Log any resulting error here also
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Print(err)
	}
	defer client.Disconnect(ctx) // defer has rest of function complete before this disconnect

	databases, err := client.ListDatabases(ctx, bson.D{})

	c.IndentedJSON(http.StatusOK, databases)
}

func GetBusStop(c *gin.Context) {

	stopNum := c.Param("stopNum")
	stopNumString := "stop " + string(stopNum)

	// Create connection to mongo server and log any resulting error
	client, err := mongo.NewClient(options.Client().
		ApplyURI("mongodb://<username>:<password>@<host>:<port>/?retryWrites=true&w=majority"))
	if err != nil {
		log.Print(err)
	}

	// Create context variable and assign time for timeout
	// Log any resulting error here also
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Print(err)
	}
	defer client.Disconnect(ctx) // defer has rest of function complete before this disconnect

	var result bson.M

	dbPointer := client.Database("BusData")
	collectionPointer := dbPointer.Collection("BusStops")

	// Find one document that matches criteria and decode results into result address
	err = collectionPointer.FindOne(ctx, bson.D{{"stop_number", string(stopNumString)}}).
		Decode(&result)

	if err != nil {
		// Specific condition for valid query without matching documents
		if err == mongo.ErrNoDocuments {
			log.Print("No matching document found")
			return
		}
		log.Print(err)
	}

	// Return result as JSON along with code 200
	c.IndentedJSON(http.StatusOK, result)
}
