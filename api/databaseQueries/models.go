package databaseQueries

// busRoute is a type that is designed to read from the trips_new collection
// from MongoDB. It contains the id fields that combine to form a unique key
// for each entry (i.e. the route id, the shape id, the direction id, the trip id
// and the stop id). It also includes coordinates for the stop associated with this
// object as well as information on the route and the shape string used to
// draw the shape on the map. All fields map to type string from the database
type busRoute struct {
	Id     string    `bson:"_id" json:"_id"`
	Stops  []BusStop `bson:"stops" json:"stops"`
	Shapes []Shape   `bson:"shapes" json:"shapes"`
}

type busRouteJSON struct {
	ID     string                `bson:"_id" json:"_id"`
	Stops  []StopWithCoordinates `bson:"stops" json:"stops"`
	Shapes []Shape               `bson:"shapes" json:"shapes"`
}

// routeStop represents the stop information contained within the trips_n_stops
// collection in MongoDB. The information contains the StopId that can be used
// to identify each stop uniquely, the name of that stop, the stop number used
// by consumers of the Dublin Bus service, the coordinates of
// the stop that can be used to mark the stop on the map and finally a sequence
// number that can be used to sort the stops to ensure that they are in the
// correct order on a given route. All fields are returned as strings from the
// database
type routeStop struct {
	StopId       string `bson:"stop_id" json:"stop_id"`
	StopName     string `bson:"stop_name" json:"stop_name"`
	StopNumber   string `bson:"stop_number" json:"stop_number"`
	StopLat      string `bson:"stop_lat" json:"stop_lat"`
	StopLon      string `bson:"stop_lon" json:"stop_lon"`
	StopSequence string `bson:"stop_sequence" json:"stop_sequence"`
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

type BusStop struct {
	StopId     string `bson:"stop_id,omitempty" json:"stop_id,omitempty"`
	StopName   string `bson:"stop_name" json:"stop_name"`
	StopNumber string `bson:"stop_number" json:"stop_number"`
	StopLat    string `bson:"stop_lat" json:"stop_lat"`
	StopLon    string `bson:"stop_lon" json:"stop_lon"`
}

type StopWithCoordinates struct {
	StopID     string  `bson:"stop_id,omitempty" json:"stop_id,omitempty"`
	StopName   string  `bson:"stop_name" json:"stop_name"`
	StopNumber string  `bson:"stop_number" json:"stop_number"`
	StopLat    float64 `bson:"stop_lat" json:"stop_lat"`
	StopLon    float64 `bson:"stop_lon" json:"stop_lon"`
}

type findByAddressResponse struct {
	Matched []StopWithCoordinates `bson:"matched" json:"matched"`
	Nearby  []StopWithCoordinates `bson:"nearby" json:"nearby"`
}
