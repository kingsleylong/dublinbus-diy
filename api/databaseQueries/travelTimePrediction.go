package databaseQueries

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetTravelTimePredictionTest(c *gin.Context) {

	resp, err := http.
		Get("http://ec2-34-239-115-43.compute-1.amazonaws.com/prediction/116/1/3/12/4/64800/2022-07-30 14:00:00")
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

	var result TravelTimePredictionFloat

	result.TransitTime, _ = strconv.ParseFloat(bodyStrings[0], 64)
	result.TransitTimePlusMAE, _ = strconv.ParseFloat(bodyStrings[1], 64)
	result.TransitTimeMinusMAE, _ = strconv.ParseFloat(bodyStrings[2], 64)

	c.IndentedJSON(http.StatusOK, result)
}

func GetTravelTimePrediction(routeNum string,
	date string,
	direction string) (TravelTimePredictionFloat, error) {

	features := FeatureExtraction(date)

	resp, err := http.
		Get(fmt.
			Sprintf("http://ec2-34-239-115-43.compute-1.amazonaws.com/prediction/%s/%s/%s/%s/%s/%s/%s",
				routeNum, direction, features[0], features[1], features[2], features[3], date))
	if err != nil {
		log.Println("Error in the GET request")
		log.Print(err)
		return TravelTimePredictionFloat{0, 0, 0}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error in the read all call on the response body")
		log.Print(err)
		return TravelTimePredictionFloat{0, 0, 0}, err
	}

	bodyString := string(body)
	bodyStringAdjusted := strings.Replace(bodyString, "[", "", 1)
	bodyStringAdjusted = strings.Replace(bodyStringAdjusted, "]\n", "", 1)
	bodyStrings := strings.Split(bodyStringAdjusted, ",")

	if len(bodyStrings) <= 1 {
		return TravelTimePredictionFloat{0, 0, 0}, errors.
			New("travel time prediction could not be generated")
	}
	var travelTime TravelTimePredictionFloat

	log.Println(bodyStrings)

	travelTime.TransitTime, _ = strconv.ParseFloat(bodyStrings[0], 64)
	travelTime.TransitTimePlusMAE, _ = strconv.ParseFloat(bodyStrings[1], 64)
	travelTime.TransitTimeMinusMAE, _ = strconv.ParseFloat(bodyStrings[2], 64)

	return travelTime, nil
}

func FeatureExtraction(date string) []string {

	dateTimeSplit := strings.Split(date, " ")
	dateSplit := strings.Split(dateTimeSplit[0], "-")
	timeSplit := strings.Split(dateTimeSplit[1], ":")

	dayOfWeek := DayOfTheWeek(dateSplit, timeSplit)
	hour := timeSplit[0]
	month := dateSplit[1]
	seconds := SecondsExtraction(timeSplit)

	featureSlice := []string{dayOfWeek, hour, month, seconds}
	return featureSlice
}

