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
          ],
        ));
    //     ],
    //   ),
    // );
  }

// function to get items from local storage
//and to make each row display the route on map
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
            // TODO: need to get all documents from storage
            // could do a for loop to iterate through all items in storage
            var items = storage.getItem('9');

            if (items != null) {
              // convert into list
              favoriteRouteList.add(items);
            }

            initialized = true;
          }

          return ListView.builder(
              scrollDirection: Axis.vertical,
              shrinkWrap: true,
              itemCount: favoriteRouteList.length,
              itemBuilder: (context, index) {
                return InkWell(
                    child: Card(
                      child: Row(
                        children: [
                          Expanded(
                            child: Padding(
                              padding: const EdgeInsets.all(20.0),
                              child: Text(
                                favoriteRouteList[index],
                                style: const TextStyle(fontSize: 19.0),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                    onTap: () {
                      // TODO: error that says "String" can't  be assigned to the parameter
                      // type 'BusRoute'.
                      print("Click to show route on map");
                      // Provider.of<PolylinesModel>(context, listen: false)
                      //     .addBusRouteAsPolyline(favoriteRouteList[index]);
                    });
              });
        });
  }
}
