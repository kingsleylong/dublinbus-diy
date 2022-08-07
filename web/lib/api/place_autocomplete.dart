import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:http/http.dart';
import 'package:uuid/uuid.dart';

// Initially I tried to use the google_maps_webservice package https://pub.dev/packages/google_maps_webservice
// However, at the time I was developing this feature, google api has a CORS restriction on all
// its web services. Same issue could be found at the package github repo. Then I have to set up
// a proxy server that delegates the google APIs and add the CORS header.
// https://github.com/lejard-h/google_maps_webservice/issues/70
Future<void> main() async {
  var list = await autocompleteAddress('ucd');
  print(list);
  await autocompleteAddress('ucd');
}

late String sessionToken;

Future<List<String>> autocompleteAddress(String filter) async {
  if (filter.isEmpty) return [];

  sessionToken = const Uuid().v4();
  print('Places autocomplete API sessionToken: $sessionToken');

  String googlemap_api_host = 'ipa-003.ucd.ie';
  Uri request = Uri.http(
    googlemap_api_host,
    '/api/googlemaps/maps/api/place/autocomplete/json',
    {'input': filter, 'inputtype': 'textquery', 'sessionToken': sessionToken},
  );
  print('request uri: $request');
  Response response = await http.get(request);
  if (response.statusCode == 200) {
    var body = response.body;
    // print(body);
    final Map<String, dynamic> predictionsBody = jsonDecode(response.body);
    print(predictionsBody);
    List predictionsJson = predictionsBody['predictions'];
    List<String> predictions = predictionsJson.map((e) => e['description'] as String).toList();
    return predictions;
  }

  return [];
}
