import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/features/home/presentation/home_screen.dart';
import 'package:mobile/features/home/presentation/home_provider.dart';
import 'package:mobile/features/home/domain/home_feed.model.dart';
import 'package:mobile/core/localization/app_localizations.dart';
import 'package:mobile/core/network/dio_client.dart';
import 'package:mobile/core/storage/secure_storage.dart';
import 'package:dio/dio.dart';
import 'secure_storage_test.dart'; // import FakeFlutterSecureStorage
import 'dio_client_test.dart'; // import FakeHttpClientAdapter

class FakeHomeFeedNotifier extends HomeFeedNotifier {
  FakeHomeFeedNotifier(AsyncValue<HomeFeedModel> initialVal) 
      : super(DioClient(Dio()..httpClientAdapter = FakeHttpClientAdapter(), SecureStorage(FakeFlutterSecureStorage()))) {
    state = initialVal;
  }

  @override
  Future<void> fetchHomeFeed() async {}
}

class FakeLocaleNotifier extends LocaleNotifier {
  FakeLocaleNotifier() {
    state = const Locale('en');
  }
}

void main() {
  testWidgets('HomeScreen renders categories, trending list, and events list', (WidgetTester tester) async {
    final mockFeed = HomeFeedModel(
      categories: const [
        CategoryModel(id: 'cat-1', name: 'Viral Videos', slug: 'viral'),
        CategoryModel(id: 'cat-2', name: 'Local News', slug: 'news'),
      ],
      trending: [
        ContentCardModel(
          id: 'trend-1',
          title: 'Madurai Temple Fest Viral Video',
          description: 'Viral video showing festival celebration',
          contentType: 'VIDEO_LINK',
          categoryId: 'cat-1',
          districtId: 'dist-1',
          publishedAt: DateTime.now(),
          createdAt: DateTime.now(),
        ),
      ],
      latest: [
        ContentCardModel(
          id: 'post-1',
          title: 'Heavy Rains in Madurai Outskirts',
          description: 'Suburbs receive light showers today',
          contentType: 'PHOTO',
          categoryId: 'cat-2',
          districtId: 'dist-1',
          publishedAt: DateTime.now(),
          createdAt: DateTime.now(),
        ),
      ],
      events: [
        EventCardModel(
          id: 'event-1',
          title: 'College Fest Announcement',
          description: 'Annual cultural meet',
          eventName: 'Madurai Cultural Fest 2026',
          eventDate: DateTime.now().add(const Duration(days: 1)),
          startTime: '10:00 AM',
          endTime: '5:00 PM',
          organizerName: 'Thiagarajar College',
        ),
      ],
    );

    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          homeFeedNotifierProvider.overrideWith(
            (ref) => FakeHomeFeedNotifier(AsyncValue.data(mockFeed)),
          ),
          localeProvider.overrideWith((ref) => FakeLocaleNotifier()),
        ],
        child: const MaterialApp(
          home: HomeScreen(),
        ),
      ),
    );

    // Verify app title and location chip
    expect(find.text('TN NOW'), findsOneWidget);
    expect(find.text('Madurai'), findsOneWidget);

    // Verify categories chips
    expect(find.text('Viral Videos'), findsOneWidget);
    expect(find.text('Local News'), findsOneWidget);

    // Verify section headers
    expect(find.text('🔥 Trending Now'), findsOneWidget);
    expect(find.text('📅 Today\'s Events'), findsOneWidget);
    expect(find.text('What\'s Happening'), findsOneWidget);

    // Verify actual list elements
    expect(find.text('Madurai Temple Fest Viral Video'), findsOneWidget);
    expect(find.text('Madurai Cultural Fest 2026'), findsOneWidget);
    expect(find.text('Heavy Rains in Madurai Outskirts'), findsOneWidget);
  });
}
