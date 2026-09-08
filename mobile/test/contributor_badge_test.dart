import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/features/reputation/presentation/contributor_badge.dart';
import 'package:mobile/features/reputation/presentation/contributor_profile_screen.dart';

void main() {
  testWidgets('ContributorBadge renders level and founding badge', (WidgetTester tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: Scaffold(
          body: ContributorBadge(
            level: 'Trusted Contributor',
            isFounding: true,
          ),
        ),
      ),
    );

    // Verify Level badge text
    expect(find.text('Trusted Contributor'), findsOneWidget);

    // Verify Founding badge text
    expect(find.text('Founding 2026'), findsOneWidget);
  });

  testWidgets('ContributorProfileScreen renders user profile, trust score meter, and stats', (WidgetTester tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: ContributorProfileScreen(
          userId: 'user-77',
          username: 'Madurai Reporter',
          level: 'Verified Contributor',
          points: 250,
          trustScore: 85,
          approvedCount: 25,
          isFounding: true,
        ),
      ),
    );

    // Verify Username
    expect(find.text('Madurai Reporter'), findsOneWidget);

    // Verify Trust Score Meter Text
    expect(find.text('85 / 100'), findsOneWidget);

    // Verify Points Stat Card
    expect(find.text('250'), findsOneWidget);
    expect(find.text('Points'), findsOneWidget);
  });
}
