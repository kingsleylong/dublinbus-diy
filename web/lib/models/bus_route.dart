import 'package:web/models/bus_stop.dart';
import 'package:web/models/shape.dart';

class BusRoute {
  String routeNumber;
  List<BusStop> stops;
  List<Shape> shapes;
  int? travelTime;

  BusRoute(this.routeNumber, this.stops, this.shapes, this.travelTime);

  factory BusRoute.fromJson(Map<String, dynamic> json) {
    List<dynamic> stopsJson = json['stops'] as List;
    List<dynamic> shapesJson = json['shapes'] as List;
    return BusRoute(
      json['route_num'],
      stopsJson.map((busStopJson) => BusStop.fromJson(busStopJson)).toList(),
      shapesJson.map((shapeJson) => Shape.fromJson(shapeJson)).toList(),
      json['travel_time']
    );
  }
}