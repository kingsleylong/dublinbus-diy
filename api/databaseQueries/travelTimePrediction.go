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

// GetTravelTimePrediction takes in the route number as a string, the
// date for prediction as a string in the format 'yyyy-MM-dd hh:mm:ss'
// (including the whitespace) and the direction of travel as a string and
// then returns the travel time prediction with two other values adjusted
// for the mean absolute error within the TravelTimePredictionFloat model
// as well as an error to be checked when generating travel time predictions
func GetTravelTimePrediction(routeNum string,
	date string,
	direction string) (TravelTimePredictionFloat, error) {

	// Features for prediction separated out from date here into an
	// array of strings
	features := FeatureExtraction(date)

	// URL is encoded here to prevent there being an issue with
	// whitespace in the path with some error checks also present
	baseUrl, err := url.Parse("https://dublinbus-diy.site/ml/prediction/")
	if err != nil {
		log.Println("Url Issue: ")
		log.Println(err.Error())
	}
	baseUrl.Path += strings.ToUpper(routeNum) + "/" + direction + "/" + features[0] + "/" +
		features[1] + "/" + features[2] + "/" + features[3] + "/" + date
	log.Println(baseUrl.String())
	resp, err := http.
		Get(baseUrl.String())
	if err != nil {
		log.Println("Error in the GET request")
		log.Print(err)
		return TravelTimePredictionFloat{0, 0, 0}, err
	}

	// Response is read in and stored in an object here before transformation
	// into a string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error in the read all call on the response body")
		log.Print(err)
		return TravelTimePredictionFloat{0, 0, 0}, err
	}

	bodyString := string(body)

	// String manipulation used here to have prediction values in correct format
	// to turn into floating point numbers
	bodyStringAdjusted := strings.Replace(bodyString, "[", "", 1)
	bodyStringAdjusted = strings.Replace(bodyStringAdjusted, "]\n", "", 1)
	bodyStrings := strings.Split(bodyStringAdjusted, ",")

	// Final check that travel time object was created correctly
	if len(bodyStrings) <= 1 {
		return TravelTimePredictionFloat{0, 0, 0}, errors.
			New("travel time prediction could not be generated")
	}

	// Travel time values turned into floats and then returned with nil error
	var travelTime TravelTimePredictionFloat

	travelTime.TransitTime, _ = strconv.ParseFloat(bodyStrings[0], 64)
	travelTime.TransitTimePlusMAE, _ = strconv.ParseFloat(bodyStrings[1], 64)
	travelTime.TransitTimeMinusMAE, _ = strconv.ParseFloat(bodyStrings[2], 64)

	return travelTime, nil
}

// FeatureExtraction is a function that uses string manipulation to take
// in the date parameter for the travel time query and then extracts
// the necessary predictive features for the predictive models and returns
// them all in an array of strings
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

