import 'dart:convert';

import 'package:dropdown_search/dropdown_search.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:provider/provider.dart';
import 'package:web/models/bus_route.dart';
import 'package:web/models/bus_stop.dart';
import 'package:web/models/map_polylines.dart';
import 'package:web/env.dart';
import 'package:web/views/googlemap.dart';

class GetMeThereOnTimeTabView extends StatefulWidget {
  const GetMeThereOnTimeTabView({Key? key}) : super(key: key);

  @override
  State<GetMeThereOnTimeTabView> createState() => _GetMeThereOnTimeTabViewState();
}

class _GetMeThereOnTimeTabViewState extends State<GetMeThereOnTimeTabView> {
  @override
  Widget build(BuildContext context) {
    return Container(
        child: const Text("Get me there on time")
    );
  }
}

class PlanMyJourneyTabView extends StatefulWidget {
  const PlanMyJourneyTabView({Key? key, required this.googleMapComponent}) : super(key: key);
  final GoogleMapComponent googleMapComponent;

  @override
  State<PlanMyJourneyTabView> createState() => _PlanMyJourneyTabViewState();
}

class _PlanMyJourneyTabViewState extends State<PlanMyJourneyTabView> {
  String? originDropdownValue = "3161";
  String? destinationDropdownValue = "3163";
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
  late Future<List<BusRoute>> futureBusRoutes;
  // Use a flag to control the visibility of the route options
  // https://stackoverflow.com/a/46126667
  bool visibilityRouteOptions = false;

  late List<Item> items;

  // The instance field that holds the state of origin dropdown list
  final _originSelectionKey = GlobalKey<DropdownSearchState<BusStop>>();
  // The instance field that holds the state of destination dropdown list
  final _destinationSelectionKey = GlobalKey<DropdownSearchState<BusStop>>();

