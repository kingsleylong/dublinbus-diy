import 'package:flutter/material.dart';
import 'package:web/views/googlemap.dart';

class DesktopBody extends StatefulWidget {
  const DesktopBody({Key? key}) : super(key: key);

  @override
  State<DesktopBody> createState() => _DesktopBodyState();
}

class _DesktopBodyState extends State<DesktopBody>
    with TickerProviderStateMixin{
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Dublin Bus DIY")),
      body: Stack(
        alignment: Alignment.topLeft,
        children: [
          GoogleMapComponent(),

          Container(
            width: 400,
            child: Column(
              children: [
                ColoredBox(
                  color: Colors.green,
                  child: TabBar(
                    // expand the tab bar out of range and slide the bar when
                    // clicking tabs at the edges
                    isScrollable: true,
                    controller: _tabController,
                    tabs: const [
                      Tab(text: "Plan My Journey"),
                      Tab(text: "Find My Route"),
                      Tab(text: "Get Me There On-Time"),
                    ],
                  ),
                ),
              ],
          )

          //search fields

          //weather info
          ),
        ],
      )
    );
  }
}
