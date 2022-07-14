import 'dart:collection';

import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:web/models/bus_route.dart';

// We use the ChangeNotifier to manage the state of the Models
// https://docs.flutter.dev/development/data-and-backend/state-mgmt/simple#changenotifier
class PolylinesModel extends ChangeNotifier {
  /// Internal, private state of the polylines.
  final List<Polyline> _polylines = [];
  /// An unmodifiable view of the polylines.
  UnmodifiableListView<Polyline> get items => UnmodifiableListView(_polylines);

  /// Adds [polyline] to the map. This and [removeAll] are the only ways to modify the
  /// cart from the outside.
  void add(Polyline polyline) {
    _polylines.add(polyline);
    // This call tells the widgets that are listening to this model to rebuild.
    notifyListeners();
  }

  void addBusRouteListAsPolylines(List<BusRoute> busRouteList) {
    busRouteList.map((busRoute) => (){
      addBusRouteAsPolyline(busRoute);
    });
  }

  void addBusRouteAsPolyline(BusRoute busRoute) {
    PolylineId polylineId = PolylineId('$busRoute.hashCode');
    // Create a Polyline by Google Map API:
    // https://github.com/flutter/plugins/blob/main/packages/google_maps_flutter/google_maps_flutter/example/lib/main.dart
    // https://developers.google.com/maps/documentation/javascript/examples/polyline-simple
    _polylines.add(
        Polyline(
          polylineId: polylineId,
          consumeTapEvents: true,
          color: Colors.orange,
          width: 5,
          points: _buildPoints(busRoute),
          onTap: () {
            print('tapped on ${busRoute.routeNumber}');
          // _onPolylineTapped(polylineId);
          },
        )
    );
    notifyListeners();
  }

  _buildPoints(BusRoute busRoute) {
    return busRoute.shapes.map<LatLng>((shape) =>
        LatLng(shape.latitude, shape.longitude)
    ).toList();
  }

  /// Removes all Polylines from the map.
  void removeAll() {
    _polylines.clear();
    // This call tells the widgets that are listening to this model to rebuild.
    notifyListeners();
  }
}