import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/widgets/design_system.dart';
import 'category_feed_provider.dart';

class CategoryDetailScreen extends ConsumerStatefulWidget {
  final String categoryId;
  final String? categoryName;

  const CategoryDetailScreen({
    super.key,
    required this.categoryId,
    this.categoryName,
  });

  @override
  ConsumerState<CategoryDetailScreen> createState() => _CategoryDetailScreenState();
}

class _CategoryDetailScreenState extends ConsumerState<CategoryDetailScreen> {
  late final ScrollController _scrollController;

  @override
  void initState() {
    super.initState();
    _scrollController = ScrollController()..addListener(_onScroll);
  }

  void _onScroll() {
    if (_scrollController.position.pixels >=
        _scrollController.position.maxScrollExtent - 200) {
      ref.read(categoryFeedProvider(widget.categoryId).notifier).fetchNextPage();
    }
  }

  @override
  void dispose() {
    _scrollController.removeListener(_onScroll);
    _scrollController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(categoryFeedProvider(widget.categoryId));

    return Scaffold(
      appBar: AppBar(
        title: Text(
          widget.categoryName ?? 'Category Updates',
          style: const TextStyle(fontWeight: FontWeight.bold, color: Colors.white),
        ),
      ),
      body: RefreshIndicator(
        onRefresh: () =>
            ref.read(categoryFeedProvider(widget.categoryId).notifier).fetchInitial(),
        color: AppTheme.primaryOrange,
        backgroundColor: AppTheme.darkSurface,
        child: _buildBody(state),
      ),
    );
  }

  Widget _buildBody(CategoryFeedState state) {
    if (state.isLoading) {
      return ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: 5,
        itemBuilder: (context, index) => const Padding(
          padding: EdgeInsets.only(bottom: 16),
          child: TNSkeleton(width: double.infinity, height: 120, borderRadius: 12),
        ),
      );
    }

    if (state.errorMessage != null && state.items.isEmpty) {
      return TNErrorState(
        message: state.errorMessage!,
        onRetry: () =>
            ref.read(categoryFeedProvider(widget.categoryId).notifier).fetchInitial(),
      );
    }

    if (state.items.isEmpty) {
      return const TNEmptyState(
        title: 'No Updates In This Category',
        description: 'Be the first contributor to share news in this category!',
      );
    }

    return ListView.builder(
      controller: _scrollController,
      physics: const AlwaysScrollableScrollPhysics(),
      padding: const EdgeInsets.all(16),
      itemCount: state.items.length + (state.isFetchingMore ? 1 : 0),
      itemBuilder: (context, index) {
        if (index == state.items.length) {
          return const Padding(
            padding: EdgeInsets.symmetric(vertical: 16),
            child: Center(
              child: CircularProgressIndicator(
                valueColor: AlwaysStoppedAnimation<Color>(AppTheme.primaryOrange),
              ),
            ),
          );
        }

        final item = state.items[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 12),
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    const CircleAvatar(
                      radius: 14,
                      backgroundColor: AppTheme.darkOverlay,
                      child: Icon(Icons.person, size: 14, color: Colors.white),
                    ),
                    const SizedBox(width: 8),
                    const Text(
                      'TN Contributor',
                      style: TextStyle(
                          fontSize: 12,
                          fontWeight: FontWeight.bold,
                          color: Colors.white),
                    ),
                    const Spacer(),
                    Container(
                      padding:
                          const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                      decoration: BoxDecoration(
                        color: AppTheme.primaryOrange.withValues(alpha: 0.2),
                        borderRadius: BorderRadius.circular(4),
                      ),
                      child: Text(
                        item.contentType,
                        style: const TextStyle(
                            fontSize: 9,
                            color: AppTheme.primaryOrange,
                            fontWeight: FontWeight.bold),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                Text(
                  item.title,
                  style: const TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.bold,
                      color: Colors.white),
                ),
                if (item.description.isNotEmpty) ...[
                  const SizedBox(height: 6),
                  Text(
                    item.description,
                    maxLines: 3,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(fontSize: 13, color: Colors.grey),
                  ),
                ],
                const SizedBox(height: 12),
                Row(
                  children: [
                    IconButton(
                      icon: const Icon(Icons.thumb_up_alt_outlined,
                          size: 18, color: Colors.grey),
                      onPressed: () {},
                    ),
                    const Text('18',
                        style: TextStyle(fontSize: 12, color: Colors.grey)),
                    const SizedBox(width: 16),
                    IconButton(
                      icon: const Icon(Icons.mode_comment_outlined,
                          size: 18, color: Colors.grey),
                      onPressed: () {},
                    ),
                    const Text('5',
                        style: TextStyle(fontSize: 12, color: Colors.grey)),
                  ],
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}
