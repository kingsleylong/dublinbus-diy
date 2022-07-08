class BusStop {
  final String stopNumber;
  final String stopName;

  const BusStop({
    required this.stopNumber,
    required this.stopName,
  });

  factory BusStop.fromJson(Map<String, dynamic> json) {
    return BusStop(
      stopNumber: json['stop_number'],
      stopName: json['stop_name'],
    );
  }

  factory BusStop.fromJsonForRoute(Map<String, dynamic> json) {
    return BusStop(
      stopNumber: json['stop_num'],
      stopName: json['stop_address'],
    );
  }
}