import 'dart:convert';

import 'package:dropdown_search/dropdown_search.dart';
import 'package:flutter/material.dart';
import 'package:geolocator/geolocator.dart';
import 'package:http/http.dart' as http;
import 'package:http/http.dart';
import 'package:localstorage/localstorage.dart';
import 'package:uuid/uuid.dart';

import '../api/location_service.dart';
import '../env.dart';
import 'bus_route.dart';
import 'bus_route_filter.dart';

/// The model for the search form
class SearchFormModel extends ChangeNotifier {
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();

  // Use a flag to control the visibility of the route options
  // https://stackoverflow.com/a/46126667
  bool visibilityRouteOptions = false;

  // Use a flag to control the visibility of the loading icon
  bool visibilityLoadingIcon = false;

  // The instance field that holds the state of origin dropdown list
  final _originSelectionKey = GlobalKey<DropdownSearchState<Prediction>>();
  PlaceDetail? originPlaceDetail;

  // The instance field that holds the state of destination dropdown list
  final _destinationSelectionKey = GlobalKey<DropdownSearchState<Prediction>>();
  PlaceDetail? destinationPlaceDetail;

  // The instance field that holds the state of the datetime picker
  final TextEditingController _dateTimePickerController =
      TextEditingController(text: DateTime.now().toString());

  // DO NOT use the late modifier because it won't initialize the variable, even as null, before
  // the request. Use the null safety to do that instead so we can take advantage of Conditional
  // property access feature of dart.
  // https://dart.dev/codelabs/null-safety
  List<BusRoute>? _busRoutes;

  List<Item>? _busRouteItems;

  // The instance field that holds the state of the time type toggle button
  int? timeTypeToggleIndex = 0;

  // The available options for the time type toggle button
  List<TimeType> timeTypes = [TimeType.departure, TimeType.arrival];

  // the state of local storage for Favorite routes
  bool storageInitialized = false;

  // the instance of local storage for Favorite routes
  final LocalStorage favoritesStorage = LocalStorage('fav_routes');

  // the favorite routes in the memory
  Map<String, RouteItem> favoriteRoutes = {};

  // getters
  TextEditingController get dateTimePickerController => _dateTimePickerController;

  get destinationSelectionKey => _destinationSelectionKey;

  get originSelectionKey => _originSelectionKey;

  GlobalKey<FormState> get formKey => _formKey;

  List<Item>? get busRouteItems => _busRouteItems;

  List<BusRoute>? get busRoutes => _busRoutes;

  // The sessionToken is very important for billing!!!
  // refer to https://developers.google.com/maps/documentation/places/web-service/autocomplete#sessiontoken
  late String sessionToken;

  SearchFormModel() {
    // initialize the token
    sessionToken = generateUuid();
  }

  String generateUuid() => const Uuid().v4();

  // Use the Place Autocomplete service to implement the address searching feature
  // https://developers.google.com/maps/documentation/places/web-service/autocomplete#place_autocomplete_requests
  Future<List<Prediction>> autocompleteAddress(String filter) async {
    if (filter.isEmpty) return [];

    print('Places autocomplete API sessionToken: $sessionToken');

    String googlemapApiHost = googleMapApiHost;
    Uri request = Uri.https(
      googlemapApiHost,
      '/api/googlemaps/maps/api/place/autocomplete/json',
      {
        'input': filter,
        'inputtype': 'textquery',
        'sessionToken': sessionToken,
        'region': 'ie', // Ireland
        'location': '53.34640516825308, -6.267271142573096', // Dublin City Center
        'radius': '35000' // 35 km
      },
    );
    print('request uri: $request');
    Response response = await http.get(request);
    if (response.statusCode == 200) {
      final Map<String, dynamic> responseBody = jsonDecode(response.body);
      // response structure: https://developers.google.com/maps/documentation/places/web-service/autocomplete#place_autocomplete_responses
      if (responseBody['status'] == 'OK') {
        List predictionsJson = responseBody['predictions'];
        List<Prediction> predictions =
            predictionsJson.map((predictJson) => Prediction.fromJson(predictJson)).toList();
        return predictions;
      }
    } else {
      print('Error fetch the predictions');
    }
    return [];
  }

