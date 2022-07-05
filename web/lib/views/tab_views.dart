import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:web/models/bus_route.dart';

import '../models/bus_stop.dart';

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
  const PlanMyJourneyTabView({Key? key, required this.futureAllBusStops}) : super(key: key);

  @override
  State<PlanMyJourneyTabView> createState() => _PlanMyJourneyTabViewState();
}

class _PlanMyJourneyTabViewState extends State<PlanMyJourneyTabView> {
  String? originDropdownValue;
  String? destinationDropdownValue;
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
  late Future<List<BusRoute>> futureBusRoutes;
  // Use a flag to control the visibility of the route options
  // https://stackoverflow.com/a/46126667
  bool visibilityRouteOptions = false;

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
              Padding(
                padding: const EdgeInsets.all(8),
                child: buildRouteOptionsListView(widget),
              ),
            Expanded(
              child: Container(),
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
      Uri.parse('http://localhost:1080/api/matchingRoute'),
      headers: {
        "Accept": "application/json",
      },
    );

    if (response.statusCode == 200) {
      // If the server did return a 200 OK response, then parse the JSON.
      final List busRoutesJson = jsonDecode(response.body);

      print("Bus route list size: ${busRoutesJson.length}");
      return List.generate(busRoutesJson.length,
              (index) => BusRoute.fromJson(busRoutesJson[index]));
    } else {
      // If the server did not return a 200 OK response, then throw an exception.
      throw Exception('Failed to load bus routes');
    }
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
                    (busRoute) => Card(
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

