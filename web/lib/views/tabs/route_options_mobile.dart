import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../models/map_polylines.dart';
import '../../models/search_form.dart';
import '../googlemap_mobile.dart';

class RouteOptionsMobile extends StatefulWidget {
  const RouteOptionsMobile({Key? key}) : super(key: key);

  @override
  State<RouteOptionsMobile> createState() => _RouteOptionsMobileState();
}

class _RouteOptionsMobileState extends State<RouteOptionsMobile> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: const Text('Route Options'),
        ),
        body: Consumer<SearchFormModel>(
          builder: (context, model, child) => SingleChildScrollView(
            child: _buildRouteOptionPanels(model.busRouteItems),
          ),
        ));
  }

  _buildRouteOptionPanels(List<Item>? items) {
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
        var fares = item.busRoute.fares;
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
                      const Icon(Icons.directions_bus),
                      Text(
                        busRoute.routeNumber,
                        textAlign: TextAlign.center,
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                  Text(
                    busRoute.stops[0].stopName,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  Row(
                    children: [
                      const Icon(Icons.timer_outlined),
                      Text(
                        '${busRoute.travelTimes?.transitTimeMin} - ${busRoute.travelTimes?.transitTimeMax} min',
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
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
                item.isExpanded = false;
              });
              Navigator.push(
                  context,
                  MaterialPageRoute(
                      builder: (BuildContext context) => const GoogleMapMobileComponent()));
            },
            title: Column(
              children: [
                Text(item.expandedValue),
                const Text('Fares:'),
                // Use Wrap to arrange the children widgets horizontally
                // https://stackoverflow.com/a/50096780
                Wrap(spacing: 10, children: [
                  // use the Null-coalescing operators to provide an alternative value
                  // when the expression evaluates to null
                  // https://dart.dev/codelabs/null-safety#exercise-null-coalescing-operators
                  Text('Adult Leap: €${fares.adultLeap ?? '-'}'),
                  Text('Adult Cash: €${fares.adultCash ?? '-'}'),
                  Text('Child Cash: €${fares.childCash ?? '-'}'),
                  Text('Child Leap: €${fares.childLeap ?? '-'}'),
                  Text('Student Leap: €${fares.studentLeap ?? '-'}'),
                ]),
              ],
            ),
            subtitle: Center(
              child: Text(item.expandedDetailsValue),
            ),
          ),
          isExpanded: item.isExpanded,
        );
      }).toList(),
    );
  }
}
