import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:provider/provider.dart';
import 'package:web/models/bus_route.dart';
import 'package:web/models/bus_stop.dart';
import 'package:web/models/map_polylines.dart';
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
  final Future<List<BusStop>> futureAllBusStops;
  const PlanMyJourneyTabView({Key? key, required this.futureAllBusStops,
    required this.googleMapComponent}) : super(key: key);
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
              child: buildFutureOriginDropdownList(widget),
            ),
            // Destination dropdown list
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
              child: buildFutureDestinationDropdownList(widget),
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
                    print('Selected origin: $originDropdownValue');
                    print('Selected destination: $destinationDropdownValue');
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

  Widget buildFutureOriginDropdownList(PlanMyJourneyTabView widget) {
    return FutureBuilder<List<BusStop>>(
      future: widget.futureAllBusStops,
      builder: (context, snapshot) {
        if (snapshot.hasData) {
          print('Building dropdown: data length = ${snapshot.data!.length}');
          return DropdownButtonFormField(
            value: originDropdownValue,
            items: snapshot.data!.map<DropdownMenuItem<String>>((BusStop value) {
              return DropdownMenuItem<String>(
                value: value.stopNumber,
                child: Text('${value.stopName} - ${value.stopNumber}'),
              );
            }).toList(),
            onChanged: (String? value) {
              setState(() {
                originDropdownValue = value!;
                visibilityRouteOptions = false;
              });
            },
            validator: (String? value) {
              if (value == null || value.isEmpty) {
                return 'Please enter the origin';
              }
              return null;
            },
            decoration: const InputDecoration(
              // icon: Icon(Icons.),
              labelText: "Origin",
              floatingLabelAlignment: FloatingLabelAlignment.start,
              hintText: 'Origin',
              // helperText: 'Select the origin',
              // counterText: '0 characters',
              border: OutlineInputBorder(),
            ),
          );
        } else if (snapshot.hasError) {
          print('${snapshot.error}');
          return Text('${snapshot.error}');
        }
        // By default, show a loading spinner.
        return const CircularProgressIndicator();
      },
    );
  }

  Widget buildFutureDestinationDropdownList(PlanMyJourneyTabView widget) {
    return FutureBuilder<List<BusStop>>(
      future: widget.futureAllBusStops,
      builder: (context, snapshot) {
        if (snapshot.hasData) {
          print('Building destination dropdown: data length = ${snapshot.data!.length}');
          return DropdownButtonFormField(
            value: destinationDropdownValue,
            items: snapshot.data!.map<DropdownMenuItem<String>>((BusStop value) {
              return DropdownMenuItem<String>(
                value: value.stopNumber,
                child: Text('${value.stopName} - ${value.stopNumber}'),
              );
            }).toList(),
            onChanged: (String? value) {
              setState(() {
                destinationDropdownValue = value!;
                visibilityRouteOptions = false;
              });
            },
            validator: (String? value) {
              if (value == null || value.isEmpty) {
                return 'Please enter the destination';
              }
              return null;
            },
            decoration: const InputDecoration(
              // icon: Icon(Icons.),
              labelText: "Destination",
              floatingLabelAlignment: FloatingLabelAlignment.start,
              hintText: 'Destination',
              // helperText: 'Select the origin',
              // counterText: '0 characters',
              border: OutlineInputBorder(),
            ),
          );
        } else if (snapshot.hasError) {
          print('${snapshot.error}');
          return Text('${snapshot.error}');
        }
        // By default, show a loading spinner.
        return const CircularProgressIndicator();
      },
    );
  }

  Future<List<BusRoute>> fetchBusRoutes() async {
    final response = await http.get(
      Uri.parse('http://ipa-003.ucd.ie/api/matchingRoute/${originDropdownValue}/${destinationDropdownValue}'),
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
      for(BusRoute busRoute in busRouteList) {
        Provider.of<PolylinesModel>(context, listen: false).addBusRouteAsPolyline(busRoute);
      }
      print('PolylinesModel size: ${Provider.of<PolylinesModel>(context, listen: false).items
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

  buildRouteOptionsListView(PlanMyJourneyTabView widget) {
    return FutureBuilder<List<BusRoute>>(
      future: futureBusRoutes,
      builder: (context, snapshot) {
        if (snapshot.hasData) {
          print("Number of routes found: ${snapshot.data!.length}");
          return ListView(
            // shrinkWrap: resolve the error of "Vertical viewport was given unbounded height."
            // https://stackoverflow.com/a/57335217/12328041
            shrinkWrap: true,
            children: snapshot.data!.map(
                    // Card layout: https://docs.flutter.dev/development/ui/layout#card
                    (busRoute) => OutlinedButton(
                      onPressed: () {
                        print('pressed');
                      },
                      child: ListTile(
                          title: Text(
                              busRoute.routeNumber,
                              style: const TextStyle(
                                fontWeight: FontWeight.w900,
                                fontSize: 18,
                              )
                          ),
                          subtitle: Text('${busRoute.stops.length} stops. From ${busRoute.stops[0]
                              .stopName}.'),
                        ),
                    )).toList(),
          );
        } else if (snapshot.hasError) {
          return Text('${snapshot.error}');
        }
        // By default, show a loading spinner.
        return const CircularProgressIndicator();
      },
    );
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
