import 'dart:convert';

import 'env.dart';
import 'models/bus_stop.dart';
import 'package:http/http.dart' as http;

Future<List<BusStop>> fetchFutureBusStopsByName(String filter) async {
  List<BusStop> busStopList = [];

  String url = '$apiHost/api/stop/findByAddress';
  // final paramsStr = (filter == '') ? '' : '&filter=$filter';
  final paramsStr = (filter == '') ? '' : '/$filter';
  final response = await http.get(
    Uri.parse('$url$paramsStr'),
    headers: {
      "Accept": "application/json",
    },
  );

  if (response.statusCode == 200) {
    // If the server did return a 200 OK response, then parse the JSON.
    final Map<String, dynamic> busStopsJson = jsonDecode(response.body);
    final List? matchedStopsJson = busStopsJson['matched'];
    final List? nearbyStopsJson = busStopsJson['nearby'];

    List<BusStop> matchedStopList;
    List<BusStop> nearbyStopList;
    if (matchedStopsJson != null) {
      matchedStopList = List.generate(
          matchedStopsJson.length,
              (index) =>
              BusStop.fromJson(matchedStopsJson[index], BusStopType.matched));
      busStopList.addAll(matchedStopList);
    }
    if (nearbyStopsJson != null) {
      nearbyStopList = List.generate(
          nearbyStopsJson.length,
              (index) =>
              BusStop.fromJson(nearbyStopsJson[index], BusStopType.nearby));
      busStopList.addAll(nearbyStopList);
    }
    return busStopList;
  } else {
    // If the server did not return a 200 OK response, then throw an exception.
    throw Exception('Failed to load bus routes');
  }
}