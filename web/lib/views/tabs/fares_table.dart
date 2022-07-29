import 'package:dublin_bus_diy/models/bus_route.dart';
import 'package:flutter/material.dart';

class FaresTable extends StatelessWidget {
  const FaresTable({Key? key, required this.fares}) : super(key: key);

  final Fares fares;

  @override
  Widget build(BuildContext context) {
    // create a table: https://api.flutter.dev/flutter/material/DataTable-class.html
    return DataTable(
      columns: const <DataColumn>[
        DataColumn(
          label: Text(
            'Fairs',
            style: TextStyle(fontWeight: FontWeight.bold),
          ),
        ),
        DataColumn(
          label: Text(
            'Adult',
            style: TextStyle(fontStyle: FontStyle.italic),
          ),
        ),
        DataColumn(
          label: Text(
            'Student',
            style: TextStyle(fontStyle: FontStyle.italic),
          ),
        ),
        DataColumn(
          label: Text(
            'Child',
            style: TextStyle(fontStyle: FontStyle.italic),
          ),
        ),
      ],
      rows: <DataRow>[
        DataRow(
          cells: <DataCell>[
            const DataCell(Text('Cash')),
            DataCell(Text('${fares.adultCash}')),
            const DataCell(Text('')),
            DataCell(Text('${fares.childLeap}')),
          ],
        ),
        DataRow(
          cells: <DataCell>[
            const DataCell(Text('Leap Card')),
            DataCell(Text('${fares.adultLeap}')),
            DataCell(Text('${fares.studentLeap}')),
            DataCell(Text('${fares.childLeap}')),
          ],
        ),
      ],
    );
  }
}
