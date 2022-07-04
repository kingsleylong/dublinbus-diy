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
  String dropdownValue = '175';
  String? originDropdownValue;
  String? destinationDropdownValue;
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();

  late List<DropdownMenuItem<int>> _menuItems;

  @override
  void initState() {
  }

  @override
  Widget build(BuildContext context) {
    const searchFieldsDecoration = InputDecoration(
      // icon: Icon(Icons.),
      labelText: "Origin",
      floatingLabelAlignment: FloatingLabelAlignment.start,
      hintText: 'Origin',
      // helperText: 'Select the origin',
      // counterText: '0 characters',
      border: OutlineInputBorder(),
    );

    return Padding(
      // padding settings https://api.flutter.dev/flutter/material/InputDecoration/contentPadding.html
      padding: const EdgeInsets.fromLTRB(5, 10, 5, 0),
      child: Form(
        key: _formKey,
        child: Column(
          children: <Widget>[
            buildFutureOriginDropdownList(widget),
            buildFutureDestinationDropdownList(widget),
            DropdownButtonFormField(
              // how to build a drop down list https://api.flutter.dev/flutter/material/DropdownButton-class.htm
              value: dropdownValue,
              // Field decoration https://api.flutter.dev/flutter/material/InputDecoration-class.html
              decoration: searchFieldsDecoration,
              items: <String>["175", "C1", "46A", "52"]
                  .map<DropdownMenuItem<String>>((String value) {
                return DropdownMenuItem<String>(
                  value: value,
                  child: Text(value),
                );
              }).toList(),
              onChanged: (String? value) {
                setState(() {
                  dropdownValue = value!;
                });
              },
              // isExpanded: true,
            ),
            DropdownButtonFormField(
              // how to build a drop down list https://api.flutter.dev/flutter/material/DropdownButton-class.html
              value: dropdownValue,
              // Field decoration https://api.flutter.dev/flutter/material/InputDecoration-class.html
              decoration: searchFieldsDecoration,
              items: <String>["175", "C1", "46A", "52"]
                  .map<DropdownMenuItem<String>>((String value) {
                return DropdownMenuItem<String>(
                  value: value,
                  child: Text(value),
                );
              }).toList(),
              onChanged: (String? value) {
                setState(() {
                  dropdownValue = value!;
                });
              },
              // isExpanded: true,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 16.0),
              child: ElevatedButton(
                onPressed: () {
                  // Validate will return true if the form is valid, or false if
                  // the form is invalid.
                  if (_formKey.currentState!.validate()) {
                    // Process data.
                  }
                },
                child: const Text('Plan'),
              ),
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
                value: value.stopId,
                child: Text(value.stopName),
              );
            }).toList(),
            onChanged: (String? value) {
              setState(() {
                originDropdownValue = value!;
              });
            },
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
                value: value.stopId,
                child: Text(value.stopName),
              );
            }).toList(),
            onChanged: (String? value) {
              setState(() {
                destinationDropdownValue = value!;
              });
            },
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

