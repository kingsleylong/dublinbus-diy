import 'package:dublin_bus_diy/models/responsive.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../models/app_model.dart';
import '../../models/bus_route.dart';
import '../../models/map_polylines.dart';
import '../../models/search_form.dart';
import '../googlemap_mobile.dart';
import 'fares_table.dart';

class RouteOptions extends StatefulWidget {
  const RouteOptions({Key? key}) : super(key: key);

  @override
  State<RouteOptions> createState() => _RouteOptionsState();
}

class _RouteOptionsState extends State<RouteOptions> {
  @override
  Widget build(BuildContext context) {
    return Consumer<SearchFormModel>(
      builder: (context, model, child) =>
          Provider.of<SearchFormModel>(context).visibilityRouteOptions
              ? SingleChildScrollView(child: _buildRouteOptionPanels(model.busRouteItems))
              : const Center(child: CircularProgressIndicator()),
    );
  }

  _buildRouteOptionPanels(List<Item>? items) {
    final TextTheme textTheme = Theme.of(context).textTheme;
    if (items == null || items.isEmpty) {
      return const Center(child: Text('No routes found.'));
    }
    print("items size: ${items.length}, first ele: ${items[0].toString()}");
    // Use ExpansionPanel to display the route options for easy use.
    // https://api.flutter.dev/flutter/material/ExpansionPanel-class.html
    return ExpansionPanelList(
      expansionCallback: (int index, bool isExpanded) {
        print("isExpanded: $isExpanded");
        setState(() {
          items[index].isExpanded = !isExpanded;
          print("new isExpanded: ${items[index].isExpanded}");
        });
        // add the polyline and marker for the selected route by changing the
        // state from the Provider and notify the Consumers.
        Provider.of<PolylinesModel>(context, listen: false)
            .addBusRouteAsPolyline(items[index].busRoute);
      },
      children: items.map<ExpansionPanel>((Item item) {
        var busRoute = item.busRoute;
        var fares = busRoute.fares;
        var travelTimes = busRoute.travelTimes;
        return ExpansionPanel(
          canTapOnHeader: true,
          headerBuilder: (BuildContext context, bool isExpanded) {
            return ListTile(
              title: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.start,
                    children: [
                      const Padding(
                        padding: EdgeInsets.only(right: 4),
                        child: Icon(Icons.directions_bus),
                      ),
                      Container(
                        color: Colors.amberAccent,
                        child: Padding(
                          padding: const EdgeInsets.all(3.0),
                          child: Text(
                            busRoute.routeNumber,
                            textAlign: TextAlign.center,
                            style: const TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                      ),
                      const Padding(
                        padding: EdgeInsets.only(right: 7),
                      ),
                      Text(busRoute.travelTimes.scheduledDepartureTime),
                    ],
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.start,
                    children: [
                      Padding(
                          padding: const EdgeInsets.only(right: 4),
                          child: travelTimes.source == TravelTimeSources.static
                              // tooltip: https://api.flutter.dev/flutter/material/Tooltip-class.html
                              // travel time from static table
                              ? const Tooltip(
                                  message: 'Travel time from static time table',
                                  child: Icon(Icons.timer_outlined),
                                )
                              // travel time from prediction
                              : Tooltip(
                                  message: 'Predicted travel time: ${travelTimes.transitTimeMin}'
                                      ' - ${travelTimes.transitTimeMax} min',
                                  child: const Icon(Icons.update),
                                )),
                      // sized box sets a fixed width of the text and align them vertically
                      SizedBox(
                        width: 60,
                        child: Text(
                          // '${busRoute.travelTimes?.transitTimeMin} - ${busRoute.travelTimes?.transitTimeMax} min',
                          '${travelTimes.transitTime} min',
                          style: const TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            );
          },
          body: ListTile(
            onTap: () {
              setState(() {
                // always collapse the panel on tapping the body
                item.isExpanded = false;
              });
              if (Provider.of<AppModel>(context, listen: false).screenSize == ScreenType.mobile) {
                Navigator.push(
                    context,
                    MaterialPageRoute(
                        builder: (BuildContext context) => const GoogleMapMobileComponent()));
              }
            },
            title: Column(
              children: [
                Text(item.expandedValue),
                FaresTable(fares: fares),
              ],
            ),
            subtitle: Column(
              children: [
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Text(
                    'Touch to see the route on map',
                    style: textTheme.bodyLarge,
                  ),
                ),
                Text(item.expandedDetailsValue),
              ],
            ),
          ),
          isExpanded: item.isExpanded,
        );
      }).toList(),
    );
  }
}
