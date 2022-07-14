class BusStop {
  String? stopNumber;
  String? stopName;
  double? latitude;
  double? longitude;

  BusStop(this.stopNumber, this.stopName, this.latitude, this.longitude);

  factory BusStop.fromJson(Map<String, dynamic> json) {
    return BusStop(
      json['stop_number'],
      json['stop_name'],
      json['stop_lat'],
      json['stop_lon'],
    );
  }

  bool isEqual(BusStop s) {
    return stopNumber == s.stopNumber;
  }

  // This method is required by DropdownSearch widget to display the BusStop object.
  @override
  String toString() {
    return '$stopName - $stopNumber';
  }
}