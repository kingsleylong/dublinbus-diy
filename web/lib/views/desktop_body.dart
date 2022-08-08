import 'package:flutter/material.dart';

import 'googlemap.dart';
import 'tab_views.dart';

class DesktopBody extends StatefulWidget {
  const DesktopBody({Key? key, required this.tabController}) : super(key: key);
  final TabController tabController;

  @override
  State<DesktopBody> createState() => _DesktopBodyState();
}

class _DesktopBodyState extends State<DesktopBody> {

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
        ));
  }

  TabBarView buildRightInformationBox() {
    print('Build right information box');
    return TabBarView(
        controller: widget.tabController,
        // disable swiping from TabBarView
        // https://flutteragency.com/how-to-disable-swipe-tabbar-in-flutter/
        physics: const NeverScrollableScrollPhysics(),
        children: const <Widget>[
          GoogleMapComponent(),
        ]);
  }

  TabBarView buildLeftTabViews() {
    return TabBarView(
        controller: widget.tabController,
        physics: const NeverScrollableScrollPhysics(),
        children: const <Widget>[
          PlanMyJourneyTabView(),
        ]);
  }

  SizedBox buildLeftBar() {
    return SizedBox(
      width: 350,
      child: Column(
        children: [
          // ColoredBox(
          //   color: Colors.red,
          //   child: TabBar(
          //     // expand the tab bar out of range and slide the bar when clicking
          //     // tabs at the edges https://stackoverflow.com/a/60636918
          //     // isScrollable: true,
          //     // Access a field of the widget in its state https://stackoverflow.com/a/58767810
          //     controller: widget.tabController,
          //     tabs: tabList,
          //   ),
          // ),
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
