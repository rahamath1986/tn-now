import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/main.dart';

void main() {
  testWidgets('TNNowApp splash screen smoke test', (WidgetTester tester) async {
    // Build our app under ProviderScope and trigger a frame.
    await tester.pumpWidget(
      const ProviderScope(
        child: TNNowApp(),
      ),
    );

    // Verify that the splash screen shows 'TN NOW' and the tagline
    expect(find.text('TN NOW'), findsOneWidget);
    expect(find.text('Everything Tamil Nadu, Right Now.'), findsOneWidget);

    // Fast-forward time for the 2-second routing timer and settle
    await tester.pump(const Duration(seconds: 2));
    await tester.pumpAndSettle();
  });
}
