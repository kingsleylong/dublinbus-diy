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
  // Box box = Hive.box(favRoutes);
  // bool _isFavorite = false;
  // @override
  // void initState() {
  //   super.initState();
  //   _isFavorite = box.get(0) ?? false;
  // }

  // trying to use the localstorage package
  // void initState() {
  //   super.initState();

  //   storage.ready.then((_) => printStorage());
  //   //storage.setItem("level", 0);
  //   printStorage();
  // }

  // void printStorage() {
  //   print("level stored: " + storage.getItem("level").toString());
  // }
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
    // if (box == null) {
    //   return CircularProgressIndicator();
    // } else
    //   Box box = Hive.box(favRoutes);
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
                          storage.setItem('1', busRoute.routeNumber);
//                        this saves to the array - it doesn't save to localstorage
                          // setState(() {
                          //   favoriteRouteList.add(busRoute.routeNumber);
                          //   print(favoriteRouteList);
                          // });
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
                      // elevated button to save fav route
                      // ElevatedButton(

                      //   onPressed: _addItem(busRoute.routeNumber),
                      //   style: ButtonStyle(
                      //     backgroundColor: MaterialStateProperty.all<Color>(
                      //       Colors.white,
                      //     ),
                      //   ),
                      //   child: const Icon(
                      //     Icons.favorite,
                      //     color: Colors.red,
                      //   ),
                      // ),
                      // {
                      //   // print(favoriteRouteList);
                      //   setState(() {
                      //     // favoriteRouteList.add(busRoute.routeNumber);
                      //     if (!initialized) {
                      //       var items = storage.getItem('favourite');

                      //       if (items != null) {
                      //         list.favoriteRouteList = List<RouteItem>.from(
                      //           (items as List).map(
                      //             (item) => RouteItem(
                      //               route: item['Route'],
                      //               favourite: item['favourite_route'],
                      //             ),
                      //           ),
                      //         );
                      //       }

                      //       initialized = true;
                      //     }
                      //     List<Widget> widgets =
                      //         list.favoriteRouteList.map((item) {
                      //       return CheckboxListTile(
                      //         value: item.favourite,
                      //         title: Text(item.route),
                      //         selected: item.favourite,
                      //         onChanged: (_) {
                      //           _toggleItem(item);
                      //         },
                      //       );
                      //     }).toList();

                      //     // print(favoriteRouteList);
                      //   });
                      // },

                      // icon button to save fav
                      // IconButton(
                      //   icon: Icon(Icons.favorite_border),
                      //   onPressed: printStorage,
                      //   color: Colors.red,)
                      //  () {
                      //   // setState(() {
                      //   //   _isFavorite = !_isFavorite;
                      //   // });

                      //   // box.put(0, _isFavorite);
                      // },

                      // ElevatedButton(
                      //   onPressed: () {},
                      //   child: Icon(Icons.favorite_border),
                      // )
                      // FavoriteButton(
                      //   isFavorite: false,
                      //   iconSize: 10,
                      //   valueChanged: (_isFavorite) {
                      //     print('Is Favorite : $_isFavorite');
                      //   },
                      // ),
                      // ],
                      // ),
                      // Text(
                      //   busRoute.stops[0].stopName,
                      //   style: const TextStyle(
                      //     fontSize: 16,
                      //     fontWeight: FontWeight.bold,
                      //   ),
                      // ),
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

  // void addItemsToLocalStorage() {
  //   final item = json.encode(RouteItem(route: , favourite: false));
  //   storage.setItem('info', item);
  // }
}

// class RouteFavOptions extends StatefulWidget {
//   const RouteFavOptions({Key? key}) : super(key: key);

//   @override
//   State<RouteFavOptions> createState() => FavoritePage();
// }

// // create the fav page
// class FavoritePage extends State<RouteFavOptions> {
//   @override
//   Widget build(BuildContext context) {
//     return Scaffold(
//         appBar: AppBar(
//           title: const Text("Favourites"),
//           centerTitle: true,
//           backgroundColor: Colors.blue,
//         ),
//         drawer: new Drawer(),
//         body: TabBarView(
//           children: [
//             favoriteRouteList.isEmpty
//                 ? const Center(
//                     child: Text(
//                       'There are no favorites yet!',
//                       style: TextStyle(color: Colors.black),
//                     ),
//                   )
//                 : ListView.builder(
//                     itemCount: favoriteRouteList.length,
//                     itemBuilder: (context, index) {
//                       return Card(
//                         child: Row(
//                           children: [
//                             Expanded(
//                               child: Padding(
//                                 padding: const EdgeInsets.all(20.0),
//                                 child: Text(
//                                   favoriteRouteList[index],
//                                   style: const TextStyle(fontSize: 19.0),
//                                 ),
//                               ),
//                             ),
//                             ElevatedButton(
//                               onPressed: () {
//                                 setState(() {
//                                   favoriteRouteList
//                                       .remove(favoriteRouteList[index]);
//                                 });
//                               },
//                               style: ButtonStyle(
//                                 backgroundColor:
//                                     MaterialStateProperty.all<Color>(
//                                   Colors.deepPurple,
//                                 ),
//                               ),
//                               child: const Icon(
//                                 Icons.remove,
//                                 color: Colors.white,
//                               ),
//                             ),
//                           ],
//                         ),
//                       );
//                     },
//                   ),
//           ],
//         ));
//   }
// }

// this is for hivebox ,
// drawer: new Drawer(),
//       body: ValueListenableBuilder(
//         valueListenable: Hive.box(favRoutes).listenable(),
//         builder: (context, box, child) {
//           List posts = List.from(box.values);
//           return ListView(
//             padding: const EdgeInsets.all(16.0),
//             children: [
//               Text("This is fav page"),
//               ...posts.map(
//                 (p) => ListTile(
//                   title: Text(p['title']),
//                   trailing: IconButton(
//                     icon: Icon(
//                       Hive.box(favRoutes).containsKey(p['id'])
//                           ? Icons.favorite
//                           : Icons.favorite_border,
//                     ),
//                     onPressed: () {
//                       if (Hive.box(favRoutes).containsKey(p['id'])) {
//                         Hive.box(favRoutes).delete(p['id']);
//                       } else {
//                         Hive.box(favRoutes).put(p['id'], p);
//                       }
//                     },
//                   ),
//                 ),
//               ),
//             ],
//           );
//         },
//       ),
//     );
//   }
// }