  // Place Details requests
  // https://developers.google.com/maps/documentation/places/web-service/detailsk
  fetchPlaceDetails(Prediction? prediction, String type) async {
    print('Places autocomplete API sessionToken: $sessionToken');

    if (type == 'origin') {
      originPlaceDetail = null;
    } else if (type == 'destination') {
      destinationPlaceDetail = null;
    }

    if (prediction == null) return;

    // TODO delay the logic to when the submit button is clicked??
    if (prediction.placeId == 'here') {
      Position position = await determinePosition();
      print('get position: $position');
      var placeDetail = PlaceDetail(position.latitude, position.longitude);
      if (type == 'origin') {
        originPlaceDetail = placeDetail;
      } else if (type == 'destination') {
        destinationPlaceDetail = placeDetail;
      }
      return;
    }
    var placeId = prediction.placeId;
    String googlemapApiHost = googleMapApiHost;
    Uri request = Uri.https(
      googlemapApiHost,
      '/api/googlemaps/maps/api/place/details/json',
      {'place_id': placeId, 'fields': 'geometry', 'sessionToken': sessionToken},
    );
    print('fetchPlaceDetails request uri: $request');
    Response response = await http.get(request);
    if (response.statusCode == 200) {
      final Map<String, dynamic> responseJson = jsonDecode(response.body);
      // response structure: https://developers.google.com/maps/documentation/places/web-service/details#PlaceDetailsResponses
      if (responseJson['status'] == 'OK') {
        var placeDetail = PlaceDetail.fromJson(responseJson['result']);
        print('placeDetail: $placeDetail');
        if (type == 'origin') {
          originPlaceDetail = placeDetail;
        } else if (type == 'destination') {
          destinationPlaceDetail = placeDetail;
        }
      }
    } else {
      print('Error fetch the predictions');
    }

    // Important!! sessionToken must be regenerated before return to finish this search session!!
    sessionToken = generateUuid();
    print('reset sessionToken: $sessionToken');
  }

  Future<void> fetchBusRoute(BusRouteSearchFilter searchFilter) async {
    // start loading, hide route options
    visibilityRouteOptions = false;
    visibilityLoadingIcon = true;
    _busRoutes?.clear();
    _busRouteItems?.clear();
    // notify immediately because the below code will block execution until api returned
    notifyListeners();

    String pathParams = '/${searchFilter.originStopNumber}/${searchFilter.destinationStopNumber}'
        '/${searchFilter.timeType.name}/${searchFilter.time}';
    final response = await http.get(
      Uri.https(apiHost, "/api/route/matchingRoute$pathParams"),
      headers: {
        "Accept": "application/json",
      },
    );

    await initializeFavoritesStorage();

    // stop loading, show route options
    visibilityRouteOptions = true;
    visibilityLoadingIcon = false;
    notifyListeners();

    if (response.statusCode == 200) {
      // If the server did return a 200 OK response, then parse the JSON.
      List<BusRoute> busRouteList = [];
      final List? busRoutesJson = jsonDecode(response.body);

      if (busRoutesJson != null && busRoutesJson.isNotEmpty) {
        busRouteList =
            List.generate(busRoutesJson.length, (index) => BusRoute.fromJson(busRoutesJson[index]));
        _busRoutes = busRouteList;
        _busRouteItems = generateItems(busRouteList, favoriteRoutes);
      } else {
        _busRoutes = [];
        _busRouteItems = [];
      }

      notifyListeners();
      // return busRouteList;
    } else {
      // If the server did not return a 200 OK response, then throw an exception.
      throw Exception('Failed to load bus routes');
    }
  }

