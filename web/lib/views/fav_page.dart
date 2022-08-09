import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:hive/hive.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:hive_flutter/adapters.dart';
import 'package:http/http.dart';
import 'package:provider/provider.dart';
import '../../models/app_model.dart';
import '../../models/bus_route.dart';
import '../../models/bus_route.dart';
import '../../models/map_polylines.dart';
import '../../models/search_form.dart';
import '../views/googlemap_mobile.dart';
import '../views/tabs/fares_table.dart';
import 'package:localstorage/localstorage.dart';

import 'tabs/route_options.dart';

class RouteFavOptions extends StatefulWidget {
  const RouteFavOptions({Key? key}) : super(key: key);

  @override
  State<RouteFavOptions> createState() => FavoritePage();
}

// create the fav page
class FavoritePage extends State<RouteFavOptions> {
  final RouteList list = new RouteList();
  final LocalStorage storage = new LocalStorage('fav_routes');
  bool initialized = false;

  _toggleItem(RouteItem route) {
    setState(() {
      route.favourite = !route.favourite;
      _saveToStorage();
    });
  }

  _saveToStorage() {
    storage.setItem('favourite', list.toJSONEncodable());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: const Text("Favourites"),
          centerTitle: true,
          backgroundColor: Colors.blue,
        ),
        body: Row(
          children: [
            Expanded(child: createListItemFromLocalStorage()),
            const Expanded(child: GoogleMapMobileComponent()),
            // Column(
            //   Container(
            // padding: EdgeInsets.all(100),
            // constraints: BoxConstraints.expand(),

            // child: FutureBuilder(
            //   future: storage.ready,
            //   builder: (BuildContext context, AsyncSnapshot snapshot) {
            //     if (snapshot.data == null) {
            //       return Center(
            //         child: CircularProgressIndicator(),
            //       );
            //     }

            // children:
            // [
            // if (!initialized) {
            //   var items = storage.getItem('favourite');

            //   if (items != null) {
            //     list.favoriteRouteList = List<RouteItem>.from(
            //       (items as List).map(
            //         (item) => RouteItem(
            //           route: item['Route'],
            //           favourite: item['favourite_route'],
            //         ),
            //       ),
            //     );
            //   }

            //   initialized = true;
            // }
            // List<Widget> widgets = list.favoriteRouteList.map((item) {
            //   return CheckboxListTile(
            //     value: item.favourite,
            //     title: Text(item.route),
            //     selected: item.favourite,
            //     onChanged: (_) {
            //       _toggleItem(item);
            //     },
            //   );
            // }).toList();

            // return Text("No fav route");
            Expanded(
              child: ListView(
                padding: EdgeInsets.only(top: 50),
                children: [
                  favoriteRouteList.isEmpty
                      ? const Center(
                          child: Text(
                            'There are no favorites routes!',
                            style: TextStyle(color: Colors.black),
                          ),
                        )
                      : ListView.builder(
                          scrollDirection: Axis.vertical,
                          shrinkWrap: true,
                          itemCount: favoriteRouteList.length,
                          itemBuilder: (context, index) {
                            return Card(
                              child: Row(
                                children: [
                                  Expanded(
                                    child: Padding(
                                      padding: const EdgeInsets.all(20.0),
                                      child: Text(
                                        // storage.ready;
                                        // storage.getItem('1'),
                                        favoriteRouteList[index],
                                        style: const TextStyle(fontSize: 19.0),
                                      ),
                                    ),
                                  ),
                                  ElevatedButton(
                                    onPressed: () {
                                      storage.ready;
                                      print("deleting the route from favourite");
                                      storage.deleteItem('1');
                                      // this deletes from the array - it doesn't save to localstorage
                                      // setState(() {
                                      //   favoriteRouteList
                                      //       .remove(favoriteRouteList[index]);
                                      // });
                                    },
                                    style: ButtonStyle(
                                      backgroundColor: MaterialStateProperty.all<Color>(
                                        Colors.white,
                                      ),
                                    ),
                                    child: const Icon(
                                      Icons.remove_circle,
                                      color: Colors.red,
                                    ),
                                  ),
                                ],
                              ),
                            );
                          },
                        ),
                ],
              ),
            ),
          ],
        ));
    //     ],
    //   ),
    // );
  }

  createListItemFromLocalStorage() {
    return FutureBuilder(
        future: storage.ready,
        builder: (BuildContext context, AsyncSnapshot snapshot) {
          if (snapshot.data == null) {
            return Center(
              child: CircularProgressIndicator(),
            );
          }

          if (!initialized) {
            var items = storage.getItem('1');

            if (items != null) {
              // convert into list
              favoriteRouteList.add(items);
              // list.items = List<TodoItem>.from(
              //   (items as List).map(
              //         (item) =>
              //         TodoItem(
              //           title: item['title'],
              //           done: item['done'],
              //         ),
              //   ),
              // );
            }
            // favoriteRouteList = [];

            initialized = true;
          }

          return ListView(
            // favoriteRouteList
            children: [
              Text(favoriteRouteList[0]),
            ],
          );
        });
  }
}
// ];
// },
// );
// ),
// );
// }
// }
