package databaseQueries

import (
	"errors"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func GetTravelTimePrediction(routeNum string,
	date string,
	direction string) (TravelTimePredictionFloat, error) {

	features := FeatureExtraction(date)

	baseUrl, err := url.Parse("http://3.250.172.35/prediction/")
	if err != nil {
		log.Println("Url Issue: ")
		log.Println(err.Error())
	}
	baseUrl.Path += strings.ToUpper(routeNum) + "/" + direction + "/" + features[0] + "/" +
		features[1] + "/" + features[2] + "/" + features[3] + "/" + date

	resp, err := http.
		Get(baseUrl.String())
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
	firstStopArrivalTime string,
	finalStopArrivalTime string) TravelTimePrediction {

	// Turn prediction values into seconds
	initialPredictionAsSeconds := initialTime.TransitTime * 60
	initialHighPredictionAsSeconds := initialTime.TransitTimePlusMAE * 60
	initialLowPredictionAsSeconds := initialTime.TransitTimeMinusMAE * 60

	originSeconds := convertStringTimeToTotalSeconds(originArrivalTime)
	destinationSeconds := convertStringTimeToTotalSeconds(destinationArrivalTime)
	firstStopSeconds := convertStringTimeToTotalSeconds(firstStopArrivalTime)
	finalStopSeconds := convertStringTimeToTotalSeconds(finalStopArrivalTime)

	originToDestinationSeconds := destinationSeconds - originSeconds

	fullTripSeconds := finalStopSeconds - firstStopSeconds
	staticTripPercentageAsDecimal := originToDestinationSeconds / fullTripSeconds

	journeyPrediction := int(math.Round(initialPredictionAsSeconds * staticTripPercentageAsDecimal))
	journeyHighPrediction := int(math.Round(initialHighPredictionAsSeconds * staticTripPercentageAsDecimal))
	journeyLowPrediction := int(math.Round(initialLowPredictionAsSeconds * staticTripPercentageAsDecimal))

	journeyPredictionInMins := journeyPrediction / 60
	journeyHighPredictionInMins := journeyHighPrediction / 60
	journeyLowPredictionInMins := journeyLowPrediction / 60

	destinationTime := createTimePredictionString(originSeconds, journeyPrediction)
	destinationHighTime := createTimePredictionString(originSeconds, journeyHighPrediction)
	destinationLowTime := createTimePredictionString(originSeconds, journeyLowPrediction)

	var transitTimePredictions TravelTimePrediction

	transitTimePredictions.TransitTime = journeyPredictionInMins
	transitTimePredictions.TransitTimePlusMAE = journeyHighPredictionInMins
	transitTimePredictions.TransitTimeMinusMAE = journeyLowPredictionInMins
	transitTimePredictions.EstimatedArrivalTime = destinationTime
	transitTimePredictions.EstimatedArrivalHighTime = destinationHighTime
	transitTimePredictions.EstimatedArrivalLowTime = destinationLowTime

	if transitTimePredictions.TransitTime == 0 && transitTimePredictions.TransitTimePlusMAE == 0 &&
		transitTimePredictions.TransitTimeMinusMAE == 0 {
		transitTimePredictions.Source = "static"
	} else {
		transitTimePredictions.Source = "prediction"
	}

	return transitTimePredictions
}

func GetStaticTime(originStopArrivalTime string, destinationStopArrivalTime string) int {

	originSeconds := convertStringTimeToTotalSeconds(originStopArrivalTime)
	destinationSeconds := convertStringTimeToTotalSeconds(destinationStopArrivalTime)

	originToDestinationSeconds := destinationSeconds - originSeconds

	originToDestinationMinutes := int(math.Round(originToDestinationSeconds / 60))

	return originToDestinationMinutes
}

func convertStringTimeToTotalSeconds(time string) float64 {

	dateAndTimeSlice := strings.Split(time, ":")
	hoursAsFloat, _ := strconv.ParseFloat(dateAndTimeSlice[0], 64)
	minutesAsFloat, _ := strconv.ParseFloat(dateAndTimeSlice[1], 64)
	secondsAsFloat, _ := strconv.ParseFloat(dateAndTimeSlice[2], 64)

	hoursAsSeconds := hoursAsFloat * 3600
	minutesAsSeconds := minutesAsFloat * 60
	totalSeconds := hoursAsSeconds + minutesAsSeconds + secondsAsFloat

	return totalSeconds
}

func createTimePredictionString(seconds float64, journeyPrediction int) string {

	timeInSeconds := int(math.Round(seconds)) + journeyPrediction
	hours := timeInSeconds / 3600
	minutes := (timeInSeconds % 3600) / 60
	hoursString := strconv.Itoa(hours)
	minutesString := strconv.Itoa(minutes)
	if hours < 10 {
		hoursString = "0" + hoursString
	}
	if minutes < 10 {
		minutesString = "0" + minutesString
	}
	destinationTime := hoursString + ":" + minutesString

	return destinationTime
}
