import 'bus_stop.dart';
import 'shape.dart';

class BusRoute {
  String routeNumber;
  List<BusStop> stops;
  List<Shape> shapes;
  Fares fares;
  TravelTimes travelTimes;

  BusRoute(this.routeNumber, this.stops, this.shapes, this.fares, this.travelTimes);

  factory BusRoute.fromJson(Map<String, dynamic> json) {
    List<dynamic> stopsJson = json['stops'] as List;
    List<dynamic> shapesJson = json['shapes'] as List;
    return BusRoute(
      json['route_num'].toString().toUpperCase(),
      stopsJson.map((busStopJson) => BusStop.fromJson(busStopJson, null)).toList(),
      shapesJson.map((shapeJson) => Shape.fromJson(shapeJson)).toList(),
      Fares.fromJson(json['fares']),
      TravelTimes.fromJson(json['travel_time']),
    );
  }

  toJson() {
    Map<dynamic, dynamic> m = {};

    m['route_num'] = routeNumber;
    m['stops'] = stops;
    m['shapes'] = shapes;
    m['fares'] = fares;
    m['travel_time'] = travelTimes;
    return m;
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
      // json[] may return an int and cause type error in iOS, need to handle carefully
      // https://stackoverflow.com/questions/71615935/typeerror-type-int-is-not-a-subtype-of-type-double-flutter
      (json['adult_leap'] as num).toDouble(),
      (json['adult_cash'] as num).toDouble(),
      (json['student_leap'] as num).toDouble(),
      (json['child_leap'] as num).toDouble(),
      (json['child_cash'] as num).toDouble(),
    );
  }

  toJson() {
    Map<dynamic, dynamic> m = {};

    m['adult_leap'] = adultLeap;
    m['adult_cash'] = adultCash;
    m['student_leap'] = studentLeap;
    m['child_leap'] = childLeap;
    m['child_cash'] = childCash;
    return m;
  }
}

class TravelTimes {
  TravelTimeSources source;
  int? transitTime;
  int? transitTimeMin;
  int? transitTimeMax;
  String scheduledDepartureTime;

  TravelTimes(this.source, this.transitTime, this.transitTimeMin, this.transitTimeMax,
      this.scheduledDepartureTime);

  factory TravelTimes.fromJson(Map<String, dynamic> json) {
    return TravelTimes(
      // lookup enum by name:
      // https://medium.com/dartlang/dart-2-15-7e7a598e508a
      TravelTimeSources.values.byName(json['source']),
      json['transit_time'],
      json['transit_time_minus_mae'],
      json['transit_time_plus_mae'],
      json['scheduled_departure_time'],
    );
  }

  toJson() {
    Map<dynamic, dynamic> m = {};

    m['source'] = source.name;
    m['transit_time'] = transitTime;
    m['transit_time_minus_mae'] = transitTimeMin;
    m['transit_time_plus_mae'] = transitTimeMax;
    m['scheduled_departure_time'] = scheduledDepartureTime;
    return m;
  }
}

enum TravelTimeSources { static, prediction }
