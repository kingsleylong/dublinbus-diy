import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:web/models/map_polylines.dart';
import 'package:web/models/search_form.dart';
import 'package:web/views/home_page.dart';

void main() => runApp(const DublinBusDiyApp());

class DublinBusDiyApp extends StatefulWidget {
  const DublinBusDiyApp({Key? key}) : super(key: key);

  @override
  _DublinBusDiyAppState createState() => _DublinBusDiyAppState();
}

class _DublinBusDiyAppState extends State<DublinBusDiyApp> {
  @override
  Widget build(BuildContext context) {
    // Move the ChangeNotifierProviders to the APP level so that all the routes can share them
    // https://stackoverflow.com/a/66269538
    return MultiProvider(
      providers: [
        // Create a model by the provider so the child can listen to the model changes
        // https://docs.flutter.dev/development/data-and-backend/state-mgmt/simple#changenotifierprovider
        ChangeNotifierProvider<PolylinesModel>(create: (context) => PolylinesModel()),
        ChangeNotifierProvider<SearchFormModel>(create: (context) => SearchFormModel()),
      ],
      child: const MaterialApp(
        // Title for web page
        title: "Dublin Bus DIY",
        home: HomePage(),
        debugShowCheckedModeBanner: true,
      ),
    );
  }

  ThemeData buildThemeData() {
    return ThemeData(
      colorScheme: ColorScheme.fromSwatch(
        primarySwatch: Colors.green,
      ),
      scaffoldBackgroundColor: Colors.green[100],
    );
  }
}
