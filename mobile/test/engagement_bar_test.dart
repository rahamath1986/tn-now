import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/features/social/presentation/engagement_bar.dart';

void main() {
  testWidgets('EngagementBar renders like count and toggles optimistic state', (WidgetTester tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: Scaffold(
          body: EngagementBar(
            contentId: 'content-1',
            initialLikeCount: 10,
            initialCommentCount: 3,
          ),
        ),
      ),
    );

    // Verify initial like count text
    expect(find.text('10'), findsOneWidget);
    expect(find.text('3'), findsOneWidget);

    // Tap Like Button
    await tester.tap(find.byIcon(Icons.favorite_border));
    await tester.pumpAndSettle();

    // Verify optimistic increment to 11
    expect(find.text('11'), findsOneWidget);
    expect(find.byIcon(Icons.favorite), findsOneWidget);
  });
}
