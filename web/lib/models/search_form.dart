import 'package:dropdown_search/dropdown_search.dart';
import 'package:flutter/material.dart';

import 'bus_stop.dart';

/// The model for the search form
class SearchFormModel extends ChangeNotifier {
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();

  // Use a flag to control the visibility of the route options
  // https://stackoverflow.com/a/46126667
  bool visibilityRouteOptions = false;

  // The instance field that holds the state of origin dropdown list
  final _originSelectionKey = GlobalKey<DropdownSearchState<BusStop>>();

  // The instance field that holds the state of destination dropdown list
  final _destinationSelectionKey = GlobalKey<DropdownSearchState<BusStop>>();

  // The instance field that holds the state of the datetime picker
  final TextEditingController _dateTimePickerController =
      TextEditingController(text: DateTime.now().toString());

  // getters
  TextEditingController get dateTimePickerController =>
      _dateTimePickerController;

  get destinationSelectionKey => _destinationSelectionKey;

  get originSelectionKey => _originSelectionKey;

  GlobalKey<FormState> get formKey => _formKey;
}