func DayOfTheWeek(dateSlice []string, timeSlice []string) string {

	//2022-07-28 14:00:00
	year, _ := strconv.ParseInt(dateSlice[0], 10, 64)
	month, _ := strconv.ParseInt(dateSlice[1], 10, 64)
	day, _ := strconv.ParseInt(dateSlice[2], 10, 64)

	hour, _ := strconv.ParseInt(timeSlice[0], 10, 64)
	minute, _ := strconv.ParseInt(timeSlice[1], 10, 64)
	second, _ := strconv.ParseInt(timeSlice[2], 10, 64)

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

func SecondsExtraction(time []string) string {

	hoursInt, _ := strconv.ParseInt(time[0], 10, 64)
	minutesInt, _ := strconv.ParseInt(time[1], 10, 64)
	secondsInt, _ := strconv.ParseInt(time[2], 10, 64)

	hoursInSeconds := hoursInt * 3600
	minutesInSeconds := minutesInt * 60

	totalSeconds := hoursInSeconds + minutesInSeconds + secondsInt

	secondsValue := totalSeconds % (24 * 3600)

	return strconv.FormatInt(secondsValue, 10)
}

func AdjustTravelTime(initialTime TravelTimePredictionFloat,
	originArrivalTime string,
	destinationArrivalTime string,
	finalStopArrivalTime string) TravelTimePrediction {

	// Turn prediction values into seconds
	initialPredictionAsSeconds := initialTime.TransitTime * 60
	initialHighPredictionAsSeconds := initialTime.TransitTimePlusMAE * 60
	initialLowPredictionAsSeconds := initialTime.TransitTimeMinusMAE * 60
	log.Println("Initial Predictions as seconds:")
	log.Println(initialPredictionAsSeconds, initialHighPredictionAsSeconds, initialLowPredictionAsSeconds)
	log.Println("")

	originArrivalStringArray := strings.Split(originArrivalTime, ":")
	destinationArrivalStringArray := strings.Split(destinationArrivalTime, ":")
	finalStopArrivalStringArray := strings.Split(finalStopArrivalTime, ":")

	originArrivalHoursAsInt, _ := strconv.ParseInt(originArrivalStringArray[0], 10, 64)
	originArrivalMinutesAsInt, _ := strconv.ParseInt(originArrivalStringArray[1], 10, 64)
	originArrivalSecondsAsInt, _ := strconv.ParseInt(originArrivalStringArray[2], 10, 64)

	destinationArrivalHoursAsInt, _ := strconv.ParseInt(destinationArrivalStringArray[0], 10, 64)
	destinationArrivalMinutesAsInt, _ := strconv.ParseInt(destinationArrivalStringArray[1], 10, 64)
	destinationArrivalSecondsAsInt, _ := strconv.ParseInt(destinationArrivalStringArray[2], 10, 64)

	finalStopArrivalHoursAsInt, _ := strconv.ParseInt(finalStopArrivalStringArray[0], 10, 64)
	finalStopArrivalMinutesAsInt, _ := strconv.ParseInt(finalStopArrivalStringArray[1], 10, 64)
	finalStopArrivalSecondsAsInt, _ := strconv.ParseInt(finalStopArrivalStringArray[2], 10, 64)

	originHoursAsSeconds := originArrivalHoursAsInt * 3600
	originMinutesAsSeconds := originArrivalMinutesAsInt * 60
	originTotalSeconds := originHoursAsSeconds + originMinutesAsSeconds + originArrivalSecondsAsInt
	log.Println("Origin times as seconds with total:")
	log.Println(originHoursAsSeconds, originMinutesAsSeconds, originArrivalSecondsAsInt, originTotalSeconds)
	log.Println("")

	destinationHoursAsSeconds := destinationArrivalHoursAsInt * 3600
	destinationMinutesAsSeconds := destinationArrivalMinutesAsInt * 60
	destinationTotalSeconds := destinationHoursAsSeconds + destinationMinutesAsSeconds + destinationArrivalSecondsAsInt
	log.Println("Destination times as seconds with total:")
	log.Println(destinationHoursAsSeconds, destinationMinutesAsSeconds,
		destinationArrivalSecondsAsInt, destinationTotalSeconds)
	log.Println("")

	finalStopHoursAsSeconds := finalStopArrivalHoursAsInt * 3600
	finalStopMinutesAsSeconds := finalStopArrivalMinutesAsInt * 60
	finalStopTotalSeconds := finalStopHoursAsSeconds + finalStopMinutesAsSeconds + finalStopArrivalSecondsAsInt

	originToDestinationSeconds := float64(destinationTotalSeconds - originTotalSeconds)
	log.Println("Origin to destination seconds:")
	log.Println(originToDestinationSeconds)
	log.Println("")

	fullTripSeconds := float64(finalStopTotalSeconds - originTotalSeconds)
	staticTripPercentageAsDecimal := originToDestinationSeconds / fullTripSeconds

	//predictedTimeComplement := initialPredictionAsSeconds - originToDestinationSeconds
	//predictedHighTimeComplement := initialHighPredictionAsSeconds - originToDestinationSeconds
	//predictedLowTimeComplement := initialLowPredictionAsSeconds - originToDestinationSeconds
	//log.Println("Prediction complements:")
	//log.Println(predictedTimeComplement, predictedHighTimeComplement, predictedLowTimeComplement)
	//log.Println("")

	journeyPrediction := int(math.Round(initialPredictionAsSeconds * staticTripPercentageAsDecimal))
	journeyHighPrediction := int(math.Round(initialHighPredictionAsSeconds * staticTripPercentageAsDecimal))
	journeyLowPrediction := int(math.Round(initialLowPredictionAsSeconds * staticTripPercentageAsDecimal))
	log.Println("Journey prediction in seconds:")
	log.Println(journeyPrediction, journeyHighPrediction, journeyLowPrediction)
	log.Println("")

	journeyPredictionInMins := journeyPrediction / 60
	journeyHighPredictionInMins := journeyHighPrediction / 60
	journeyLowPredictionInMins := journeyLowPrediction / 60
	log.Println("Journey prediction in minutes:")
	log.Println(journeyPredictionInMins, journeyHighPredictionInMins, journeyLowPredictionInMins)
	log.Println("")

	var transitTimePredictions TravelTimePrediction

	transitTimePredictions.TransitTime = journeyPredictionInMins
	transitTimePredictions.TransitTimePlusMAE = journeyHighPredictionInMins
	transitTimePredictions.TransitTimeMinusMAE = journeyLowPredictionInMins

	return transitTimePredictions
}
