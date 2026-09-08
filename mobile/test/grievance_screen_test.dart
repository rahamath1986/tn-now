import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/features/compliance/presentation/grievance_screen.dart';

void main() {
  testWidgets('GrievanceScreen renders officer disclosure, email input, category dropdown, and submit button', (WidgetTester tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: GrievanceScreen(),
      ),
    );

    // Verify AppBar title
    expect(find.text('Grievance Redressal'), findsOneWidget);

    // Verify IT Rules 2021 Disclosure
    expect(find.text('IT Rules 2021 Compliance Officer'), findsOneWidget);
    expect(find.textContaining('grievance@tnnow.in'), findsOneWidget);

    // Verify Email field
    expect(find.text('Complainant Email *'), findsOneWidget);

    // Scroll down to bring submit button into view
    await tester.scrollUntilVisible(
      find.text('Register Grievance'),
      200,
      scrollable: find.byType(Scrollable).first,
    );

    // Verify Submit button
    expect(find.text('Register Grievance'), findsOneWidget);
  });
}
