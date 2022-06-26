import 'package:flutter/material.dart';

class AppBody extends StatelessWidget {
  const AppBody({super.key});

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 3,
      child: Scaffold(
        appBar: AppBar(
          bottom: const TabBar(
            tabs: [
              Tab(text: "Plan My Journey",),
              Tab(text: "Find My Route",),
              Tab(text: "Get Me There On-Time",),
          ],
        ),
        ),
        body: const TabBarView(
          children: [
            Text("Plan My Journey"),
            Text("Find My Route"),
            Text("Get Me There On-Time"),
          ],
        ),
      ),
    );
  }

}