  @override
  Widget build(BuildContext context) {
    return Padding(
      // padding settings https://api.flutter.dev/flutter/material/InputDecoration/contentPadding.html
      padding: const EdgeInsets.fromLTRB(5, 10, 5, 0),
      child: Form(
        key: _formKey,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: <Widget>[
            // Origin dropdown list
            // The form fields should be wrapped by Padding otherwise they would overlap each other
            // https://docs.flutter.dev/cookbook/forms/text-input
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
              child: buildSearchableOriginDropdownList(),
            ),
            // Destination dropdown list
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
              child: buildSearchableDestinationDropdownList(),
            ),
            // Submit button
            Padding(
              padding: const EdgeInsets.all(8),
              child: ElevatedButton(
                style: ElevatedButton.styleFrom(
                    minimumSize: const Size.fromHeight(60),
                    textStyle: const TextStyle(fontSize: 18),
                ),
                onPressed: () {
                  // Validate will return true if the form is valid, or false if
                  // the form is invalid.
                  if (_formKey.currentState!.validate()) {
                    Provider.of<PolylinesModel>(context, listen: false).removeAll();
                    futureBusRoutes = fetchBusRoutes();
                    setState(() {
                      visibilityRouteOptions = true;
                    });
                  }
                },
                child: const Text('Plan'),
              ),
            ),
            if(visibilityRouteOptions)
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
        )
      )
    );
  }

  Widget buildSearchableOriginDropdownList() {
    // DropdownSearch widget plugin: https://pub.dev/packages/dropdown_search
    // Check the examples code for usage: https://github.com/salim-lachdhaf/searchable_dropdown
    return DropdownSearch<BusStop>(
      key: _originSelectionKey,
      asyncItems: (filter) => fetchFutureBusStopsByName(filter),
      compareFn: (i, s) => i.isEqual(s),
      popupProps: PopupProps.menu(
        showSearchBox: true,
        title: const Text('Search origin bus stop'),
        isFilterOnline: true,
        showSelectedItems: true,
        itemBuilder: _originPopupItemBuilder,
        favoriteItemProps: FavoriteItemProps(
          showFavoriteItems: true,
          // TODO This is a fake favorite feature. We need to implement in future to let the user
          //  mark the favorite stops
          favoriteItems: (us) {
            return us
              .where((e) => e.stopName!.contains("UCD"))
              .toList();
          },
        ),
      ),
    );
  }

  Widget buildSearchableDestinationDropdownList() {
    return DropdownSearch<BusStop>(
      key: _destinationSelectionKey,
      asyncItems: (filter) => fetchFutureBusStopsByName(filter),
      compareFn: (i, s) => i.isEqual(s),
      popupProps: PopupProps.menu(
        showSearchBox: true,
        title: const Text('Search destination bus stop'),
        isFilterOnline: true,
        showSelectedItems: true,
        itemBuilder: _originPopupItemBuilder,
        favoriteItemProps: FavoriteItemProps(
          showFavoriteItems: true,
          favoriteItems: (us) {
            return us
                .where((e) => e.stopName!.contains("Spire"))
                .toList();
          },
        ),
      ),
    );
  }

  Future<List<BusStop>> fetchFutureBusStopsByName(String filter) async {
    String url = '$apiHost/api/stop/findByAddress';
    final paramsStr = (filter == '') ? '' : '&filter=$filter';
    final response = await http.get(
      Uri.parse('$url$paramsStr'),
      headers: {
        "Accept": "application/json",
      },
    );

    if (response.statusCode == 200) {
      // If the server did return a 200 OK response, then parse the JSON.
      final List busStopsJson = jsonDecode(response.body);

      print("Bus stops size: ${busStopsJson.length}");
      List<BusStop> busStopList = List.generate(busStopsJson.length, (index) => BusStop.fromJson
        (busStopsJson[index]));
      // TODO There is one bug here. The only one polyline was shown on the map, but if you
      //  switch to other tabs and come back, all polylines will be there.

      // Provider.of<PolylinesModel>(context, listen: false).addBusRouteListAsPolylines(busRouteList);
      // busRouteList.map((busRoute) => (){
      //     Provider.of<PolylinesModel>(context, listen: false).addBusRouteAsPolyline(busRoute);
      //     // print(Provider.of<PolylinesModel>(context, listen: false));
      //   }
      // );
      // for(BusRoute busRoute in busStopList) {
      //   Provider.of<PolylinesModel>(context, listen: false).addBusRouteAsPolyline(busRoute);
      // }
      // print('PolylinesModel size: ${Provider.of<PolylinesModel>(context, listen: false).items
      //     .length}');

      return busStopList;
    } else {
      // If the server did not return a 200 OK response, then throw an exception.
      throw Exception('Failed to load bus routes');
    }
  }

  Widget _originPopupItemBuilder(BuildContext context, BusStop item, bool isSelected) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 8),
      decoration: !isSelected
          ? null
          : BoxDecoration(
        border: Border.all(color: Theme.of(context).primaryColor),
        borderRadius: BorderRadius.circular(5),
        color: Colors.white,
      ),
      child: ListTile(
        selected: isSelected,
        title: Text(item.stopName!),
        subtitle: Text(item.stopNumber.toString()),
        leading: CircleAvatar(
          // this does not work - throws 404 error
          // backgroundImage: NetworkImage(item.avatar ?? ''),
        ),
      ),
    );
  }

  Future<List<BusRoute>> fetchBusRoutes() async {
    print('Selected origin: ${_originSelectionKey.currentState?.getSelectedItem?.stopNumber}');
    print('Selected destination: ${_destinationSelectionKey.currentState?.getSelectedItem?.stopNumber}');
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
      List<BusRoute> busRouteList = List.generate(busRoutesJson.length, (index) => BusRoute.fromJson
        (busRoutesJson[index]));
      items = generateItems(busRouteList);
      // TODO There is one bug here. The only one polyline was shown on the map, but if you
      //  switch to other tabs and come back, all polylines will be there.

      // Provider.of<PolylinesModel>(context, listen: false).addBusRouteListAsPolylines(busRouteList);
      // busRouteList.map((busRoute) => (){
      //     Provider.of<PolylinesModel>(context, listen: false).addBusRouteAsPolyline(busRoute);
      //     // print(Provider.of<PolylinesModel>(context, listen: false));
      //   }
      // );
      Provider.of<PolylinesModel>(context, listen: false).addBusRouteListAsPolylines(busRouteList);
      // for(BusRoute busRoute in busRouteList) {
      //   Provider.of<PolylinesModel>(context, listen: false).addBusRouteAsPolyline(busRoute);
        // Provider.of<MarkersModel>(context, listen: false).addBusRouteAsMarker(busRoute);
      // }
      print('PolylinesModel size: ${Provider.of<PolylinesModel>(context, listen: false).itemsOfPolylines
          .length}');

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
                  title: Text(item.headerValue),
                  onTap: () {
                    setState(() {
                      item.isExpanded = !isExpanded;
                    });
                  },
                );
              },
              body: ListTile(
                title: Text(item.expandedValue),
              ),
              isExpanded: item.isExpanded,
          );
        }).toList(),
    );
  }

  List<Item> generateItems(List<BusRoute> data) {
    return List<Item>.generate(data.length, (int index) {
      return Item(
        headerValue: data[index].routeNumber,
        expandedValue: '${data[index].stops.length} stops. Starts from ${data[index].stops[0].stopName}',
      );
    });
  }
}

// stores ExpansionPanel state information
class Item {
  Item({
    required this.expandedValue,
    required this.headerValue,
    this.isExpanded = false,
  });

  String expandedValue;
  String headerValue;
  bool isExpanded;
}
