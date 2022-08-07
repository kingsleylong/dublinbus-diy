import 'package:flutter/material.dart';

import 'googlemap.dart';
import 'tab_views.dart';
import 'about_us.dart';

class DesktopBody extends StatefulWidget {
  const DesktopBody({Key? key, required this.tabController}) : super(key: key);
  final TabController tabController;

  @override
  State<DesktopBody> createState() => _DesktopBodyState();
}

class _DesktopBodyState extends State<DesktopBody> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: const Text("Dublin Bus DIY"),
          actions: <Widget>[
            ElevatedButton(
              child: const Text('About Us'),
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => const AboutUs()),
                );
              },
            ),
            // ElevatedButton(
            //   child: const Text('Favourites'),
            //   onPressed: () {
            //     Navigator.push(
            //       context,
            //       MaterialPageRoute(builder: (context) => FavoritePage()),
            //     );
            //   },
            // ),
            // IconButton(icon: Icon(Icons.favorite), onPressed: () {}),
          ],
        ),
        body: Row(
          // alignment: Alignment.topLeft,
          children: [
            // left bar
            buildLeftBar(),
            // right information box, use Expanded class to take the rest of space
            // const Expanded(
            //   child: GoogleMapComponent(),
            // )
            // buildLeftTabViews(),

            Expanded(
              child: buildRightInformationBox(),
            ),
          ],
        ));
  }

  TabBarView buildRightInformationBox() {
    print('Build right information box');
    return TabBarView(
        controller: widget.tabController,
        // disable swiping from TabBarView
        // https://flutteragency.com/how-to-disable-swipe-tabbar-in-flutter/
        physics: const NeverScrollableScrollPhysics(),
        children: const <Widget>[
          GoogleMapComponent(),
        ]);
  }

  TabBarView buildLeftTabViews() {
    return TabBarView(
        controller: widget.tabController,
        physics: const NeverScrollableScrollPhysics(),
        children: const <Widget>[
          PlanMyJourneyTabView(),
        ]);
  }

  SizedBox buildLeftBar() {
    return SizedBox(
      width: 350,
      child: Column(
        children: [
          // ColoredBox(
          //   color: Colors.red,
          //   child: TabBar(
          //     // expand the tab bar out of range and slide the bar when clicking
          //     // tabs at the edges https://stackoverflow.com/a/60636918
          //     // isScrollable: true,
          //     // Access a field of the widget in its state https://stackoverflow.com/a/58767810
          //     controller: widget.tabController,
          //     tabs: tabList,
          //   ),
          // ),
          //
          // // tab views
          Expanded(
            child: buildLeftTabViews(),
          ),

          //weather info
          // TODO
        ],
      ),
    );
  }
}

// change the homepage to the about us page
void onItemPressed(BuildContext context, {required int index}) {
  Navigator.pop(context);

  switch (index) {
    case 0:
      Navigator.push(
          context, MaterialPageRoute(builder: (context) => const AboutUs()));
      break;
    default:
      Navigator.pop(context);
      break;
  }
}

// change the homepage to the favourite routes page
// void onItemPressed2(BuildContext context, {required int index}) {
//   Navigator.pop(context);

//   switch (index) {
//     case 0:
//       Navigator.push(
//           context, MaterialPageRoute(builder: (context) => FavoritePage()));
//       break;
//     default:
//       Navigator.pop(context);
//       break;
//   }
// }
// // }
