import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:web/models/map_polylines.dart';
import 'package:web/models/search_form.dart';

class RouteOptions extends StatefulWidget {
  const RouteOptions({Key? key}) : super(key: key);

  @override
  State<RouteOptions> createState() => _RouteOptionsState();
}

class _RouteOptionsState extends State<RouteOptions> {
  @override
  Widget build(BuildContext context) {
    if (Provider.of<SearchFormModel>(context).visibilityRouteOptions) {
      return Consumer<SearchFormModel>(
        builder: (context, model, child) => SingleChildScrollView(
          child: _buildRouteOptionPanels(model.busRouteItems),
        ),
      );
      //   const Expanded(
      //     child: Padding(
      //       padding: EdgeInsets.all(8),
      //       child: RouteOptions(),
      //     ),
      //   ),
      // ConstrainedBox(
      //   constraints: const BoxConstraints(
      //     minHeight: 2.0,
      //   ),
      // ),

      // return FutureBuilder<List<BusRoute>>(
      //   future: Provider.of<SearchFormModel>(context).busRoutes,
      //   builder: (context, snapshot) {
      //     if (snapshot.hasData) {
      //       return SingleChildScrollView(
      //         child: _buildRouteOptionPanels(snapshot.data!),
      //       );
      //     } else if (snapshot.hasError) {
      //       return Text('${snapshot.error}');
      //     }
      //     // By default, show a loading spinner.
      //     return const Center(
      //       child: CircularProgressIndicator(),
      //     );
      //   },
      // );
    } else {
      return Container();
    }
  }

  _buildRouteOptionPanels(List<Item> items) {
    // List<Item> items = generateItems(data);
    print("items: ${items[0].toString()}");
    // Use ExpansionPanel to display the route options for easy use.
    // https://api.flutter.dev/flutter/material/ExpansionPanel-class.html
    return ExpansionPanelList(
      expansionCallback: (int index, bool isExpanded) {
        print("isExpanded: $isExpanded");
        setState(() {
          items[index].isExpanded = !isExpanded;
          print("new isExpanded: ${items[index].isExpanded}");
        });
      },
      children: items.map<ExpansionPanel>((Item item) {
        return ExpansionPanel(
          headerBuilder: (BuildContext context, bool isExpanded) {
            return ListTile(
              title: Text(
                item.headerValue,
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
              onTap: () {
                setState(() {
                  item.isExpanded = !isExpanded;
                });
                print('item.isExpanded: ${item.isExpanded}');
                // add the polyline and marker for the selected route by changing
                // the state from the Provider and notify the Consumers.
                Provider.of<PolylinesModel>(context, listen: false)
                    .addBusRouteAsPolyline(item.busRoute);
              },
            );
          },
          body: ListTile(
            title: Text(item.expandedValue),
            subtitle: Center(
              child: Text(item.expandedDetailsValue),
            ),
          ),
          isExpanded: item.isExpanded,
        );
      }).toList(),
    );
  }
}
