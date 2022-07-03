import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:web/models/bus_stop.dart';
import 'package:web/views/googlemap.dart';

import 'tab_views.dart';
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
      Uri.parse('http://ipa-003.ucd.ie/api/busStop/2'),
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
          // const Expanded(
          //   child: GoogleMapComponent(),
          // )
          // buildLeftTabViews(),

          Expanded(
            child: TabBarView(
              controller: widget.tabController,
              children: const <Widget>[
                // PlanMyJourneyTabView(),
                Center(
                  child: Text("It's rainy here"),
                ),
                Center(
                  child: Text("It's sunny here"),
                ),
                Center(
                  child: Text("It's sunny here"),
                ),
              ]
            )
          ),
        ],
      )
    );
  }

  TabBarView buildRightInformationBox() {
    return TabBarView(
      controller: widget.tabController,
      children: const <Widget>[
        PlanMyJourneyTabView(),
        Center(
          child: Text("It's rainy here"),
        ),
        Center(
          child: Text("It's sunny here"),
        ),
      ]
    );
  }

  TabBarView buildLeftTabViews() {
    return TabBarView(
        controller: widget.tabController,
        children: const <Widget>[
          PlanMyJourneyTabView(),
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
              tabs: tabList,
            ),
          ),
          //
          // // tab views
          // buildLeftTabViews(),

          Expanded(
              child: FutureBuilder<BusStop>(
              future: futureBusStop,
              builder: (context, snapshot) {
                if (snapshot.hasData) {
                  print(snapshot.data!.stopName);
                  return Text(snapshot.data!.stopName);
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
