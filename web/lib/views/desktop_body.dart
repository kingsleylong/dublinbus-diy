import 'dart:convert';

import 'package:flutter/material.dart';

import 'package:web/models/bus_route.dart';
import 'package:web/models/bus_stop.dart';
import 'package:web/views/googlemap.dart';

import 'tab_views.dart';
import 'tabs.dart';

class DesktopBody extends StatefulWidget {
  const DesktopBody(
      {Key? key, required this.tabController, required this.futureAllBusStops}
      ) : super(key: key);
  final TabController tabController;
  final Future<List<BusStop>> futureAllBusStops;

  @override
  State<DesktopBody> createState() => _DesktopBodyState();
}

class _DesktopBodyState extends State<DesktopBody> {
  final GoogleMapComponent googleMapComponent = const GoogleMapComponent();

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
    print('Build right information box');
    return TabBarView(
      controller: widget.tabController,
      children: <Widget>[
        googleMapComponent,
        const Center(
          child: Text("It's rainy here"),
        ),
        const Center(
          child: Text("It's sunny here"),
        ),
      ]
    );
  }

  TabBarView buildLeftTabViews() {
    return TabBarView(
        controller: widget.tabController,
        children: <Widget>[
          PlanMyJourneyTabView(
              futureAllBusStops: widget.futureAllBusStops,
              googleMapComponent: googleMapComponent),
          const Center(
            child: Text("It's rainy here"),
          ),
          const Center(
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
}
