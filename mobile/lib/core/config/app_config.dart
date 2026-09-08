enum AppEnvironment {
  development,
  staging,
  production,
}

class AppConfig {
  final AppEnvironment environment;
  final String apiBaseUrl;
  final Duration connectTimeout;
  final Duration receiveTimeout;

  const AppConfig({
    required this.environment,
    required this.apiBaseUrl,
    this.connectTimeout = const Duration(seconds: 10),
    this.receiveTimeout = const Duration(seconds: 10),
  });

  static AppConfig? _instance;

  static void initialize(AppEnvironment env) {
    switch (env) {
      case AppEnvironment.development:
        _instance = const AppConfig(
          environment: AppEnvironment.development,
          // Using 10.0.2.2 for local Android emulator, localhost for others
          apiBaseUrl: 'http://localhost:8080/api/v1',
        );
        break;
      case AppEnvironment.staging:
        _instance = const AppConfig(
          environment: AppEnvironment.staging,
          apiBaseUrl: 'https://staging.tnnow.in/api/v1',
        );
        break;
      case AppEnvironment.production:
        _instance = const AppConfig(
          environment: AppEnvironment.production,
          apiBaseUrl: 'https://api.tnnow.in/api/v1',
        );
        break;
    }
  }

  static AppConfig get instance {
    if (_instance == null) {
      // Default fallback
      initialize(AppEnvironment.development);
    }
    return _instance!;
  }
}
