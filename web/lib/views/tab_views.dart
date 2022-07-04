import 'package:flutter/material.dart';

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

  late List<DropdownMenuItem<int>> _menuItems;

  @override
  void initState() {
  }

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
                  }
                },
                child: const Text('Plan'),
              ),
            ),
            // Padding(
            //   padding: const EdgeInsets.all(8),
            //   child: buildRouteOptionsListView(widget),
            // ),
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
}