// DayOfTheWeek is a function that takes in the slice of strings
// separating the date and the slice of strings separating the time
// that was created within the FeatureExtraction function and then
// uses the built-in time package to determine the day of the week
// of a given date and return a number from 0-6 inclusive (0 being
// Sunday). This number is returned as a string to make it suitable
// for return from FeatureExtraction and for use in the url path for
// creating travel time predictions
func DayOfTheWeek(dateSlice []string, timeSlice []string) string {

	// Individual fields from each portion of the date and time
	// are initialised as integers and then used to create a day
	// of the week as a string using built-in functions from the
	// time package
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

	// Switch statement used to match return value to given day of the week
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

// SecondsExtraction takes in an array of strings representing the time
// of day in the format 'hh:mm:ss' and then returns a string representation
// of the total number of seconds contained in each portion of that time
// added together to make it suitable for a feature within the travel time
// prediction model
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

// AdjustTravelTime is a function that takes in the intial values
// of the travel time prediction in the TravelTimePredictionFloat format
// as well as string representations of different arrival times for stops
// along the route and from there returns the full travel time object in
// the form of the TravelTimePrediction model
func AdjustTravelTime(initialTime TravelTimePredictionFloat,
	originArrivalTime string,
	destinationArrivalTime string,
	firstStopArrivalTime string,
	finalStopArrivalTime string) TravelTimePrediction {

	// Turn prediction values into seconds
	initialPredictionAsSeconds := initialTime.TransitTime * 60
	initialHighPredictionAsSeconds := initialTime.TransitTimePlusMAE * 60
	initialLowPredictionAsSeconds := initialTime.TransitTimeMinusMAE * 60

	// Convert string time representations into seconds
	originSeconds := convertStringTimeToTotalSeconds(originArrivalTime)
	destinationSeconds := convertStringTimeToTotalSeconds(destinationArrivalTime)
	firstStopSeconds := convertStringTimeToTotalSeconds(firstStopArrivalTime)
	finalStopSeconds := convertStringTimeToTotalSeconds(finalStopArrivalTime)

	// Get number of seconds for trip traversal in static timetable
	originToDestinationSeconds := destinationSeconds - originSeconds

	// Get seconds for first to last stop on whole route so that the specified
	// trip can be represented as a percentage of that route
	fullTripSeconds := finalStopSeconds - firstStopSeconds
	staticTripPercentageAsDecimal := originToDestinationSeconds / fullTripSeconds

	// Values for the travel time are taken in as integers
	journeyPrediction := int(math.Round(initialPredictionAsSeconds * staticTripPercentageAsDecimal))
	journeyHighPrediction := int(math.Round(initialHighPredictionAsSeconds * staticTripPercentageAsDecimal))
	journeyLowPrediction := int(math.Round(initialLowPredictionAsSeconds * staticTripPercentageAsDecimal))

	// Integer values of seconds for travel converted to minute values
	journeyPredictionInMins := journeyPrediction / 60
	journeyHighPredictionInMins := journeyHighPrediction / 60
	journeyLowPredictionInMins := journeyLowPrediction / 60

	// Destination times of arrival in string representation are created using
	// the number of seconds at the origin and the prediction integers generated
	destinationTime := createTimePredictionString(originSeconds, journeyPrediction)
	destinationHighTime := createTimePredictionString(originSeconds, journeyHighPrediction)
	destinationLowTime := createTimePredictionString(originSeconds, journeyLowPrediction)

	// Prediction values and time representations stored in TravelTimePrediction object
	var transitTimePredictions TravelTimePrediction

	transitTimePredictions.TransitTime = journeyPredictionInMins
	transitTimePredictions.TransitTimePlusMAE = journeyHighPredictionInMins
	transitTimePredictions.TransitTimeMinusMAE = journeyLowPredictionInMins
	transitTimePredictions.EstimatedArrivalTime = destinationTime
	transitTimePredictions.EstimatedArrivalHighTime = destinationHighTime
	transitTimePredictions.EstimatedArrivalLowTime = destinationLowTime

	// If prediction values are still 0 then travel time prediction from model
	// was unsuccessful so source set as static to enable generating predictions
	// from static information
	if transitTimePredictions.TransitTime == 0 && transitTimePredictions.TransitTimePlusMAE == 0 &&
		transitTimePredictions.TransitTimeMinusMAE == 0 {
		transitTimePredictions.Source = "static"
	} else {
		transitTimePredictions.Source = "prediction"
	}

	return transitTimePredictions
}

// GetStaticTime is a function that takes in the origin and destination stop
// times from the static timetable and calculates the difference between them
// to provide the user with some estimation as to the duration of the bus trip
func GetStaticTime(originStopArrivalTime string, destinationStopArrivalTime string) int {

	originSeconds := convertStringTimeToTotalSeconds(originStopArrivalTime)
	destinationSeconds := convertStringTimeToTotalSeconds(destinationStopArrivalTime)

	originToDestinationSeconds := destinationSeconds - originSeconds

	originToDestinationMinutes := int(math.Round(originToDestinationSeconds / 60))

	return originToDestinationMinutes
}

// convertStringTimeToTotalSeconds takes in a string representation
// of time in the format "yyyy-MM-dd hh:mm:ss" and returns a floating
// point number for the total number of seconds of the hours, minutes
// and seconds represented in the second half of the time string
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

// createTimePredictionString takes in the number of seconds for the origin
// and then integer value of the journey prediction in minutes and then returns
// a string representation of a time of day for arrival at a given destination
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
