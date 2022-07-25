class BusRouteSearchFilter {
  String originStopNumber;
  String destinationStopNumber;
  TimeType timeType;
  String time;

  BusRouteSearchFilter(
      this.originStopNumber, this.destinationStopNumber, this.timeType, this.time);
}

enum TimeType {
  departure,
  arrival,
}