import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class AppLocalizations {
  final Locale locale;

  AppLocalizations(this.locale);

  static const Map<String, Map<String, String>> _localizedValues = {
    'en': {
      'app_title': 'TN NOW',
      'tagline': 'Everything Tamil Nadu, Right Now.',
      'login': 'Login',
      'register': 'Register',
      'language_select': 'Select Language',
      'tamil': 'தமிழ்',
      'english': 'English',
      'home': 'Home',
      'trending': 'Trending',
      'nearby': 'Nearby',
      'events': 'Events',
      'profile': 'Profile',
      'submit': 'Submit',
      'madurai_launch': 'Welcome to TN NOW\n(Madurai Launch)',
    },
    'ta': {
      'app_title': 'டிஎன் நவ்',
      'tagline': 'தமிழ்நாட்டின் அண்மைச் செய்திகள், உடனுக்குடன்.',
      'login': 'உள்நுழை',
      'register': 'பதிவுசெய்',
      'language_select': 'மொழியைத் தேர்வுசெய்க',
      'tamil': 'தமிழ்',
      'english': 'English',
      'home': 'முகப்பு',
      'trending': 'பிரபலமானவை',
      'nearby': 'அருகிலுள்ளவை',
      'events': 'நிகழ்வுகள்',
      'profile': 'சுயவிவரம்',
      'submit': 'பகிர்க',
      'madurai_launch': 'டிஎன் நவ்-க்கு வரவேற்கிறோம்\n(மதுரை அறிமுகம்)',
    },
  };

  String translate(String key) {
    return _localizedValues[locale.languageCode]?[key] ?? key;
  }
}

// Riverpod locale state provider
final localeProvider = StateNotifierProvider<LocaleNotifier, Locale>((ref) {
  return LocaleNotifier();
});

class LocaleNotifier extends StateNotifier<Locale> {
  LocaleNotifier() : super(const Locale('en'));

  void setLocale(Locale locale) {
    if (locale.languageCode == 'ta' || locale.languageCode == 'en') {
      state = locale;
    }
  }

  void toggleLocale() {
    state = state.languageCode == 'en' ? const Locale('ta') : const Locale('en');
  }
}

// AppLocalizations provider for easy watching in components
final localizationsProvider = Provider<AppLocalizations>((ref) {
  final locale = ref.watch(localeProvider);
  return AppLocalizations(locale);
});
