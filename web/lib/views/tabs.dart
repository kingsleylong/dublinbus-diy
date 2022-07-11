import 'package:flutter/material.dart';

const List<Tab> tabList = [
  Tab(text: "Plan My Journey"),
  Tab(text: "Find My Route"),
  Tab(text: "Get Me There On-Time"),
];

class PlanMyJourneyTab extends StatefulWidget {
  const PlanMyJourneyTab({Key? key}) : super(key: key);

  @override
  State<PlanMyJourneyTab> createState() => _PlanMyJourneyTabState();
}

class _PlanMyJourneyTabState extends State<PlanMyJourneyTab> {
  @override
  Widget build(BuildContext context) {
    return Container();
  }
}


