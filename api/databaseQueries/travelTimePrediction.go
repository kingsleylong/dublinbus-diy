package databaseQueries

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetTravelTimePredictionTest(c *gin.Context) {

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

	resp, err := http.
		Get(fmt.
			Sprintf("http://ec2-34-239-115-43.compute-1.amazonaws.com/prediction/%s/%s/%s/%s/%s/%s/%s",
				routeNum, direction, dayOfWeek, hour, month, departureTime, date))
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

	var travelTime TravelTimePrediction

	travelTime.TransitTime, _ = strconv.ParseFloat(bodyStrings[0], 64)
	travelTime.TransitTimePlusMAE, _ = strconv.ParseFloat(bodyStrings[1], 64)
	travelTime.TransitTimeMinusMAE, _ = strconv.ParseFloat(bodyStrings[2], 64)

	return travelTime
}

func DayOfTheWeek(date string) string {

	//2022-07-28 14:00:00

	dateTimeSplit := strings.Split(date, " ")
	dateSplit := strings.Split(dateTimeSplit[0], "-")
	timeSplit := strings.Split(dateTimeSplit[1], ":")

	year, _ := strconv.ParseInt(dateSplit[0], 10, 64)
	month, _ := strconv.ParseInt(dateSplit[1], 10, 64)
	day, _ := strconv.ParseInt(dateSplit[2], 10, 64)

	hour, _ := strconv.ParseInt(timeSplit[0], 10, 64)
	minute, _ := strconv.ParseInt(timeSplit[1], 10, 64)
	second, _ := strconv.ParseInt(timeSplit[2], 10, 64)

	dayOfWeek := time.Date(int(year),
		time.Month(month),
		int(day),
		int(hour),
		int(minute), int(second), 0, time.Local).Weekday().String()

	var dayNum string
	switch dayOfWeek {
	case "Sunday":
		dayNum = "0"
	case "Monday":
		dayNum = "1"
	case "Tuesday":
		dayNum = "2"
	case "Wednesday":
		dayNum = "3"
	case "Thursday":
		dayNum = "4"
	case "Friday":
		dayNum = "5"
	case "Saturday":
		dayNum = "6"
	}

	return dayNum
}

func DateExtraction(date string) (string, string) {

}

func SecondsExtraction(date string) string {

}