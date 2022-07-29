import 'package:flutter/material.dart';
import 'package:package_info_plus/package_info_plus.dart';

class ResponsiveLayout extends StatelessWidget {
  const ResponsiveLayout({Key? key, required this.mobileBody, required this.desktopBody})
      : super(key: key);

  final Widget mobileBody;
  final Widget desktopBody;

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        PackageInfo.fromPlatform().then((PackageInfo packageInfo) {
          String appName = packageInfo.appName;
          String packageName = packageInfo.packageName;
          String version = packageInfo.version;
          String buildNumber = packageInfo.buildNumber;
          print('$packageName $appName: $version+$buildNumber');
        });
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
