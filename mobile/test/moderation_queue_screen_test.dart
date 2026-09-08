import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/features/moderation/presentation/moderation_queue_screen.dart';

void main() {
  testWidgets('ModerationQueueScreen renders title, queue items, and action buttons', (WidgetTester tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: ModerationQueueScreen(),
      ),
    );

    // Verify AppBar title
    expect(find.text('Moderation Queue'), findsOneWidget);

    // Verify initial item titles
    expect(find.text('Unverified Madurai Incident Report'), findsOneWidget);
    expect(find.text('Viral Video Submission with Heavy Audio'), findsOneWidget);

    // Verify Approve and Reject buttons exist
    expect(find.text('Approve'), findsNWidgets(2));
    expect(find.text('Reject'), findsNWidgets(2));

    // Tap Approve on first item
    await tester.tap(find.text('Approve').first);
    await tester.pumpAndSettle();

    // Verify first item is removed from queue
    expect(find.text('Unverified Madurai Incident Report'), findsNothing);
  });
}
