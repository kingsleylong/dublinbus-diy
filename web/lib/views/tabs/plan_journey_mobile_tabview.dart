import 'package:expandable/expandable.dart';
import 'package:flutter/material.dart';
import 'package:web/views/googlemap.dart';
import 'package:web/views/tabs/search_panel.dart';

/// This is a stateless widget because we don't need to maintain a state here.
/// It just creates the page structure and the state is managed by the imported components.
class PlanMyJourneyTabMobileView extends StatelessWidget {
  const PlanMyJourneyTabMobileView({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(
      // mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        const Expanded(child: GoogleMapComponent()),
        buildSearchFilterPanel(context),
      ],
    );
  }

  buildSearchFilterPanel(BuildContext context) {
    // Expandable widget: https://pub.dev/packages/expandable/example
    return ExpandableNotifier(
      initialExpanded: true,
      child: Column(
        children: <Widget>[
          ScrollOnExpand(
            scrollOnExpand: true,
            scrollOnCollapse: false,
            child: ExpandablePanel(
              theme: const ExpandableThemeData(
                headerAlignment: ExpandablePanelHeaderAlignment.center,
                tapBodyToCollapse: false,
              ),
              header: Center(
                  child: Text(
                "Search Filters",
                style: Theme.of(context).textTheme.titleSmall,
              )),
              collapsed: const Text(
                "loremIpsum",
                softWrap: true,
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
              expanded: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: const <Widget>[
                  SearchForm(),
                ],
              ),
              builder: (_, collapsed, expanded) {
                return Padding(
                  padding:
                      const EdgeInsets.only(left: 10, right: 10, bottom: 10),
                  child: Expandable(
                    collapsed: collapsed,
                    expanded: expanded,
                    theme: const ExpandableThemeData(crossFadePoint: 0),
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  buildSearchFilterBody() {
    return <Widget>[
      const SearchForm(),
    ];
  }
}
