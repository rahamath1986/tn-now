import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/features/category/presentation/category_detail_screen.dart';
import 'package:mobile/features/category/presentation/category_feed_provider.dart';
import 'package:mobile/features/home/domain/home_feed.model.dart';
import 'package:mobile/core/localization/app_localizations.dart';
import 'package:mobile/core/network/dio_client.dart';
import 'package:mobile/core/storage/secure_storage.dart';
import 'package:dio/dio.dart';
import 'secure_storage_test.dart';
import 'dio_client_test.dart';

class FakeCategoryFeedNotifier extends CategoryFeedNotifier {
  FakeCategoryFeedNotifier(CategoryFeedState initialVal, String categoryId)
      : super(
          DioClient(
            Dio()..httpClientAdapter = FakeHttpClientAdapter(),
            SecureStorage(FakeFlutterSecureStorage()),
          ),
          categoryId,
        ) {
    state = initialVal;
  }

  @override
  Future<void> fetchInitial() async {}

  @override
  Future<void> fetchNextPage() async {}
}

class FakeLocaleNotifier extends LocaleNotifier {
  FakeLocaleNotifier() {
    state = const Locale('en');
  }
}

void main() {
  testWidgets('CategoryDetailScreen renders title and list of articles', (WidgetTester tester) async {
    const categoryId = 'cat-123';
    final mockState = CategoryFeedState(
      isLoading: false,
      hasMore: false,
      items: [
        ContentCardModel(
          id: 'post-1',
          title: 'Category Exclusive Article 1',
          description: 'Description for article 1',
          contentType: 'PHOTO',
          categoryId: categoryId,
          districtId: 'dist-1',
          publishedAt: DateTime.now(),
          createdAt: DateTime.now(),
        ),
        ContentCardModel(
          id: 'post-2',
          title: 'Category Exclusive Article 2',
          description: 'Description for article 2',
          contentType: 'VIDEO_LINK',
          categoryId: categoryId,
          districtId: 'dist-1',
          publishedAt: DateTime.now(),
          createdAt: DateTime.now(),
        ),
      ],
    );

    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          categoryFeedProvider(categoryId).overrideWith(
            (ref) => FakeCategoryFeedNotifier(mockState, categoryId),
          ),
          localeProvider.overrideWith((ref) => FakeLocaleNotifier()),
        ],
        child: const MaterialApp(
          home: CategoryDetailScreen(
            categoryId: categoryId,
            categoryName: 'Viral Videos',
          ),
        ),
      ),
    );

    // Verify AppBar Title
    expect(find.text('Viral Videos'), findsOneWidget);

    // Verify Article Titles
    expect(find.text('Category Exclusive Article 1'), findsOneWidget);
    expect(find.text('Category Exclusive Article 2'), findsOneWidget);
  });
}
