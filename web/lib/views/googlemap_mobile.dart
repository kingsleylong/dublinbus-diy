import 'package:flutter/material.dart';
import 'package:web/views/googlemap.dart';

class GoogleMapMobileComponent extends StatelessWidget {
  const GoogleMapMobileComponent({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: const Text('Route Map'),
        ),
        body: const GoogleMapComponent());
  }
}
