import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/features/search/presentation/search_screen.dart';

void main() {
  testWidgets('SearchScreen renders search bar, district filter chips, and results', (WidgetTester tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: SearchScreen(),
      ),
    );

    // Verify initial search prompt text
    expect(find.text('Type a keyword to discover Tamil Nadu updates 🔍'), findsOneWidget);

    // Verify District filter chips
    expect(find.text('All'), findsOneWidget);
    expect(find.text('Madurai'), findsOneWidget);
    expect(find.text('Chennai'), findsOneWidget);

    // Enter query in search bar
    await tester.enterText(find.byType(TextField), 'Madurai');
    await tester.pumpAndSettle();

    // Verify matching search result card
    expect(find.text('Madurai Meenakshi Temple Annual Chithirai Festival 2026'), findsOneWidget);
  });
}
