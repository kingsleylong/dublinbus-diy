import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../models/search_form.dart';
import '../models/map_polylines.dart';
import '../views/googlemap_mobile.dart';

class RouteFavOptions extends StatefulWidget {
  const RouteFavOptions({Key? key}) : super(key: key);

  @override
  State<RouteFavOptions> createState() => FavoritePage();
}

// create the fav page
class FavoritePage extends State<RouteFavOptions> {
  @override
  Widget build(BuildContext context) {
    // init the local storage everytime creating this page to make sure it's loaded
    Provider.of<SearchFormModel>(context, listen: false).initializeFavoritesStorage();
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
  }

// function to get items from local storage
//and to make each row display the route on map
  createListItemFromLocalStorage() {
    return Consumer<SearchFormModel>(
      builder: (context, searchFormModel, child) {
        var favoriteRoutes = searchFormModel.favoriteRoutes.values.toList(growable: false);
        return ListView.builder(
            scrollDirection: Axis.vertical,
            shrinkWrap: true,
            itemCount: favoriteRoutes.length,
            itemBuilder: (context, index) {
              return InkWell(
                  child: Card(
                    child: Row(
                      children: [
                        Expanded(
                          child: Padding(
                            padding: const EdgeInsets.all(20.0),
                            child: Text(
                              favoriteRoutes[index].route.routeNumber,
                              style: const TextStyle(fontSize: 19.0),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  onTap: () {
                    print("Click to show route on map");
                    Provider.of<PolylinesModel>(context, listen: false)
                        .addBusRouteAsPolyline(favoriteRoutes[index].route);
                  });
            });
      },
    );
  }
}
