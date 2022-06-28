import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

import 'tabs.dart';

class DesktopBody extends StatefulWidget {
  const DesktopBody({Key? key, required this.tabController}) : super(key: key);

  final TabController tabController;

  @override
  State<DesktopBody> createState() => _DesktopBodyState();
}

class _DesktopBodyState extends State<DesktopBody> {
  final _lines = <String>["175", "C1", "46A", "52"];

  Future<http.Response> fetchLines() {
    return http.get(Uri.parse(''));
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
          )

          //weather info
          // TODO

        ],
      ),
    );
  }
}
