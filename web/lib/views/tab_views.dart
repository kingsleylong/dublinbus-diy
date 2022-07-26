import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:web/models/search_form.dart';
import 'package:web/views/googlemap.dart';
import 'package:web/views/tabs/route_options.dart';
import 'package:web/views/tabs/search_panel.dart';

class GetMeThereOnTimeTabView extends StatefulWidget {
  const GetMeThereOnTimeTabView({Key? key}) : super(key: key);

  @override
  State<GetMeThereOnTimeTabView> createState() => _GetMeThereOnTimeTabViewState();
}

class _GetMeThereOnTimeTabViewState extends State<GetMeThereOnTimeTabView> {
  @override
  Widget build(BuildContext context) {
    return Container(child: const Text("Get me there on time"));
  }
}

class PlanMyJourneyTabView extends StatelessWidget {
  const PlanMyJourneyTabView({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
        // padding settings https://api.flutter.dev/flutter/material/InputDecoration/contentPadding.html
        padding: const EdgeInsets.fromLTRB(5, 10, 5, 0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: <Widget>[
            const SearchForm(),
            const Expanded(
              child: RouteOptions(),
            ),
            ConstrainedBox(
              constraints: const BoxConstraints(
                minHeight: 2.0,
              ),
            ),
          ],
        ));
  }
}
