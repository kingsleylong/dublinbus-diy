class Shape {
  double latitude;
  double longitude;

  Shape(this.latitude, this.longitude);

  factory Shape.fromJson(Map<String, dynamic> json) {
    return Shape(
      json['shape_pt_lat'],
      json['shape_pt_lon'],
    );
  }
}