import 'dart:collection';
import 'dart:math';

import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';

import 'bus_route.dart';

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
    for (BusRoute busRoute in busRouteList) {
      addBusRouteAsPolyline(busRoute);
    }
    addBusRouteAsMarker(busRouteList[0]);
    // This call tells the widgets that are listening to this model to rebuild.
    notifyListeners();
  }

  void addBusRouteAsPolyline(BusRoute busRoute) {
    // remove the current polylines and markers
    removeAll();

    PolylineId polylineId = PolylineId(busRoute.routeNumber);
    // Create a Polyline by Google Map API:
    // https://github.com/flutter/plugins/blob/main/packages/google_maps_flutter/google_maps_flutter/example/lib/main.dart
    // https://developers.google.com/maps/documentation/javascript/examples/polyline-simple
    _polylines.add(Polyline(
      polylineId: polylineId,
      consumeTapEvents: true,
      color: Colors.orange,
      width: 5,
      points: _buildPoints(busRoute),
      onTap: () {
        print('tapped on ${busRoute.routeNumber}');
        // _onPolylineTapped(polylineId);
      },
    ));
    addBusRouteAsMarker(busRoute);
    notifyListeners();
  }

  void addBusRouteAsMarker(BusRoute busRoute) {
    var origin = busRoute.stops.first;
    var destination = busRoute.stops.last;
    // Create a Marker by Google Map API:
    // https://github.com/flutter/plugins/blob/main/packages/google_maps_flutter/google_maps_flutter/example/lib/place_marker.dart
    // https://developers.google.com/maps/documentation/javascript/examples/marker-simple
    _markers.add(Marker(
      // The MarkerId is the unique identifier for one marker, so it should be related to the unique
      // property of the marker
      markerId: MarkerId(origin.stopNumber),
      position: LatLng(
        origin.latitude!,
        origin.longitude!,
      ),
      infoWindow: InfoWindow(title: origin.toString()),
      onTap: () => print('tapped'),
    ));
    _markers.add(Marker(
      markerId: MarkerId(destination.stopNumber),
      position: LatLng(
        destination.latitude!,
        destination.longitude!,
      ),
      infoWindow: InfoWindow(title: destination.toString()),
      onTap: () => print('tapped'),
    ));
    // notifyListeners();
  }

  _buildPoints(BusRoute busRoute) {
    return busRoute.shapes.map<LatLng>((shape) => LatLng(shape.latitude, shape.longitude)).toList();
  }

  void showSingleMarkerPosition(double latitude, double longitude) {
    removeAll();
    _markers.add(Marker(
      markerId: MarkerId(Random().toString()),
      position: LatLng(latitude, longitude),
      infoWindow: const InfoWindow(title: 'You are here.'),
      consumeTapEvents: false,
    ));
    notifyListeners();
  }

  /// Removes all Polylines and Markers from the map.
  void removeAll() {
    _markers.clear();
    _polylines.clear();
    // This call tells the widgets that are listening to this model to rebuild.
    notifyListeners();
  }
}
