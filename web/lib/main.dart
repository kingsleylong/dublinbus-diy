import 'package:flutter/material.dart';
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
    return MaterialApp(
      // Title for web page
      title: "Dublin Bus DIY",
      home: const HomePage(),
      debugShowCheckedModeBanner: true,
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