  List<Item> generateItems(List<BusRoute> data, Map<String, RouteItem> favoriteRouteList) {
    return List<Item>.generate(data.length, (int index) {
      return Item(
        headerValue: data[index].routeNumber,
        expandedValue:
            '${data[index].stops.length} stops. Starts from ${data[index].stops[0].stopName}.',
        expandedDetailsValue: data[index]
            .stops
            .map((stop) => '${stop.stopName} - stop ${stop.stopNumber}')
            .reduce((value, element) => '$value\n$element'),
        busRoute: data[index],
        favorite: favoriteRouteList[data[index].routeNumber]?.favourite ?? false,
      );
    });
  }

  // toggle the state of the selected route and update the route options state
  toggleFavorite(BusRoute route) {
    var favoriteRoute = favoriteRoutes[route.routeNumber];
    if (favoriteRoute == null) {
      favoriteRoutes[route.routeNumber] = RouteItem(favourite: true, route: route);
    } else {
      favoriteRoutes.remove(route.routeNumber);
    }

    updateRouteOptions();

    saveToStorage();
    notifyListeners();
  }

  // sync the favorite state to the route options
  updateRouteOptions() {
    if (_busRoutes != null) {
      _busRouteItems = generateItems(_busRoutes!, favoriteRoutes);
    }
  }

  // save the favorites to the local storage
  saveToStorage() {
    favoritesStorage.setItem('items', favoriteRoutes);
    notifyListeners();
  }

  // create the local storage and load the saved routes into memory
  initializeFavoritesStorage() async {
    await favoritesStorage.ready;

    if (!storageInitialized) {
      final Map<String, dynamic> favoriteRoutesJson = await favoritesStorage.getItem("items") ?? {};
      favoriteRoutesJson.forEach((key, value) {
        favoriteRoutes.putIfAbsent(key, () => RouteItem.fromJson(value));
      });

      storageInitialized = true;
      notifyListeners();
    }
  }
}

// The class that wraps the bus route and favorite status
class RouteItem {
  BusRoute route;

  // The state that control the icon of the Favorite button
  bool favourite;

  RouteItem({required this.route, required this.favourite});

  factory RouteItem.fromJson(Map<String, dynamic> json) {
    return RouteItem(route: BusRoute.fromJson(json['route']), favourite: json['favourite']);
  }

  // This is the default method used to serialize the object to JSON before it's
  // saved into the local storage
  toJson() {
    Map<dynamic, dynamic> m = {};

    m['route'] = route;
    m['favourite'] = favourite;
    return m;
  }
}

class Prediction {
  String placeId;
  String mainText;
  String secondaryText;
  late PlaceDetail placeDetail;

  // List<String> items;

  Prediction(this.placeId, this.mainText, this.secondaryText);

  // response structure
  // https://developers.google.com/maps/documentation/places/web-service/autocomplete#place_autocomplete_responses
  factory Prediction.fromJson(Map<String, dynamic> json) {
    return Prediction(
      json['place_id'],
      json['structured_formatting']['main_text'],
      json['structured_formatting']['secondary_text'],
      // json['items'],
    );
  }

  @override
  String toString() {
    return mainText;
  }
}

class PlaceDetail {
  double latitude;
  double longitude;

  PlaceDetail(this.latitude, this.longitude);

  // response structure
  // https://developers.google.com/maps/documentation/places/web-service/details#PlaceDetailsResponses
  factory PlaceDetail.fromJson(Map<String, dynamic> json) {
    var location = json['geometry']['location'];
    return PlaceDetail(
      (location['lat'] as num).toDouble(),
      (location['lng'] as num).toDouble(),
    );
  }

  @override
  String toString() {
    return '$latitude,$longitude';
  }
}

class Item {
  Item({
    required this.headerValue,
    required this.busRoute,
    this.isExpanded = false,
    required this.expandedValue,
    required this.expandedDetailsValue,
    required this.favorite,
  });

  String headerValue;
  bool isExpanded;
  BusRoute busRoute;
  String expandedValue;
  String expandedDetailsValue;
  bool favorite;
}
