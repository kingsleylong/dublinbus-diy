
import 'bus_stop.dart';
import 'shape.dart';

class BusRoute {
  String routeNumber;
  List<BusStop> stops;
  List<Shape> shapes;
  Fares fares;
  TravelTimes? travelTimes;

  BusRoute(
      this.routeNumber, this.stops, this.shapes, this.fares, this.travelTimes);

  factory BusRoute.fromJson(Map<String, dynamic> json) {
    List<dynamic> stopsJson = json['stops'] as List;
    List<dynamic> shapesJson = json['shapes'] as List;
    return BusRoute(
      json['route_num'],
      stopsJson.map((busStopJson) => BusStop.fromJson(busStopJson, null)).toList(),
      shapesJson.map((shapeJson) => Shape.fromJson(shapeJson)).toList(),
      Fares.fromJson(json['fares'] ?? {}),
      TravelTimes.fromJson(json['travel_time'] ?? {}),
    );
  }
}

class Fares {
  double? adultLeap;
  double? adultCash;
  double? studentLeap;
  double? childLeap;
  double? childCash;

  Fares(this.adultLeap, this.adultCash, this.studentLeap, this.childLeap, this.childCash);

  factory Fares.fromJson(Map<String, dynamic> json) {
    return Fares(
      json['adult_leap'],
      json['adult_cash'],
      json['student_leap'],
      json['child_leap'],
      json['child_cash'],
    );
  }
}

class TravelTimes {
  int? transitTime;
  int? transitTimeMin;
  int? transitTimeMax;

  TravelTimes(this.transitTime, this.transitTimeMin, this.transitTimeMax);

  factory TravelTimes.fromJson(Map<String, dynamic> json) {
    return TravelTimes(
      json['transit_time'],
      json['transit_time_plus_mae'],
      json['transit_time_minus_mae'],
    );
  }
}
