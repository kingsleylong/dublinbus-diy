class BusStop {
  final int busStopId;
  final String name;

  const BusStop({
    required this.busStopId,
    required this.name,
  });

  factory BusStop.fromJson(Map<String, dynamic> json) {
    return BusStop(
      busStopId: json['busStopId'],
      name: json['name'],
    );
  }
}