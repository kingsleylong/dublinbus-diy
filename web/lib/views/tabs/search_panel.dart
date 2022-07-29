import 'package:date_time_picker/date_time_picker.dart';
import 'package:dropdown_search/dropdown_search.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:toggle_switch/toggle_switch.dart';

import '../../api/fetch_bus_stop.dart';
import '../../models/bus_route.dart';
import '../../models/bus_route_filter.dart';
import '../../models/bus_stop.dart';
import '../../models/map_polylines.dart';
import '../../models/responsive.dart';
import '../../models/search_form.dart';
import 'route_options_mobile.dart';

class SearchForm extends StatefulWidget {
  const SearchForm({Key? key, required this.screenSize}) : super(key: key);

  final ScreenSize screenSize;

  @override
  State<SearchForm> createState() => _SearchFormState();
}

class _SearchFormState extends State<SearchForm> {
  late Future<List<BusRoute>> futureBusRoutes;

  // late List<Item> items;

  @override
  Widget build(BuildContext context) {
    return Padding(
      // padding settings https://api.flutter.dev/flutter/material/InputDecoration/contentPadding.html
      padding: const EdgeInsets.fromLTRB(5, 10, 5, 0),
      // Listen to the SearchFormModel and share the state of the form between multiple views
      child: Consumer<SearchFormModel>(
        builder: (context, searchFormModel, child) => Form(
            key: searchFormModel.formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: <Widget>[
                // Origin dropdown list
                // The form fields should be wrapped by Padding otherwise they would overlap each other
                // https://docs.flutter.dev/cookbook/forms/text-input
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
                  child: buildSearchableOriginDropdownList(searchFormModel),
                ),
                // Destination dropdown list
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
                  child: buildSearchableDestinationDropdownList(searchFormModel),
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
                  child: ToggleSwitch(
                    // https://pub.dev/packages/toggle_switch
                    // Here, default theme colors are used for activeBgColor, activeFgColor, inactiveBgColor and inactiveFgColor
                    minWidth: 150,
                    inactiveBgColor: Colors.grey[300],
                    initialLabelIndex: searchFormModel.timeTypeToggleIndex,
                    totalSwitches: 2,
                    labels: const ['Departure', 'Arrival'],
                    onToggle: (index) {
                      searchFormModel.timeTypeToggleIndex = index;
                      print('switched to: ${searchFormModel.timeTypeToggleIndex}');
                    },
                  ),
                ),
                // Departure/Arrival time
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
                  child: DateTimePicker(
                    // Data time picker: https://pub.dev/packages/date_time_picker
                    type: DateTimePickerType.dateTimeSeparate,
                    // Date format: https://api.flutter.dev/flutter/intl/DateFormat-class.html
                    dateMask: 'E d MMM, yyyy',
                    controller: searchFormModel.dateTimePickerController,
                    firstDate: DateTime.now(),
                    // We allow travel planning ahead of 4 days
                    lastDate: DateTime.now().add(const Duration(hours: 4 * 24)),
                    icon: const Icon(Icons.event),
                    dateLabelText: 'Date',
                    timeLabelText: "Hour",
                    onChanged: (val) => print(val),
                  ),
                ),
                // Submit button
                Padding(
                  padding: const EdgeInsets.all(8),
                  child: ElevatedButton(
                    style: ElevatedButton.styleFrom(
                      minimumSize: const Size.fromHeight(60),
                      textStyle: const TextStyle(fontSize: 18),
                    ),
                    onPressed: () {
                      // Validate will return true if the form is valid, or false if
                      // the form is invalid.
                      if (searchFormModel.formKey.currentState!.validate()) {
                        Provider.of<PolylinesModel>(context, listen: false).removeAll();
                        print(
                            'Selected origin: ${searchFormModel.originSelectionKey.currentState?.getSelectedItem?.stopNumber}');
                        print(
                            'Selected destination: ${searchFormModel.destinationSelectionKey.currentState?.getSelectedItem?.stopNumber}');
                        print(
                            'Selected datetime: ${searchFormModel.dateTimePickerController.value.text}');

                        DateTime parseTime =
                            DateTime.parse(searchFormModel.dateTimePickerController.value.text);

                        // Date format: https://api.flutter.dev/flutter/intl/DateFormat-class.html
                        print(
                            'parseTime: $parseTime  ${DateFormat('MM-dd-yyyy HH:mm:ss').format(parseTime)}');

                        BusRouteSearchFilter searchFilter = BusRouteSearchFilter(
                            searchFormModel
                                .originSelectionKey.currentState?.getSelectedItem?.stopNumber,
                            searchFormModel
                                .destinationSelectionKey.currentState?.getSelectedItem?.stopNumber,
                            searchFormModel.timeTypes[searchFormModel.timeTypeToggleIndex ?? 0],
                            DateFormat('MM-dd-yyyy HH:mm:ss').format(parseTime));

                        // futureBusRoutes = fetchBusRoutes(searchFilter);
                        Provider.of<SearchFormModel>(context, listen: false)
                            .fetchBusRoute(searchFilter);
                        // Provider.of<SearchFormModel>(context, listen: false)
                        //     .visibilityRouteOptions = true;
                        // searchFormModel.busRoutes = futureBusRoutes;

                        // Use a new route to show the route options
                        // https://docs.flutter.dev/cookbook/navigation/navigation-basics
                        if (widget.screenSize == ScreenSize.mobile) {
                          Navigator.push(
                              context,
                              MaterialPageRoute(
                                  builder: (BuildContext context) => const RouteOptionsMobile()));
                        }
                      }
                    },
                    child: const Text('Plan'),
                  ),
                ),
              ],
            )),
      ),
    );
  }

  Widget buildSearchableOriginDropdownList(SearchFormModel searchFormModel) {
    // DropdownSearch widget plugin: https://pub.dev/packages/dropdown_search
    // Check the examples code for usage: https://github.com/salim-lachdhaf/searchable_dropdown
    return DropdownSearch<BusStop>(
      key: searchFormModel.originSelectionKey,
      asyncItems: (filter) => fetchFutureBusStopsByName(filter == '' ? 'donnybrook' : filter),
      compareFn: (i, s) => i.isEqual(s),
      dropdownDecoratorProps: const DropDownDecoratorProps(
        dropdownSearchDecoration: InputDecoration(
          labelText: 'Origin',
          border: OutlineInputBorder(),
          icon: Icon(Icons.map),
        ),
      ),
      popupProps: PopupProps.menu(
        showSearchBox: true,
        title: const Text('Search origin bus stop'),
        isFilterOnline: true,
        showSelectedItems: true,
        itemBuilder: _dropdownPopupItemBuilder,
        // favoriteItemProps: FavoriteItemProps(
        //   showFavoriteItems: true,
        //   // TODO This is a fake favorite feature. We need to implement in future to let the user
        //   //  mark the favorite stops
        //   favoriteItems: (us) {
        //     return us
        //       .where((e) => e.stopName!.contains("UCD"))
        //       .toList();
        //   },
        // ),
      ),
    );
  }

  Widget buildSearchableDestinationDropdownList(SearchFormModel searchFormModel) {
    return DropdownSearch<BusStop>(
      key: searchFormModel.destinationSelectionKey,
      asyncItems: (filter) => fetchFutureBusStopsByName(filter == '' ? 'ucd' : filter),
      compareFn: (i, s) => i.isEqual(s),
      dropdownDecoratorProps: const DropDownDecoratorProps(
        dropdownSearchDecoration: InputDecoration(
          labelText: 'Destination',
          border: OutlineInputBorder(),
          icon: Icon(Icons.map),
        ),
      ),
      popupProps: PopupProps.menu(
        showSearchBox: true,
        title: const Text('Search destination bus stop'),
        isFilterOnline: true,
        showSelectedItems: true,
        itemBuilder: _dropdownPopupItemBuilder,
        // favoriteItemProps: FavoriteItemProps(
        //   showFavoriteItems: true,
        //   favoriteItems: (us) {
        //     return us
        //         .where((e) => e.stopName!.contains("Spire"))
        //         .toList();
        //   },
        // ),
      ),
    );
  }

  Widget _dropdownPopupItemBuilder(BuildContext context, BusStop item, bool isSelected) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 8),
      decoration: !isSelected
          ? null
          : BoxDecoration(
              border: Border.all(color: Theme.of(context).primaryColor),
              borderRadius: BorderRadius.circular(5),
              color: Colors.white,
            ),
      child: ListTile(
        selected: isSelected,
        title: Text(item.stopName),
        subtitle: Text(item.stopNumber.toString()),
        leading: CircleAvatar(
          child: buildBusStopAvatarByType(item),
        ),
      ),
    );
  }

  // use icon to distinguish the bus stop type
  buildBusStopAvatarByType(BusStop item) {
    if (item.type == BusStopType.matched) {
      return const Icon(Icons.search);
    } else {
      return const Icon(Icons.location_searching);
    }
  }
}
