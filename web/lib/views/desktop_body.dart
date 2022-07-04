import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:web/models/bus_route.dart';
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
  late Future<List<BusStop>> futureAllBusStops;
  late Future<List<BusRoute>> futureBusRoutes;

  @override
  void initState() {
    super.initState();
    futureAllBusStops = fetchAllBusStops();
    futureBusRoutes = fetchBusRoutes();
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
            child: buildRightInformationBox(),
          ),
        ],
      )
    );
  }

  TabBarView buildRightInformationBox() {
    return TabBarView(
      controller: widget.tabController,
      children: const <Widget>[
        GoogleMapComponent(),
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
        children: <Widget>[
          PlanMyJourneyTabView(futureAllBusStops: futureAllBusStops),
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
          Expanded(
            child: buildLeftTabViews(),
          ),

        //weather info
          // TODO

        ],
      ),
    );
  }

  Future<List<BusStop>> fetchAllBusStops() async {
    final response = await http.get(
      Uri.parse('http://localhost:1080/api/allStops'),
      headers: {
        "Accept": "application/json",
      },
    );

    if (response.statusCode == 200) {
      // If the server did return a 200 OK response, then parse the JSON.
      final List allBusStops = jsonDecode(response.body);

      print("Bus Stop list size: ${allBusStops.length}");
      return List.generate(allBusStops.length,
              (index) => BusStop.fromJson(allBusStops[index]));
    } else {
      // If the server did not return a 200 OK response, then throw an exception.
      throw Exception('Failed to load bus stop');
    }
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
}
