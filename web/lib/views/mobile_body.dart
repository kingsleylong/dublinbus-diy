import 'package:flutter/material.dart';
import 'package:web/views/tabs/plan_journey_mobile_tabview.dart';

class MobileBody extends StatefulWidget {
  const MobileBody({Key? key, required this.tabController}) : super(key: key);

  final TabController tabController;

  @override
  State<MobileBody> createState() => _MobileBodyState();
}

class _MobileBodyState extends State<MobileBody> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // Create an AppBar https://docs.flutter.dev/cookbook/design/tabs#interactive-example
      appBar: AppBar(
        title: const Text("Dublin Bus DIYs"),
        // bottom: TabBar(
        //   // expand the tab bar out of range and slide the bar when clicking
        //   // tabs at the edges https://stackoverflow.com/a/60636918
        //   // isScrollable: true,
        //   // Access a field of the widget in its state https://stackoverflow.com/a/58767810
        //   controller: widget.tabController,
        //   tabs: const [
        //     Tab(text: "Plan My Journey"),
        //   ],
        // ),
      ),
      body: TabBarView(controller: widget.tabController, children: const <Widget>[
        PlanMyJourneyTabMobileView(),
      ]),
    );
  }
}
