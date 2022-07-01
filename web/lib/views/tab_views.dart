import 'package:flutter/material.dart';

class GetMeThereOnTimeTabView extends StatefulWidget {
  const GetMeThereOnTimeTabView({Key? key}) : super(key: key);

  @override
  State<GetMeThereOnTimeTabView> createState() => _GetMeThereOnTimeTabViewState();
}

class _GetMeThereOnTimeTabViewState extends State<GetMeThereOnTimeTabView> {
  @override
  Widget build(BuildContext context) {
    return Container(
        child: const Text("Get me there on time")
    );
  }
}

class PlanMyJourneyTabView extends StatefulWidget {
  const PlanMyJourneyTabView({Key? key}) : super(key: key);

  @override
  State<PlanMyJourneyTabView> createState() => _PlanMyJourneyTabViewState();
}

class _PlanMyJourneyTabViewState extends State<PlanMyJourneyTabView> {
  String dropdownValue = '175';
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();

  @override
  Widget build(BuildContext context) {
    const searchFieldsDecoration = InputDecoration(
      // icon: Icon(Icons.),
      labelText: "Origin",
      floatingLabelAlignment: FloatingLabelAlignment.start,
      hintText: 'Origin',
      // helperText: 'Select the origin',
      // counterText: '0 characters',
      border: OutlineInputBorder(),
    );

    return Padding(
      // padding settings https://api.flutter.dev/flutter/material/InputDecoration/contentPadding.html
      padding: const EdgeInsets.fromLTRB(5, 10, 5, 0),
      child: Form(
        key: _formKey,
        child: Column(
          children: <Widget>[
            DropdownButtonFormField(
              // how to build a drop down list https://api.flutter.dev/flutter/material/DropdownButton-class.htm
              value: dropdownValue,
              // Field decoration https://api.flutter.dev/flutter/material/InputDecoration-class.html
              decoration: searchFieldsDecoration,
              items: <String>["175", "C1", "46A", "52"]
                  .map<DropdownMenuItem<String>>((String value) {
                return DropdownMenuItem<String>(
                  value: value,
                  child: Text(value),
                );
              }).toList(),
              onChanged: (String? value) {
                setState(() {
                  dropdownValue = value!;
                });
              },
              // isExpanded: true,
            ),
            DropdownButtonFormField(
              // how to build a drop down list https://api.flutter.dev/flutter/material/DropdownButton-class.htm
              value: dropdownValue,
              // Field decoration https://api.flutter.dev/flutter/material/InputDecoration-class.html
              decoration: searchFieldsDecoration,
              items: <String>["175", "C1", "46A", "52"]
                  .map<DropdownMenuItem<String>>((String value) {
                return DropdownMenuItem<String>(
                  value: value,
                  child: Text(value),
                );
              }).toList(),
              onChanged: (String? value) {
                setState(() {
                  dropdownValue = value!;
                });
              },
              // isExpanded: true,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 16.0),
              child: ElevatedButton(
                onPressed: () {
                  // Validate will return true if the form is valid, or false if
                  // the form is invalid.
                  if (_formKey.currentState!.validate()) {
                    // Process data.
                  }
                },
                child: const Text('Plan'),
              ),
            ),
          ],
        )
      )
    );
  }
}


// for building dynamic lists
// ListView.builder(
// // specify the length of data, without this an index out of range error will be
// // thrown. https://stackoverflow.com/a/58850610
// itemCount: _lines.length,
// itemBuilder: (context, index) {
// return ListTile(
// title: Text(
// _lines[index],
// // strutStyle: ,
// ),
// );
// }
// ),

