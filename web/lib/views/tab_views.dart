import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/responsive.dart';
import '../models/search_form.dart';
import 'tabs/route_options.dart';
import 'tabs/search_panel.dart';

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
            const SearchForm(screenSize: ScreenType.desktop),
            Expanded(child: showRouteOptionsOrLoadingIcon(context)),
          ],
        ));
    // ;
  }

  Widget showRouteOptionsOrLoadingIcon(BuildContext context) {
    if (Provider.of<SearchFormModel>(context).visibilityRouteOptions) {
      return const RouteOptions();
    } else if (Provider.of<SearchFormModel>(context).visibilityLoadingIcon) {
      return const Center(child: CircularProgressIndicator());
    } else {
      return Container();
    }
  }
}
