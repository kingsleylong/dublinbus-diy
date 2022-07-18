import 'dart:collection';

import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:web/models/bus_route.dart';

// We use the ChangeNotifier to manage the state of the Models
// https://docs.flutter.dev/development/data-and-backend/state-mgmt/simple#changenotifier
class PolylinesModel extends ChangeNotifier {
  /// Internal, private state of the polylines.
  final List<Polyline> _polylines = [];
  /// Internal, private state of the markers.
  final List<Marker> _markers = [];
  /// An unmodifiable view of the polylines.
  UnmodifiableListView<Polyline> get itemsOfPolylines => UnmodifiableListView(_polylines);
  /// An unmodifiable view of the markers.
  UnmodifiableListView<Marker> get itemsOfMarkers => UnmodifiableListView(_markers);

  /// Adds [polyline] to the map. This and [removeAll] are the only ways to modify the
  /// cart from the outside.
  void add(Polyline polyline) {
    _polylines.add(polyline);
    // This call tells the widgets that are listening to this model to rebuild.
    notifyListeners();
  }

  void addBusRouteListAsPolylines(List<BusRoute> busRouteList) {
    for(BusRoute busRoute in busRouteList) {
      addBusRouteAsPolyline(busRoute);
    }
    addBusRouteAsMarker(busRouteList[0]);
    // This call tells the widgets that are listening to this model to rebuild.
    notifyListeners();
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
    // notifyListeners();
  }

  void addBusRouteAsMarker(BusRoute busRoute) {
    MarkerId originMarkerId = const MarkerId('origin');
    // Create a Marker by Google Map API:
    // https://github.com/flutter/plugins/blob/main/packages/google_maps_flutter/google_maps_flutter/example/lib/place_marker.dart
    // https://developers.google.com/maps/documentation/javascript/examples/marker-simple
    _markers.add(
        Marker(
          markerId: originMarkerId,
          position: LatLng(
            busRoute.stops[0].latitude!,
            busRoute.stops[0].longitude!,
          ),
          infoWindow: InfoWindow(title: '${busRoute.stops[0].stopName}'
              ' - ${busRoute.stops[0].stopNumber}'),
          onTap: () => print('tapped'),
        )
    );
    MarkerId destinationMarkerId = const MarkerId('destination');
    _markers.add(
        Marker(
          markerId: destinationMarkerId,
          position: LatLng(
            busRoute.stops[busRoute.stops.length - 1].latitude!,
            busRoute.stops[busRoute.stops.length - 1].longitude!,
          ),
          infoWindow: InfoWindow(title: '${busRoute.stops[busRoute.stops.length - 1].stopName}'
              ' - ${busRoute.stops[busRoute.stops.length - 1].stopNumber}'),
          onTap: () => print('tapped'),
        )
    );
    // notifyListeners();
  }

  _buildPoints(BusRoute busRoute) {
    // final List<LatLng> points = <LatLng>[];
    return busRoute.stops.map<LatLng>((stop) =>
        LatLng(stop.latitude!, stop.longitude!)
    ).toList();
  }

  /// Removes all Polylines and Markers from the map.
  void removeAll() {
    _markers.clear();
    _polylines.clear();
    // This call tells the widgets that are listening to this model to rebuild.
    notifyListeners();
  }
}