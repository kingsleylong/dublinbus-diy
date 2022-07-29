import 'package:flutter/material.dart';

class ResponsiveLayout extends StatelessWidget {
  const ResponsiveLayout({Key? key, required this.mobileBody, required this.desktopBody})
      : super(key: key);

  final Widget mobileBody;
  final Widget desktopBody;

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        // The most common mobile screen sizes for 2021
        // https://worship.agency/mobile-screen-sizes-for-2021
        // https://www.youtube.com/watch?v=MrPJBAOzKTQ&list=PLn3LDx3baxQHtXo8_5p1KCB5MLvbdonSU&index=1&t=76s
        if (constraints.maxWidth < 650) {
          return mobileBody;
        } else {
          return desktopBody;
        }
      },
    );
  }
}
