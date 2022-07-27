import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:web/models/map_polylines.dart';
import 'package:web/models/search_form.dart';

class RouteOptions extends StatefulWidget {
  const RouteOptions({Key? key}) : super(key: key);

  @override
  State<RouteOptions> createState() => _RouteOptionsState();
}

class _RouteOptionsState extends State<RouteOptions> {
  @override
  Widget build(BuildContext context) {
    if (Provider.of<SearchFormModel>(context).visibilityRouteOptions) {
      return Consumer<SearchFormModel>(
        builder: (context, model, child) => SingleChildScrollView(
          child: _buildRouteOptionPanels(model.busRouteItems),
        ),
      );
      //   const Expanded(
      //     child: Padding(
      //       padding: EdgeInsets.all(8),
      //       child: RouteOptions(),
      //     ),
      //   ),
      // ConstrainedBox(
      //   constraints: const BoxConstraints(
      //     minHeight: 2.0,
      //   ),
      // ),

      // return FutureBuilder<List<BusRoute>>(
      //   future: Provider.of<SearchFormModel>(context).busRoutes,
      //   builder: (context, snapshot) {
      //     if (snapshot.hasData) {
      //       return SingleChildScrollView(
      //         child: _buildRouteOptionPanels(snapshot.data!),
      //       );
      //     } else if (snapshot.hasError) {
      //       return Text('${snapshot.error}');
      //     }
      //     // By default, show a loading spinner.
      //     return const Center(
      //       child: CircularProgressIndicator(),
      //     );
      //   },
      // );
    } else {
      return Container();
    }
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
            title: Column(
              children: [
                Text(item.expandedValue),
                const Text('Fares:'),
                Row(mainAxisAlignment: MainAxisAlignment.spaceEvenly, children: [
                  // use the Null-coalescing operators to provide an alternative value
                  // when the expression evaluates to null
                  // https://dart.dev/codelabs/null-safety#exercise-null-coalescing-operators
                  Text('Adult Leap: ${fares.adultLeap ?? '-'}'),
                  Text('Adult Cash: ${fares.adultCash ?? '-'}'),
                ]),
                Row(mainAxisAlignment: MainAxisAlignment.spaceEvenly, children: [
                  Text('Child Cash: ${fares.childCash ?? '-'}'),
                  Text('Child Leap: ${fares.childLeap ?? '-'}'),
                ]),
                Text('Student Leap: ${fares.studentLeap ?? '-'}'),
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
