import 'package:flutter/material.dart';
import '../../../core/theme/app_theme.dart';

class ContributorBadge extends StatelessWidget {
  final String level;
  final bool isFounding;

  const ContributorBadge({
    super.key,
    required this.level,
    this.isFounding = false,
  });

  @override
  Widget build(BuildContext context) {
    return Wrap(
      spacing: 6,
      runSpacing: 4,
      children: [
        // Level Badge
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
          decoration: BoxDecoration(
            color: _getLevelColor(level).withValues(alpha: 0.2),
            borderRadius: BorderRadius.circular(12),
            border: Border.all(color: _getLevelColor(level)),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(_getLevelIcon(level), size: 12, color: _getLevelColor(level)),
              const SizedBox(width: 4),
              Text(
                level,
                style: TextStyle(
                  fontSize: 11,
                  fontWeight: FontWeight.bold,
                  color: _getLevelColor(level),
                ),
              ),
            ],
          ),
        ),

        // Special Founding Contributor Badge
        if (isFounding)
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [Color(0xFFFFD700), Color(0xFFFFA500)],
              ),
              borderRadius: BorderRadius.circular(12),
              boxShadow: [
                BoxShadow(
                  color: const Color(0xFFFFD700).withValues(alpha: 0.4),
                  blurRadius: 6,
                  spreadRadius: 1,
                ),
              ],
            ),
            child: const Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(Icons.workspace_premium, size: 12, color: Colors.black),
                SizedBox(width: 4),
                Text(
                  'Founding 2026',
                  style: TextStyle(
                    fontSize: 11,
                    fontWeight: FontWeight.bold,
                    color: Colors.black,
                  ),
                ),
              ],
            ),
          ),
      ],
    );
  }

  Color _getLevelColor(String level) {
    switch (level) {
      case 'Core Contributor':
        return Colors.purpleAccent;
      case 'Verified Contributor':
        return Colors.blueAccent;
      case 'Trusted Contributor':
        return Colors.green;
      default:
        return AppTheme.primaryOrange;
    }
  }

  IconData _getLevelIcon(String level) {
    switch (level) {
      case 'Core Contributor':
        return Icons.stars;
      case 'Verified Contributor':
        return Icons.verified;
      case 'Trusted Contributor':
        return Icons.shield;
      default:
        return Icons.person;
    }
  }
}
