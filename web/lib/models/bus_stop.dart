enum BusStopType {
  matched, nearby
}

class BusStop {
  String? stopNumber;
  String? stopName;
  double? latitude;
  double? longitude;
  Enum? type;

  BusStop(this.stopNumber, this.stopName, this.latitude, this.longitude, this.type);

  factory BusStop.fromJson(Map<String, dynamic> json, Enum? type) {
    return BusStop(
      json['stop_number'],
      json['stop_name'],
      json['stop_lat'],
      json['stop_lon'],
      type,
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