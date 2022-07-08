import 'dart:convert';

import 'package:web/models/bus_stop.dart';

class BusRoute {
  final String routeNumber;
  final List<BusStop> stops;

  const BusRoute({
    required this.routeNumber,
    required this.stops,
  });

  factory BusRoute.fromJson(Map<String, dynamic> json) {
    List<dynamic> stopsJson = json['route_stops'] as List;
    return BusRoute(
      routeNumber: json['route_num'],
      stops: stopsJson.map((busStopJson) => BusStop.fromJsonForRoute(busStopJson)).toList(),
    );
  }
}