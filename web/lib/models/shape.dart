class Shape {
  double latitude;
  double longitude;
  int sequence;
  String distanceTraveled;

  Shape(this.latitude, this.longitude, this.sequence, this.distanceTraveled);

  factory Shape.fromJson(Map<String, dynamic> json) {
    return Shape(
      json['shape_pt_lat'],
      json['shape_pt_lon'],
      json['shape_pt_sequence'],
      json['shape_dist_traveled'],
    );
  }
}