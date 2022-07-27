package databaseQueries

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func GetTravelTimePredictionTest(c *gin.Context) {

	resp, err := http.Get("http://localhost:5000/prediction/145/1/3/12/4/64800/date")
	if err != nil {
		log.Print(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	bodyString := string(body)
	c.IndentedJSON(http.StatusOK, bodyString)
}

//func GetTravelTimePrediction() []float64 {
//
//}
