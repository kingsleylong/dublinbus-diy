import 'dart:convert';

import 'package:date_time_picker/date_time_picker.dart';
import 'package:dropdown_search/dropdown_search.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:web/models/bus_route.dart';
import 'package:web/models/bus_stop.dart';
import 'package:web/models/map_polylines.dart';
import 'package:web/env.dart';
import 'package:web/views/googlemap.dart';
import 'package:web/views/tabs/search_panel.dart';

class GetMeThereOnTimeTabView extends StatefulWidget {
  const GetMeThereOnTimeTabView({Key? key}) : super(key: key);

  @override
  State<GetMeThereOnTimeTabView> createState() =>
      _GetMeThereOnTimeTabViewState();
}

class _GetMeThereOnTimeTabViewState extends State<GetMeThereOnTimeTabView> {
  @override
  Widget build(BuildContext context) {
    return Container(child: const Text("Get me there on time"));
  }
}

class PlanMyJourneyTabView extends StatefulWidget {
  const PlanMyJourneyTabView({Key? key, required this.googleMapComponent})
      : super(key: key);
  final GoogleMapComponent googleMapComponent;

  @override
  State<PlanMyJourneyTabView> createState() => _PlanMyJourneyTabViewState();
}

class _PlanMyJourneyTabViewState extends State<PlanMyJourneyTabView> {
  late Future<List<BusRoute>> futureBusRoutes;

  // Use a flag to control the visibility of the route options
  // https://stackoverflow.com/a/46126667
  bool visibilityRouteOptions = false;

  late List<Item> items;

  // The instance field that holds the state of origin dropdown list
  final _originSelectionKey = GlobalKey<DropdownSearchState<BusStop>>();

  // The instance field that holds the state of destination dropdown list
  final _destinationSelectionKey = GlobalKey<DropdownSearchState<BusStop>>();
  late TextEditingController _dateTimePickerController;

  @override
  void initState() {
    super.initState();
    _dateTimePickerController =
        TextEditingController(text: DateTime.now().toString());
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
        // padding settings https://api.flutter.dev/flutter/material/InputDecoration/contentPadding.html
        padding: const EdgeInsets.fromLTRB(5, 10, 5, 0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: <Widget>[
            const SearchForm(),
            if (visibilityRouteOptions)
              Expanded(
                child: Padding(
                  padding: const EdgeInsets.all(8),
                  child: buildRouteOptionsListView2(widget),
                ),
              ),
            ConstrainedBox(
              constraints: const BoxConstraints(
                minHeight: 2.0,
              ),
            ),
          ],
        ));
  }

  Future<List<BusRoute>> fetchBusRoutes() async {
    print(
        'Selected origin: ${_originSelectionKey.currentState?.getSelectedItem?.stopNumber}');
    print(
        'Selected destination: ${_destinationSelectionKey.currentState?.getSelectedItem?.stopNumber}');
    print('Selected datetime: ${_dateTimePickerController.value.text}');
    var parseTime = DateTime.parse(_dateTimePickerController.value.text);
    // Date format: https://api.flutter.dev/flutter/intl/DateFormat-class.html
    print(
        'parseTime: $parseTime  ${DateFormat('MM-dd-yyyy HH:mm:ss').format(parseTime)}');
    final response = await http.get(
      Uri.parse('$apiHost/api/route/matchingRoute'),
      headers: {
        "Accept": "application/json",
      },
    );

    if (response.statusCode == 200) {
      // If the server did return a 200 OK response, then parse the JSON.
      final List busRoutesJson = jsonDecode(response.body);

      print("Bus route list size: ${busRoutesJson.length}");
      List<BusRoute> busRouteList = List.generate(busRoutesJson.length,
          (index) => BusRoute.fromJson(busRoutesJson[index]));
      items = generateItems(busRouteList);
      return busRouteList;
    } else {
      // If the server did not return a 200 OK response, then throw an exception.
      throw Exception('Failed to load bus routes');
    }
  }

  buildRouteOptionsListView2(PlanMyJourneyTabView widget) {
    return FutureBuilder<List<BusRoute>>(
      future: futureBusRoutes,
      builder: (context, snapshot) {
        if (snapshot.hasData) {
          return SingleChildScrollView(
            child: _buildRouteOptionPanels(items),
          );
        } else if (snapshot.hasError) {
          return Text('${snapshot.error}');
        }
        // By default, show a loading spinner.
        return const Center(
          child: CircularProgressIndicator(),
        );
      },
    );
  }

  _buildRouteOptionPanels(List<Item> data) {
    // Use ExpansionPanel to display the route options for easy use.
    // https://api.flutter.dev/flutter/material/ExpansionPanel-class.html
    return ExpansionPanelList(
      expansionCallback: (int index, bool isExpanded) {
        print("isExpanded: $isExpanded");
        setState(() {
          data[index].isExpanded = !isExpanded;
          print("new isExpanded: ${data[index].isExpanded}");
        });
      },
      children: data.map<ExpansionPanel>((Item item) {
        return ExpansionPanel(
          headerBuilder: (BuildContext context, bool isExpanded) {
            return ListTile(
              title: Text(
                item.headerValue,
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
              onTap: () {
                setState(() {
                  item.isExpanded = !isExpanded;
                });
                // add the polyline and marker for the selected route by changing
                // the state from the Provider and notify the Consumers.
                Provider.of<PolylinesModel>(context, listen: false)
                    .addBusRouteAsPolyline(item.busRoute);
              },
            );
          },
          body: ListTile(
            title: Text(item.expandedValue),
            subtitle: Center(
              child: Text(item.expandedDetailsValue),
            ),
          ),
          isExpanded: item.isExpanded,
        );
      }).toList(),
    );
  }

  List<Item> generateItems(List<BusRoute> data) {
    return List<Item>.generate(data.length, (int index) {
      return Item(
        // TODO integrate the travel time
        headerValue: '${data[index].routeNumber}      30 min',
        expandedValue:
            '${data[index].stops.length} stops. Starts from ${data[index].stops[0].stopName}',
        expandedDetailsValue: data[index]
            .stops
            .map((stop) => '${stop.stopName} - ${stop.stopNumber}')
            .reduce((value, element) => '${value}\n${element}'),
        busRoute: data[index],
      );
    });
  }
}

// stores ExpansionPanel state information
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
