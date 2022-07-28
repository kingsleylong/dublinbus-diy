package databaseQueries

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetTravelTimePredictionTest(c *gin.Context) {

	resp, err := http.
		Get("http://ec2-34-239-115-43.compute-1.amazonaws.com/prediction/102/1/3/12/4/64800/2022-07-28 14:00:00")
	if err != nil {
		log.Print(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	bodyString := string(body)
	bodyStringAdjusted := strings.Replace(bodyString, "[", "", 1)
	bodyStringAdjusted = strings.Replace(bodyStringAdjusted, "]\n", "", 1)
	bodyStrings := strings.Split(bodyStringAdjusted, ",")

	var result TravelTimePrediction

	result.TransitTime, _ = strconv.ParseFloat(bodyStrings[0], 64)
	result.TransitTimePlusMAE, _ = strconv.ParseFloat(bodyStrings[1], 64)
	result.TransitTimeMinusMAE, _ = strconv.ParseFloat(bodyStrings[2], 64)

	c.IndentedJSON(http.StatusOK, result)
}

func GetTravelTimePrediction(routeNum string,
	date string,
	direction string) TravelTimePrediction {

	dayOfWeek := DayOfTheWeek(date)
	hour, month := DateExtraction(date)
	departureTime := SecondsExtraction(date)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	bodyString := string(body)
	bodyStringAdjusted := strings.Replace(bodyString, "[", "", 1)
	bodyStringAdjusted = strings.Replace(bodyStringAdjusted, "]\n", "", 1)
	bodyStrings := strings.Split(bodyStringAdjusted, ",")

	var travelTime TravelTimePrediction

	travelTime.TransitTime, _ = strconv.ParseFloat(bodyStrings[0], 64)
	travelTime.TransitTimePlusMAE, _ = strconv.ParseFloat(bodyStrings[1], 64)
	travelTime.TransitTimeMinusMAE, _ = strconv.ParseFloat(bodyStrings[2], 64)

	return travelTime
}
