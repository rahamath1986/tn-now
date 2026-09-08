import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/features/submission/presentation/submission_screen.dart';

void main() {
  testWidgets('SubmissionScreen renders title, format tabs, dropdowns, and submit button', (WidgetTester tester) async {
    await tester.pumpWidget(
      const ProviderScope(
        child: MaterialApp(
          home: SubmissionScreen(),
        ),
      ),
    );

    // Verify AppBar title
    expect(find.text('Submit Update'), findsOneWidget);

    // Verify Format tabs
    expect(find.text('Video'), findsOneWidget);
    expect(find.text('Photo'), findsOneWidget);
    expect(find.text('Story'), findsOneWidget);
    expect(find.text('Event'), findsOneWidget);

    // Scroll down outer ListView to reveal submit button
    await tester.scrollUntilVisible(
      find.text('Submit for Moderation'),
      200,
      scrollable: find.byType(Scrollable).first,
    );
    expect(find.text('Submit for Moderation'), findsOneWidget);

    // Tap Photo tab
    await tester.tap(find.text('Photo'));
    await tester.pumpAndSettle();

    // Verify Photo specific input field
    expect(find.text('Photo Image URL *'), findsOneWidget);
  });
}
