import 'dart:collection';

import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:web/models/bus_route.dart';

// We use the ChangeNotifier to manage the state of the Models
// https://docs.flutter.dev/development/data-and-backend/state-mgmt/simple#changenotifier
class MarkersModel extends ChangeNotifier {
  /// Internal, private state of the polylines.
  final List<Marker> _markers = [];
  /// An unmodifiable view of the polylines.
  UnmodifiableListView<Marker> get items => UnmodifiableListView(_markers);

  /// Adds [marker] to the map. This and [removeAll] are the only ways to modify the
  /// cart from the outside.
  void add(Marker marker) {
    _markers.add(marker);
    // This call tells the widgets that are listening to this model to rebuild.
    notifyListeners();
  }

  void addBusRouteAsMarker(BusRoute busRoute) {
    MarkerId markerId = MarkerId('$busRoute.hashCode');
    // Create a Marker by Google Map API:
    // https://github.com/flutter/plugins/blob/main/packages/google_maps_flutter/google_maps_flutter/example/lib/place_marker.dart
    // https://developers.google.com/maps/documentation/javascript/examples/marker-simple
    _markers.add(
        Marker(
          markerId: markerId,
          position: LatLng(
            busRoute.stops[0].latitude!,
            busRoute.stops[0].longitude!,
          ),
          infoWindow: InfoWindow(title: '${busRoute.stops[0].stopName} - ${busRoute.stops[0].stopNumber}'),
          onTap: () => print('tapped'),
        )
    );
    notifyListeners();
  }

  /// Removes all Markers from the map.
  void removeAll() {
    _markers.clear();
    // This call tells the widgets that are listening to this model to rebuild.
    notifyListeners();
  }
}