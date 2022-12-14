import 'package:flutter/material.dart';

import 'desktop_body.dart';
import 'mobile_body.dart';
import 'responsive_layout.dart';
import 'tabs.dart';

class HomePage extends StatefulWidget {
  const HomePage({Key? key}) : super(key: key);

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> with TickerProviderStateMixin {
  // The tab controller will be shared by all responsive views to keep the tab selection
  // consistent when the screen size changes
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: tabList.length, vsync: this);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: ResponsiveLayout(
        mobileBody: MobileBody(
            tabController: _tabController,
        ),
        desktopBody: DesktopBody(
          tabController: _tabController,
        ),
      ),
    );
  }
}
