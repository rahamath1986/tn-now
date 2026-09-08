import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/widgets/design_system.dart';
import '../../../core/localization/app_localizations.dart';
import 'home_provider.dart';
import '../domain/home_feed.model.dart';

class HomeScreen extends ConsumerWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final localizations = ref.watch(localizationsProvider);
    final homeFeedState = ref.watch(homeFeedNotifierProvider);

    return Scaffold(
      appBar: AppBar(
        automaticallyImplyLeading: false,
        title: Row(
          children: [
            Text(
              localizations.translate('app_title'),
              style: const TextStyle(
                color: AppTheme.primaryOrange,
                fontWeight: FontWeight.bold,
                fontSize: 24,
                letterSpacing: 0.5,
              ),
            ),
            const SizedBox(width: 8),
            // Location Chip selector
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
              decoration: BoxDecoration(
                color: AppTheme.darkOverlay,
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: AppTheme.primaryOrange.withValues(alpha: 0.3)),
              ),
              child: const Row(
                children: [
                  Icon(Icons.location_on, size: 14, color: AppTheme.primaryOrange),
                  SizedBox(width: 4),
                  Text(
                    'Madurai',
                    style: TextStyle(fontSize: 12, fontWeight: FontWeight.bold, color: Colors.white),
                  ),
                ],
              ),
            ),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.search, color: Colors.white),
            onPressed: () => context.push('/search'),
          ),
          IconButton(
            icon: const Icon(Icons.notifications_none, color: Colors.white),
            onPressed: () => context.push('/notifications'),
          ),
          const SizedBox(width: 8),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () => ref.read(homeFeedNotifierProvider.notifier).fetchHomeFeed(),
        color: AppTheme.primaryOrange,
        backgroundColor: AppTheme.darkSurface,
        child: homeFeedState.when(
          data: (feed) {
            if (feed.categories.isEmpty && feed.latest.isEmpty && feed.events.isEmpty) {
              return const TNEmptyState(
                title: 'No Updates Found',
                description: 'Check back later for trending updates in Tamil Nadu.',
              );
            }
            return _buildFeedContent(context, ref, feed);
          },
          loading: () => _buildLoadingShimmer(),
          error: (error, _) => TNErrorState(
            message: error.toString(),
            onRetry: () => ref.read(homeFeedNotifierProvider.notifier).fetchHomeFeed(),
          ),
        ),
      ),
    );
  }

  Widget _buildFeedContent(BuildContext context, WidgetRef ref, HomeFeedModel feed) {
    return ListView(
      physics: const AlwaysScrollableScrollPhysics(),
      children: [
        // Categories list
        if (feed.categories.isNotEmpty) ...[
          const SizedBox(height: 12),
          SizedBox(
            height: 40,
            child: ListView.builder(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 16),
              itemCount: feed.categories.length,
              itemBuilder: (context, index) {
                final category = feed.categories[index];
                return Padding(
                  padding: const EdgeInsets.only(right: 8.0),
                  child: TNCategoryChip(
                    label: category.name,
                    isSelected: index == 0, // Mock highlighting first category
                    onSelected: (selected) {},
                  ),
                );
              },
            ),
          ),
        ],

        // Trending Carousel section
        if (feed.trending.isNotEmpty) ...[
          const Padding(
            padding: EdgeInsets.only(left: 16.0, top: 20.0, bottom: 10.0),
            child: Text(
              '🔥 Trending Now',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Colors.white),
            ),
          ),
          SizedBox(
            height: 160,
            child: ListView.builder(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 16),
              itemCount: feed.trending.length,
              itemBuilder: (context, index) {
                final item = feed.trending[index];
                return Container(
                  width: 280,
                  margin: const EdgeInsets.only(right: 12.0),
                  decoration: BoxDecoration(
                    color: AppTheme.darkSurface,
                    borderRadius: BorderRadius.circular(16),
                    gradient: LinearGradient(
                      begin: Alignment.topRight,
                      end: Alignment.bottomLeft,
                      colors: [
                        AppTheme.primaryOrange.withValues(alpha: 0.15),
                        AppTheme.darkSurface,
                      ],
                    ),
                  ),
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Container(
                              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                              decoration: BoxDecoration(
                                color: AppTheme.primaryOrange.withValues(alpha: 0.2),
                                borderRadius: BorderRadius.circular(8),
                              ),
                              child: Text(
                                item.contentType,
                                style: const TextStyle(fontSize: 10, color: AppTheme.primaryOrange, fontWeight: FontWeight.bold),
                              ),
                            ),
                            const Text(
                              'via YouTube', // mock attribution
                              style: TextStyle(fontSize: 11, color: Colors.grey),
                            ),
                          ],
                        ),
                        Text(
                          item.title,
                          style: const TextStyle(fontSize: 14, fontWeight: FontWeight.bold, color: Colors.white),
                          maxLines: 2,
                          overflow: TextOverflow.ellipsis,
                        ),
                        const Row(
                          children: [
                            Icon(Icons.thumb_up_alt_outlined, size: 14, color: Colors.grey),
                            SizedBox(width: 4),
                            Text('24 likes', style: TextStyle(fontSize: 12, color: Colors.grey)),
                          ],
                        ),
                      ],
                    ),
                  ),
                );
              },
            ),
          ),
        ],

        // Upcoming Events section
        if (feed.events.isNotEmpty) ...[
          const Padding(
            padding: EdgeInsets.only(left: 16.0, top: 24.0, bottom: 10.0),
            child: Text(
              '📅 Today\'s Events',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Colors.white),
            ),
          ),
          ListView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            padding: const EdgeInsets.symmetric(horizontal: 16),
            itemCount: feed.events.length,
            itemBuilder: (context, index) {
              final event = feed.events[index];
              return Card(
                margin: const EdgeInsets.only(bottom: 12.0),
                child: ListTile(
                  contentPadding: const EdgeInsets.all(12),
                  leading: Container(
                    width: 56,
                    height: 56,
                    decoration: BoxDecoration(
                      color: AppTheme.primaryOrange.withValues(alpha: 0.2),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: const Icon(Icons.event, color: AppTheme.primaryOrange, size: 28),
                  ),
                  title: Text(
                    event.eventName,
                    style: const TextStyle(fontWeight: FontWeight.bold, color: Colors.white),
                  ),
                  subtitle: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const SizedBox(height: 4),
                      Text(event.organizerName, style: const TextStyle(fontSize: 12, color: Colors.grey)),
                      const SizedBox(height: 2),
                      Row(
                        children: [
                          const Icon(Icons.schedule, size: 12, color: Colors.grey),
                          const SizedBox(width: 4),
                          Text(event.startTime, style: const TextStyle(fontSize: 12, color: Colors.grey)),
                        ],
                      ),
                    ],
                  ),
                  trailing: const Icon(Icons.chevron_right, color: Colors.grey),
                  onTap: () => context.push('/events/${event.id}'),
                ),
              );
            },
          ),
        ],

        // What's Happening (Latest posts) list
        if (feed.latest.isNotEmpty) ...[
          const Padding(
            padding: EdgeInsets.only(left: 16.0, top: 24.0, bottom: 10.0),
            child: Text(
              'What\'s Happening',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Colors.white),
            ),
          ),
          ListView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            padding: const EdgeInsets.symmetric(horizontal: 16),
            itemCount: feed.latest.length,
            itemBuilder: (context, index) {
              final item = feed.latest[index];
              return Card(
                margin: const EdgeInsets.only(bottom: 12.0),
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          const CircleAvatar(
                            radius: 16,
                            backgroundColor: AppTheme.darkOverlay,
                            child: Icon(Icons.person, size: 16, color: Colors.white),
                          ),
                          const SizedBox(width: 8),
                          const Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text('Local Contributor', style: TextStyle(fontSize: 12, fontWeight: FontWeight.bold, color: Colors.white)),
                              Text('10 minutes ago', style: TextStyle(fontSize: 10, color: Colors.grey)),
                            ],
                          ),
                          const Spacer(),
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                            decoration: BoxDecoration(
                              color: AppTheme.darkOverlay,
                              borderRadius: BorderRadius.circular(4),
                            ),
                            child: Text(
                              item.contentType,
                              style: const TextStyle(fontSize: 9, color: Colors.grey),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 12),
                      Text(
                        item.title,
                        style: const TextStyle(fontWeight: FontWeight.bold, color: Colors.white),
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
                            icon: const Icon(Icons.thumb_up_alt_outlined, size: 20, color: Colors.grey),
                            onPressed: () {},
                          ),
                          const Text('12', style: TextStyle(fontSize: 12, color: Colors.grey)),
                          const SizedBox(width: 16),
                          IconButton(
                            icon: const Icon(Icons.mode_comment_outlined, size: 20, color: Colors.grey),
                            onPressed: () {},
                          ),
                          const Text('3', style: TextStyle(fontSize: 12, color: Colors.grey)),
                        ],
                      ),
                    ],
                  ),
                ),
              );
            },
          ),
        ],
        const SizedBox(height: 24),
      ],
    );
  }

  Widget _buildLoadingShimmer() {
    return ListView(
      padding: const EdgeInsets.all(16.0),
      children: [
        const SizedBox(height: 8),
        Row(
          children: List.generate(4, (index) => const Padding(
            padding: EdgeInsets.only(right: 8.0),
            child: TNSkeleton(width: 80, height: 36, borderRadius: 18),
          )),
        ),
        const SizedBox(height: 32),
        const TNSkeleton(width: 150, height: 20),
        const SizedBox(height: 16),
        SizedBox(
          height: 150,
          child: ListView.builder(
            scrollDirection: Axis.horizontal,
            itemCount: 3,
            itemBuilder: (context, index) => const Padding(
              padding: EdgeInsets.only(right: 12.0),
              child: TNSkeleton(width: 250, height: 150, borderRadius: 16),
            ),
          ),
        ),
        const SizedBox(height: 32),
        const TNSkeleton(width: 150, height: 20),
        const SizedBox(height: 16),
        ...List.generate(3, (index) => const Padding(
          padding: EdgeInsets.only(bottom: 12.0),
          child: TNSkeleton(width: double.infinity, height: 80, borderRadius: 12),
        )),
      ],
    );
  }
}
