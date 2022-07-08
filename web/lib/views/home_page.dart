import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:provider/provider.dart';
import 'package:web/models/bus_stop.dart';
import 'package:web/models/map_polylines.dart';
import 'package:web/views/responsive_layout.dart';

import 'desktop_body.dart';
import 'mobile_body.dart';

class HomePage extends StatefulWidget {
  const HomePage({Key? key}) : super(key: key);

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> with TickerProviderStateMixin {
  // The tab controller will be shared by all responsive views to keep the tab selection
  // consistent when the screen size changes
  late TabController _tabController;
  late Future<List<BusStop>> futureAllBusStops;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    futureAllBusStops = fetchAllBusStops();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // Create a model by the provider so the child can listen to the model changes
      // https://docs.flutter.dev/development/data-and-backend/state-mgmt/simple#changenotifierprovider
      body: ChangeNotifierProvider(
        create: (BuildContext context) => PolylinesModel(),
        child: ResponsiveLayout(
          mobileBody: MobileBody(
              tabController: _tabController,
              futureAllBusStops: futureAllBusStops
          ),
          desktopBody: DesktopBody(
            tabController: _tabController,
            futureAllBusStops: futureAllBusStops,
          ),
        ),
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
              (index) => BusStop.fromJsonForRoute(allBusStops[index]));
    } else {
      // If the server did not return a 200 OK response, then throw an exception.
      throw Exception('Failed to load bus stop');
    }
  }
}
