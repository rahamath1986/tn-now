import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

// Import Screens to be defined
import '../../../main.dart'; // fallback stubs for now
import '../../features/home/presentation/home_screen.dart';
import '../../features/category/presentation/category_detail_screen.dart';
import '../../features/submission/presentation/submission_screen.dart';
import '../../features/moderation/presentation/moderation_queue_screen.dart';
import '../../features/reputation/presentation/contributor_profile_screen.dart';
import '../../features/compliance/presentation/grievance_screen.dart';
import '../../features/search/presentation/search_screen.dart';

final routerProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: '/',
    debugLogDiagnostics: true,
    routes: [
      GoRoute(
        path: '/',
        name: 'splash',
        builder: (context, state) => const SplashScreen(),
      ),
      GoRoute(
        path: '/language',
        name: 'language',
        builder: (context, state) => const LanguageScreen(),
      ),
      GoRoute(
        path: '/location',
        name: 'location',
        builder: (context, state) => const LocationOnboardingScreen(),
      ),
      GoRoute(
        path: '/login',
        name: 'login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/home',
        name: 'home',
        builder: (context, state) => const HomeScreen(),
      ),
      GoRoute(
        path: '/category/:id',
        name: 'category',
        builder: (context, state) {
          final id = state.pathParameters['id'] ?? '';
          final name = state.uri.queryParameters['name'];
          return CategoryDetailScreen(categoryId: id, categoryName: name);
        },
      ),
      GoRoute(
        path: '/submit',
        name: 'submit',
        builder: (context, state) => const SubmissionScreen(),
      ),
      GoRoute(
        path: '/moderation',
        name: 'moderation',
        builder: (context, state) => const ModerationQueueScreen(),
      ),
      GoRoute(
        path: '/contributor/:id',
        name: 'contributor',
        builder: (context, state) {
          final id = state.pathParameters['id'] ?? '';
          return ContributorProfileScreen(userId: id);
        },
      ),
      GoRoute(
        path: '/grievance',
        name: 'grievance',
        builder: (context, state) => const GrievanceScreen(),
      ),
      GoRoute(
        path: '/search',
        name: 'search',
        builder: (context, state) => const SearchScreen(),
      ),
    ],
  );
});

// Mock screen stubs representing standard flows for routing
class LanguageScreen extends StatelessWidget {
  const LanguageScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(child: Text('Language Screen')),
    );
  }
}

class LocationOnboardingScreen extends StatelessWidget {
  const LocationOnboardingScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(child: Text('Location Onboarding')),
    );
  }
}

class LoginScreen extends StatelessWidget {
  const LoginScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(child: Text('Login')),
    );
  }
}
