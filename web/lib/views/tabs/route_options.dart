import 'package:dublin_bus_diy/models/responsive.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../models/app_model.dart';
import '../../models/bus_route.dart';
import '../../models/map_polylines.dart';
import '../../models/search_form.dart';
import '../googlemap_mobile.dart';
import 'fares_table.dart';
import 'package:localstorage/localstorage.dart';

List<String> favoriteRouteList = [];

// Map<String, dynamic> favoriteRouteList = {};

class RouteItem {
  String route;
  bool favourite;

  RouteItem({required this.route, required this.favourite});

  toJSONEncodable() {
    Map<String, dynamic> m = new Map();

    m['route_num'] = route;
    m['favourite'] = favourite;
    return m;
  }
}

class RouteList {
  List<RouteItem> favoriteRouteList = [];

  toJSONEncodable() {
    return favoriteRouteList.map((item) {
      return item.toJSONEncodable();
    }).toList();
  }
}

class RouteOptions extends StatefulWidget {
  const RouteOptions({Key? key}) : super(key: key);

  @override
  State<RouteOptions> createState() => _RouteOptionsState();
}

class _RouteOptionsState extends State<RouteOptions> {
  final RouteList list = new RouteList();
  final LocalStorage storage = new LocalStorage('fav_routes');
  bool initialized = false;

  _toggleItem(RouteItem route) {
    setState(() {
      route.favourite = !route.favourite;
      _saveToStorage();
    });
  }

  _addItem(String route) {
    setState(() {
      final item = new RouteItem(route: route, favourite: false);
      list.favoriteRouteList.add(item);
      _saveToStorage();
    });
  }

  _deleteItem(String route) {
    setState(() {
      final item = new RouteItem(route: route, favourite: false);
      list.favoriteRouteList.remove(item);
      _deleteFromStorage();
    });
  }

  _saveToStorage() {
    storage.setItem('favourite', list.toJSONEncodable());
  }

  _deleteFromStorage() {
    storage.deleteItem('favourite');
  }

  _clearStorage() async {
    await storage.clear();

    setState(() {
      list.favoriteRouteList = storage.getItem('favourite') ?? [];
    });
  }

  @override
  Widget build(BuildContext context) {
    if (Provider.of<SearchFormModel>(context).visibilityRouteOptions) {
      return Consumer<SearchFormModel>(
        builder: (context, model, child) => SingleChildScrollView(
          child: _buildRouteOptionPanels(model.busRouteItems),
        ),
      );
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
        // print(busRoute.routeNumber);
        var fares = item.busRoute.fares;
        return ExpansionPanel(
          canTapOnHeader: true,
          headerBuilder: (BuildContext context, bool isExpanded) {
            return ListTile(
              // TODO: here need to add the fav icon button and create the list
              // create the list for the fav route
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
                      ElevatedButton(
                        onPressed: () {
                          // this should save the route to localstorage
                          storage.ready;
                          print("setting the route as favourite");
                          storage.setItem('=', busRoute.routeNumber);
                        },
                        style: ButtonStyle(
                          backgroundColor: MaterialStateProperty.all<Color>(
                            Colors.white,
                          ),
                        ),
                        child: const Icon(
                          Icons.favorite,
                          color: Colors.red,
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
                ],
              ),
            );
          },
          body: ListTile(
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
