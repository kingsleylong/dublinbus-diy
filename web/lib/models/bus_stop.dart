class BusStop {
  final String stopId;
  final String stopName;

  const BusStop({
    required this.stopId,
    required this.stopName,
  });

  factory BusStop.fromJson(Map<String, dynamic> json) {
    return BusStop(
      stopId: json['stop_id'],
      stopName: json['stop_name'],
    );
  }
}