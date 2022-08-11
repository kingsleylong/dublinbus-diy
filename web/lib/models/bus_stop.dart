enum BusStopType { matched, nearby }

class BusStop {
  String stopNumber;
  String stopName;
  double? latitude;
  double? longitude;
  Enum? type;
  String arrivalTime;

  BusStop(
      this.stopNumber, this.stopName, this.latitude, this.longitude, this.type, this.arrivalTime);

  factory BusStop.fromJson(Map<String, dynamic> json, Enum? type) {
    return BusStop(
      json['stop_number'],
      json['stop_name'],
      json['stop_lat'],
      json['stop_lon'],
      type,
      json['arrival_time'],
    );
  }

  bool isEqual(BusStop s) {
    return stopNumber == s.stopNumber;
  }

  // This method is required by DropdownSearch widget to display the BusStop object.
  @override
  String toString() {
    return '$stopName, stop $stopNumber';
  }
}
