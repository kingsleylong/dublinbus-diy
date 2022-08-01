import 'dart:convert';

import 'package:dropdown_search/dropdown_search.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

import '../env.dart';
import 'bus_route.dart';
import 'bus_route_filter.dart';
import 'bus_stop.dart';

/// The model for the search form
class SearchFormModel extends ChangeNotifier {
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();

  // Use a flag to control the visibility of the route options
  // https://stackoverflow.com/a/46126667
  bool visibilityRouteOptions = false;

  // Use a flag to control the visibility of the loading icon
  bool visibilityLoadingIcon = false;

  // The instance field that holds the state of origin dropdown list
  final _originSelectionKey = GlobalKey<DropdownSearchState<BusStop>>();

  // The instance field that holds the state of destination dropdown list
  final _destinationSelectionKey = GlobalKey<DropdownSearchState<BusStop>>();

  // The instance field that holds the state of the datetime picker
  final TextEditingController _dateTimePickerController =
      TextEditingController(text: DateTime.now().toString());

  late List<BusRoute> _busRoutes;

  late List<Item> _busRouteItems;

  // The instance field that holds the state of the time type toggle button
  int? timeTypeToggleIndex = 0;

  // The available options for the time type toggle button
  List<TimeType> timeTypes = [TimeType.departure, TimeType.arrival];

  // The route option items used for the ExpandablePanel
  // late List<Item> _items;

  // getters
  TextEditingController get dateTimePickerController => _dateTimePickerController;

  get destinationSelectionKey => _destinationSelectionKey;

  get originSelectionKey => _originSelectionKey;

  GlobalKey<FormState> get formKey => _formKey;

  List<Item> get busRouteItems => _busRouteItems; //
  // List<Item> get items => _items;

  List<BusRoute> get busRoutes => _busRoutes;

  Future<void> fetchBusRoute(BusRouteSearchFilter searchFilter) async {
    // start loading, hide route options
    visibilityRouteOptions = false;
    visibilityLoadingIcon = true;
    _busRoutes.clear();
    _busRouteItems.clear();
    // notify immediately because the below code will block execution until api returned
    notifyListeners();

    String url = '$apiHost/api/route/matchingRoute';
    url += '/${searchFilter.originStopNumber}/${searchFilter.destinationStopNumber}'
        '/${searchFilter.timeType.name}/${searchFilter.time}';
    final response = await http.get(
      Uri.parse(url),
      headers: {
        "Accept": "application/json",
      },
    );

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
        _busRouteItems = generateItems(busRouteList);
      } else {
        _busRoutes = [];
        _busRouteItems = [];
      }

      notifyListeners();
      // return busRouteList;
    } else {
      // TODO display some error message where this happens
      // If the server did not return a 200 OK response, then throw an exception.
      throw Exception('Failed to load bus routes');
    }
  }

  List<Item> generateItems(List<BusRoute> data) {
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
      );
    });
  }
}

class Item {
  Item({
    required this.headerValue,
    required this.busRoute,
    this.isExpanded = false,
    required this.expandedValue,
    required this.expandedDetailsValue,
  });

  String headerValue;
  bool isExpanded;
  BusRoute busRoute;
  String expandedValue;
  String expandedDetailsValue;
}
