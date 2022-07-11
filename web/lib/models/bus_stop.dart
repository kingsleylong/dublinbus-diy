class BusStop {
  final String stopNumber;
  final String stopName;
  final double latitude;
  final double longitude;

  const BusStop({
    required this.stopNumber,
    required this.stopName,
    required this.latitude,
    required this.longitude,
  });

  factory BusStop.fromJson(Map<String, dynamic> json) {
    return BusStop(
      stopNumber: json['stop_number'],
      stopName: json['stop_name'],
      // latitude: json['stop_lat'],
      // longitude: json['stop_lon'],
      // TODO we may not need this
      latitude: 0.0,
      longitude: 0.0,
    );
  }

  factory BusStop.fromJsonForRoute(Map<String, dynamic> json) {
    return BusStop(
      stopNumber: json['stop_num'],
      stopName: json['stop_address'],
      // parse String to double: https://stackoverflow.com/a/13167498
      latitude: double.parse(json['stop_lat']),
      longitude: double.parse(json['stop_lon']),
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