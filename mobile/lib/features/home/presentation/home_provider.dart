import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/network/dio_client.dart';
import '../domain/home_feed.model.dart';

class HomeFeedNotifier extends StateNotifier<AsyncValue<HomeFeedModel>> {
  final DioClient _dioClient;

  HomeFeedNotifier(this._dioClient) : super(const AsyncValue.loading()) {
    fetchHomeFeed();
  }

  Future<void> fetchHomeFeed() async {
    state = const AsyncValue.loading();
    try {
      final response = await _dioClient.get('/home');
      
      // The backend returns standard: { success: true, data: { ... } }
      final data = response.data;
      if (data != null && data['success'] == true) {
        final feedData = data['data'] as Map<String, dynamic>;
        final feed = HomeFeedModel.fromJson(feedData);
        state = AsyncValue.data(feed);
      } else {
        state = AsyncValue.error(
          data['message'] ?? 'Failed to load home feed data',
          StackTrace.current,
        );
      }
    } catch (err, stack) {
      state = AsyncValue.error(err, stack);
    }
  }
}

// Riverpod provider for HomeFeedNotifier
final homeFeedNotifierProvider =
    StateNotifierProvider<HomeFeedNotifier, AsyncValue<HomeFeedModel>>((ref) {
  final dioClient = ref.read(dioClientProvider);
  return HomeFeedNotifier(dioClient);
});
