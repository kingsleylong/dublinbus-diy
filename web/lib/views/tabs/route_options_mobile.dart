import 'package:dublin_bus_diy/views/tabs/route_options.dart';
import 'package:flutter/material.dart';

class RouteOptionsMobile extends StatefulWidget {
  const RouteOptionsMobile({Key? key}) : super(key: key);

  @override
  State<RouteOptionsMobile> createState() => _RouteOptionsMobileState();
}

class _RouteOptionsMobileState extends State<RouteOptionsMobile> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Route Options'),
      ),
      body: const RouteOptions(),
    );
  }
}
