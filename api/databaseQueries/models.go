package databaseQueries

// busRoute is a type that is designed to read from the trips_n_stops collection
// from MongoDB following passing the aggregation pipeline defined in the
// busRouteQueries file. The route short name is used as the id feature
// for each busRoute while this structure contains arrays of nested structures.
// The Stops array is made of type BusStop while the Shapes array is made of type
// Shape.
type busRoute struct {
	Id     string    `bson:"_id" json:"_id"`
	Stops  []BusStop `bson:"stops" json:"stops"`
	Shapes []Shape   `bson:"shapes" json:"shapes"`
}

// busRouteJSON is designed in a very similar fashion to the busRoute structure.
// The ID field mirrors that of the busRoute struct and the Shapes array is exactly
// the same also. The main difference between these structures is in the Stops array.
// In the busRouteJSON this array is made of type RouteStop which as a key difference
// returns the coordinates of each bus stop as type float as opposed to strings.
type busRouteJSON struct {
	RouteNum        string      `bson:"route_num" json:"route_num"`
	Stops           []RouteStop `bson:"stops" json:"stops"`
	Shapes          []ShapeJSON `bson:"shapes" json:"shapes"`
	FareCalculation float64     `bson:"fare_calculation" json:"fare_calculation"`
}

// RouteStop represents the stop information contained within the trips_n_stops
// collection in MongoDB. The information contains the StopId that can be used
// to identify each stop uniquely, the name of that stop, the stop number used
// by consumers of the Dublin Bus service, the coordinates of
// the stop that can be used to mark the stop on the map, a sequence
// number that can be used to sort the stops to ensure that they are in the
// correct order on a given route and finally arrival and departure times for
// when a bus arrived and departed that particular stop for a given trip.
// All fields are returned as strings from the database, apart from the
// coordinates of the stop that come back as a float each
type RouteStop struct {
	StopId            string  `bson:"stop_id" json:"stop_id"`
	StopName          string  `bson:"stop_name" json:"stop_name"`
	StopNumber        string  `bson:"stop_number" json:"stop_number"`
	StopLat           float64 `bson:"stop_lat" json:"stop_lat"`
	StopLon           float64 `bson:"stop_lon" json:"stop_lon"`
	StopSequence      string  `bson:"stop_sequence" json:"stop_sequence"`
	ArrivalTime       string  `bson:"arrival_time" json:"arrival_time"`
	DepartureTime     string  `bson:"departure_time" json:"departure_time"`
	DistanceTravelled float64 `bson:"shape_dist_traveled" json:"shape_dist_traveled"`
}

// route is a struct that contains a means of matching the route number (referred to
// as RouteShortName in this object) to the route id (i.e. the RouteId). All
// fields map to type string from the database
type route struct {
	RouteId        string `bson:"route_id" json:"route_id"`
	RouteShortName string `bson:"route_short_name" json:"route_short_name"`
}

// Shape is struct that contains the coordinates for each turn in a bus
// line as it travels its designated route that combined together allow
// the bus route to be drawn on a map matching the road network of Dublin.
// All fields map to type string from the database
type Shape struct {
	//ShapeId         string `bson:"shape_id" json:"shape_id"`
	ShapePtLat      string `bson:"shape_pt_lat" json:"shape_pt_lat"`
	ShapePtLon      string `bson:"shape_pt_lon" json:"shape_pt_lon"`
	ShapePtSequence string `bson:"shape_pt_sequence" json:"shape_pt_sequence"`
	ShapeDistTravel string `bson:"shape_dist_traveled" json:"shape_dist_traveled"`
}

type ShapeJSON struct {
	ShapePtLat      float64 `bson:"shape_pt_lat" json:"shape_pt_lat"`
	ShapePtLon      float64 `bson:"shape_pt_lon" json:"shape_pt_lon"`
	ShapePtSequence string  `bson:"shape_pt_sequence" json:"shape_pt_sequence"`
	ShapeDistTravel string  `bson:"shape_dist_travel" json:"shape_dist_travel"`
}

// BusStop contains all the necessary information from the mongo collection
// trips_n_stops to provide information on each stop on a given bus route
// for a certain trip, including the stop id, its name and number, its
// coordinates, its number in the sequence of bus stops on the trip and the
// arrival and departure times for the bus making the given trip. All fields
// are returned as strings from the database
type BusStop struct {
	StopId            string `bson:"stop_id,omitempty" json:"stop_id,omitempty"`
	StopName          string `bson:"stop_name" json:"stop_name"`
	StopNumber        string `bson:"stop_number" json:"stop_number"`
	StopLat           string `bson:"stop_lat" json:"stop_lat"`
	StopLon           string `bson:"stop_lon" json:"stop_lon"`
	StopSequence      string `bson:"stop_sequence" json:"stop_sequence"`
	ArrivalTime       string `bson:"arrival_time" json:"arrival_time"`
	DepartureTime     string `bson:"departure_time" json:"departure_time"`
	DistanceTravelled string `bson:"shape_dist_traveled" json:"shape_dist_traveled"`
}

// StopWithCoordinates contains the fields necessary to map out a route
// on a map by including identifying information for each stop (its id,
// name and number) as well as the coordinates for that stop as floating
// point numbers
type StopWithCoordinates struct {
	StopID     string  `bson:"stop_id,omitempty" json:"stop_id,omitempty"`
	StopName   string  `bson:"stop_name" json:"stop_name"`
	StopNumber string  `bson:"stop_number" json:"stop_number"`
	StopLat    float64 `bson:"stop_lat" json:"stop_lat"`
	StopLon    float64 `bson:"stop_lon" json:"stop_lon"`
}

// findByAddressResponse is a simple structure that just contains two arrays
// of nested structures - the StopWithCoordinates structure. It separates the
// response into two different arrays - an array of stops that were sourced
// by matching the search keyword in our database and then an array of
// stops that were sourced by finding the coordinates of the search keyword
// using Google Maps' geocoding service and then finding stops with nearby
// coordinates in our database
type findByAddressResponse struct {
	Matched []StopWithCoordinates `bson:"matched" json:"matched"`
	Nearby  []StopWithCoordinates `bson:"nearby" json:"nearby"`
}
