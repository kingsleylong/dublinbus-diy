import 'package:date_time_picker/date_time_picker.dart';
import 'package:dropdown_search/dropdown_search.dart';
import 'package:flutter/material.dart';
import 'package:geolocator/geolocator.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:toggle_switch/toggle_switch.dart';

import '../../api/location_service.dart';
import '../../models/bus_route.dart';
import '../../models/bus_route_filter.dart';
import '../../models/map_polylines.dart';
import '../../models/responsive.dart';
import '../../models/search_form.dart';
import 'route_options_mobile.dart';

class SearchForm extends StatefulWidget {
  const SearchForm({Key? key, required this.screenSize}) : super(key: key);

  final ScreenType screenSize;

  @override
  State<SearchForm> createState() => _SearchFormState();
}

class _SearchFormState extends State<SearchForm> {
  late Future<List<BusRoute>> futureBusRoutes;

  var defaultDropdown = [Prediction('here', 'Use my location', '')];

  @override
  void initState() {
    super.initState();
    checkLocation();
  }

  checkLocation() async {
    Position here = await determinePosition();
    print(here);
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      // padding settings https://api.flutter.dev/flutter/material/InputDecoration/contentPadding.html
      padding: const EdgeInsets.fromLTRB(0, 10, 0, 0),
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
                  padding: const EdgeInsets.symmetric(horizontal: 2, vertical: 8),
                  child: Row(
                    children: [
                      Expanded(child: buildSearchableOriginDropdownList(searchFormModel)),
                    ],
                  ),
                ),
                // Destination dropdown list
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 2, vertical: 8),
                  child: buildSearchableDestinationDropdownList(searchFormModel),
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 2, vertical: 8),
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
                  padding: const EdgeInsets.symmetric(horizontal: 2, vertical: 8),
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
                            'Selected origin: ${searchFormModel.originSelectionKey.currentState?.getSelectedItem}');
                        print(
                            'coordinates: ${searchFormModel.originPlaceDetail} ${searchFormModel.destinationPlaceDetail}');
                        print(
                            'Selected destination: ${searchFormModel.destinationSelectionKey.currentState?.getSelectedItem}');
                        print(
                            'Selected datetime: ${searchFormModel.dateTimePickerController.value.text}');

                        DateTime parseTime =
                            DateTime.parse(searchFormModel.dateTimePickerController.value.text);

                        // Date format: https://api.flutter.dev/flutter/intl/DateFormat-class.html
                        print(
                            'parseTime: $parseTime  ${DateFormat('yyyy-MM-dd HH:mm:ss').format(parseTime)}');

                        BusRouteSearchFilter searchFilter = BusRouteSearchFilter(
                            searchFormModel.originPlaceDetail.toString(),
                            searchFormModel.destinationPlaceDetail.toString(),
                            searchFormModel.timeTypes[searchFormModel.timeTypeToggleIndex ?? 0],
                            DateFormat('yyyy-MM-dd HH:mm:ss').format(parseTime));

                        // futureBusRoutes = fetchBusRoutes(searchFilter);
                        Provider.of<SearchFormModel>(context, listen: false)
                            .fetchBusRoute(searchFilter);

                        // Use a new route to show the route options
                        // https://docs.flutter.dev/cookbook/navigation/navigation-basics
                        if (widget.screenSize == ScreenType.mobile) {
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
    // TODO replace secondary text by location
    return DropdownSearch<Prediction>(
      key: searchFormModel.originSelectionKey,
      asyncItems: (filter) => searchFormModel.autocompleteAddress(filter),
      compareFn: (i, s) => i.placeId == s.placeId,
      onChanged: (value) {
        searchFormModel.fetchPlaceDetails(value, 'origin');
      },
      // TODO may add a history list
      items: defaultDropdown,
      dropdownDecoratorProps: const DropDownDecoratorProps(
        dropdownSearchDecoration: InputDecoration(
          labelText: 'Origin',
          border: OutlineInputBorder(),
          icon: Icon(Icons.map),
        ),
      ),
      popupProps: PopupProps.menu(
        showSearchBox: true,
        title: const Text('Search origin address'),
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
      validator: (value) {
        if (value == null) {
          return "Please search the origin address";
        }
        return null;
      },
    );
  }

  Widget buildSearchableDestinationDropdownList(SearchFormModel searchFormModel) {
    return DropdownSearch<Prediction>(
      key: searchFormModel.destinationSelectionKey,
      asyncItems: (filter) => searchFormModel.autocompleteAddress(filter),
      compareFn: (i, s) => i.placeId == s.placeId,
      items: defaultDropdown,
      onChanged: (value) {
        searchFormModel.fetchPlaceDetails(value, 'destination');
      },
      dropdownDecoratorProps: const DropDownDecoratorProps(
        dropdownSearchDecoration: InputDecoration(
          labelText: 'Destination',
          border: OutlineInputBorder(),
          icon: Icon(Icons.map),
        ),
      ),
      popupProps: PopupProps.menu(
        showSearchBox: true,
        title: const Text('Search destination address'),
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
      validator: (value) {
        if (value == null) {
          return "Please search the destination address";
        }
        return null;
      },
    );
  }

  Widget _dropdownPopupItemBuilder(BuildContext context, Prediction prediction, bool isSelected) {
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
        // title: Text(item.terms.map((e) => e.value).reduce((value, element) => '$value, $element')),
        title: Text(prediction.mainText),
        subtitle: Text(prediction.secondaryText),
        leading: CircleAvatar(
          child: buildBusStopAvatarByType(prediction),
        ),
      ),
    );
  }

  // use icon to distinguish the bus stop type
  buildBusStopAvatarByType(Prediction prediction) {
    if (prediction.placeId == 'here') {
      return const Icon(Icons.place);
    } else {
      return const Icon(Icons.location_searching);
    }
  }
}
