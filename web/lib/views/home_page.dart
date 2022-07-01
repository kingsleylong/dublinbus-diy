import 'package:flutter/material.dart';
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

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: ResponsiveLayout(
        mobileBody: MobileBody(tabController: _tabController),
        desktopBody: DesktopBody(tabController: _tabController),
      )
    );
  }
}
