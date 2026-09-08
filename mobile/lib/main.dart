import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'core/routing/app_router.dart';
import 'core/theme/app_theme.dart';
import 'core/localization/app_localizations.dart';

void main() {
  runApp(
    const ProviderScope(
      child: TNNowApp(),
    ),
  );
}

class TNNowApp extends ConsumerWidget {
  const TNNowApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      title: 'TN NOW',
      themeMode: ThemeMode.dark,
      darkTheme: AppTheme.darkTheme,
      routerConfig: router,
      debugShowCheckedModeBanner: false,
    );
  }
}

class SplashScreen extends ConsumerStatefulWidget {
  const SplashScreen({super.key});

  @override
  ConsumerState<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends ConsumerState<SplashScreen> {
  @override
  void initState() {
    super.initState();
    _startOnboardingFlow();
  }

  void _startOnboardingFlow() {
    Future.delayed(const Duration(seconds: 2), () {
      if (mounted) {
        // Route to language selection page for onboarding
        context.go('/language');
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final localizations = ref.watch(localizationsProvider);

    return Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              localizations.translate('app_title'),
              style: const TextStyle(
                fontSize: 36,
                fontWeight: FontWeight.bold,
                color: AppTheme.primaryOrange,
                letterSpacing: 2.0,
              ),
            ),
            const SizedBox(height: 10),
            Text(
              localizations.translate('tagline'),
              style: const TextStyle(
                fontSize: 14,
                color: Colors.grey,
              ),
            ),
            const SizedBox(height: 40),
            const CircularProgressIndicator(
              valueColor: AlwaysStoppedAnimation<Color>(AppTheme.primaryOrange),
            ),
          ],
        ),
      ),
    );
  }
}

