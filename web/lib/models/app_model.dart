import 'package:flutter/material.dart';

import 'responsive.dart';

// This model holds the global state in the application level such as the current screen type
class AppModel extends ChangeNotifier {
  late ScreenType screenSize;
}