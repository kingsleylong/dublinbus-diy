import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:web/models/bus_stop.dart';

import 'tabs.dart';

class DesktopBody extends StatefulWidget {
  const DesktopBody({Key? key, required this.tabController}) : super(key: key);

  final TabController tabController;

  @override
  State<DesktopBody> createState() => _DesktopBodyState();
}

class _DesktopBodyState extends State<DesktopBody> {
  final _lines = <String>["175", "C1", "46A", "52"];
  late Future<BusStop> futureBusStop;

  @override
  void initState() {
    super.initState();
    futureBusStop = fetchBusStop();
  }

  Future<BusStop> fetchBusStop() async {
    final response = await http.get(
      Uri.parse('http://localhost:1080/busStop/1'),
      headers: {
        "Accept": "application/json",
      },
    );

    if (response.statusCode == 200) {
      // If the server did return a 200 OK response,
      // then parse the JSON.
      return BusStop.fromJson(jsonDecode(response.body));
    } else {
      // If the server did not return a 200 OK response,
      // then throw an exception.
      throw Exception('Failed to load bus stop');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Dublin Bus DIY")),
      body: Row(
        // alignment: Alignment.topLeft,
        children: [
          // left bar
          buildLeftBar(),
          // right information box, use Expanded class to take the rest of space
          Expanded(
            child: buildRightInformationBox(),
          )
        ],
      )
    );
  }

  TabBarView buildRightInformationBox() {
    return TabBarView(
      controller: widget.tabController,
      children: <Widget>[
        PlanMyJourneyTab(),
        Center(
          child: Text("It's rainy here"),
        ),
        Center(
          child: Text("It's sunny here"),
        ),
      ]
    );
  }

  SizedBox buildLeftBar() {
    return SizedBox(
      width: 350,
      child: Column(
        children: [
          ColoredBox(
            color: Colors.red,
            child: TabBar(
              // expand the tab bar out of range and slide the bar when clicking
              // tabs at the edges https://stackoverflow.com/a/60636918
              isScrollable: true,
              // Access a field of the widget in its state https://stackoverflow.com/a/58767810
              controller: widget.tabController,
              tabs: const [
                Tab(text: "Plan My Journey"),
                Tab(text: "Find My Route"),
                Tab(text: "Get Me There On-Time"),
              ],
            ),
          ),

          //search fields
          Expanded(
            // build a list view from data
            // https://codelabs.developers.google.com/codelabs/first-flutter-app-pt1/#5
            child: ListView.builder(
              // specify the length of data, without this an index out of range error will be
              // thrown. https://stackoverflow.com/a/58850610
              itemCount: _lines.length,
              itemBuilder: (context, index) {
                return ListTile(
                  title: Text(
                    _lines[index],
                    // strutStyle: ,
                  ),
                );
              }
            ),
          ),

          Expanded(
              child: FutureBuilder<BusStop>(
              future: futureBusStop,
              builder: (context, snapshot) {
                if (snapshot.hasData) {
                  print(snapshot.data!.name);
                  return Text(snapshot.data!.name);
                } else if (snapshot.hasError) {
                  print('${snapshot.error}');
                  return Text('${snapshot.error}');
                }

                // By default, show a loading spinner.
                return const CircularProgressIndicator();
              },
            )
          ),

        //weather info
          // TODO

        ],
      ),
    );
  }
}
