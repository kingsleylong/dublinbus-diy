import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';

import 'responsive.dart';

// This model holds the global state in the application level such as the current screen type
class AppModel extends ChangeNotifier {
  late ScreenType screenSize;

  late GoogleMapController mapController;
}