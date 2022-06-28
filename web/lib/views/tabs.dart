import 'package:flutter/material.dart';
import 'package:web/views/googlemap.dart';

class GetMeThereOnTimeTab extends StatefulWidget {
  const GetMeThereOnTimeTab({Key? key}) : super(key: key);

  @override
  State<GetMeThereOnTimeTab> createState() => _GetMeThereOnTimeTabState();
}

class _GetMeThereOnTimeTabState extends State<GetMeThereOnTimeTab> {
  @override
  Widget build(BuildContext context) {
    return Container(
        child: const Text("Get me there on time")
    );
  }
}

class PlanMyJourneyTab extends StatefulWidget {
  const PlanMyJourneyTab({Key? key}) : super(key: key);

  @override
  State<PlanMyJourneyTab> createState() => _PlanMyJourneyTabState();
}

class _PlanMyJourneyTabState extends State<PlanMyJourneyTab> {
  @override
  Widget build(BuildContext context) {
    return Container(
        child: GoogleMapComponent(),
    );
  }
}

