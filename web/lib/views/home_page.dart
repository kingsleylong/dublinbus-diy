import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:web/models/map_polylines.dart';
import 'package:web/models/search_form.dart';
import 'package:web/views/responsive_layout.dart';
import 'package:web/views/tabs.dart';

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
    _tabController = TabController(length: tabList.length, vsync: this);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // Create a model by the provider so the child can listen to the model changes
      // https://docs.flutter.dev/development/data-and-backend/state-mgmt/simple#changenotifierprovider
      body: MultiProvider(
        providers: [
          ChangeNotifierProvider(create: (context) => PolylinesModel()),
          ChangeNotifierProvider(create: (context) => SearchFormModel()),
        ],
        child: ResponsiveLayout(
          mobileBody: MobileBody(
              tabController: _tabController,
          ),
          desktopBody: DesktopBody(
            tabController: _tabController,
          ),
        ),
      ),
    );
  }
}
