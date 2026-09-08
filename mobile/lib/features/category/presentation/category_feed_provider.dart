import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/network/dio_client.dart';
import '../../home/domain/home_feed.model.dart';

class CategoryFeedState {
  final List<ContentCardModel> items;
  final bool isLoading;
  final bool isFetchingMore;
  final bool hasMore;
  final String? nextCursorTime;
  final String? nextCursorId;
  final String? errorMessage;

  const CategoryFeedState({
    this.items = const [],
    this.isLoading = true,
    this.isFetchingMore = false,
    this.hasMore = true,
    this.nextCursorTime,
    this.nextCursorId,
    this.errorMessage,
  });

  CategoryFeedState copyWith({
    List<ContentCardModel>? items,
    bool? isLoading,
    bool? isFetchingMore,
    bool? hasMore,
    String? nextCursorTime,
    String? nextCursorId,
    String? errorMessage,
  }) {
    return CategoryFeedState(
      items: items ?? this.items,
      isLoading: isLoading ?? this.isLoading,
      isFetchingMore: isFetchingMore ?? this.isFetchingMore,
      hasMore: hasMore ?? this.hasMore,
      nextCursorTime: nextCursorTime ?? this.nextCursorTime,
      nextCursorId: nextCursorId ?? this.nextCursorId,
      errorMessage: errorMessage,
    );
  }
}

class CategoryFeedNotifier extends StateNotifier<CategoryFeedState> {
  final DioClient _dioClient;
  final String categoryId;

  CategoryFeedNotifier(this._dioClient, this.categoryId)
      : super(const CategoryFeedState()) {
    fetchInitial();
  }

  Future<void> fetchInitial() async {
    state = state.copyWith(isLoading: true, errorMessage: null);
    try {
      final response = await _dioClient.get(
        '/content',
        queryParameters: {
          'category_id': categoryId,
          'limit': 10,
        },
      );

      final data = response.data;
      if (data != null && data['success'] == true) {
        final resData = data['data'] as Map<String, dynamic>;
        final itemsRaw = (resData['items'] as List? ?? [])
            .map((e) => ContentCardModel.fromJson(e as Map<String, dynamic>))
            .toList();

        state = state.copyWith(
          items: itemsRaw,
          isLoading: false,
          hasMore: resData['hasMore'] as bool? ?? false,
          nextCursorTime: resData['nextCursorTime'] as String?,
          nextCursorId: resData['nextCursorId'] as String?,
        );
      } else {
        state = state.copyWith(
          isLoading: false,
          errorMessage: data['message'] ?? 'Failed to fetch category posts',
        );
      }
    } catch (err) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: err.toString(),
      );
    }
  }

  Future<void> fetchNextPage() async {
    if (!state.hasMore || state.isFetchingMore || state.isLoading) return;

    state = state.copyWith(isFetchingMore: true);
    try {
      final queryParams = <String, dynamic>{
        'category_id': categoryId,
        'limit': 10,
      };
      if (state.nextCursorTime != null) {
        queryParams['cursor_time'] = state.nextCursorTime;
      }
      if (state.nextCursorId != null) {
        queryParams['cursor_id'] = state.nextCursorId;
      }

      final response = await _dioClient.get(
        '/content',
        queryParameters: queryParams,
      );

      final data = response.data;
      if (data != null && data['success'] == true) {
        final resData = data['data'] as Map<String, dynamic>;
        final newItems = (resData['items'] as List? ?? [])
            .map((e) => ContentCardModel.fromJson(e as Map<String, dynamic>))
            .toList();

        state = state.copyWith(
          items: [...state.items, ...newItems],
          isFetchingMore: false,
          hasMore: resData['hasMore'] as bool? ?? false,
          nextCursorTime: resData['nextCursorTime'] as String?,
          nextCursorId: resData['nextCursorId'] as String?,
        );
      } else {
        state = state.copyWith(isFetchingMore: false);
      }
    } catch (err) {
      state = state.copyWith(isFetchingMore: false);
    }
  }
}

final categoryFeedProvider = StateNotifierProvider.family<CategoryFeedNotifier,
    CategoryFeedState, String>((ref, categoryId) {
  final dioClient = ref.read(dioClientProvider);
  return CategoryFeedNotifier(dioClient, categoryId);
});
