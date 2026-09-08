import 'package:flutter/material.dart';
import '../../../core/theme/app_theme.dart';
import 'contributor_badge.dart';

class ContributorProfileScreen extends StatelessWidget {
  final String userId;
  final String username;
  final String level;
  final int points;
  final int trustScore;
  final int approvedCount;
  final bool isFounding;

  const ContributorProfileScreen({
    super.key,
    required this.userId,
    this.username = 'TN Local Reporter',
    this.level = 'Trusted Contributor',
    this.points = 120,
    this.trustScore = 75,
    this.approvedCount = 12,
    this.isFounding = true,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(
          'Contributor Profile',
          style: TextStyle(fontWeight: FontWeight.bold, color: Colors.white),
        ),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // Header Profile Card
          Card(
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                children: [
                  const CircleAvatar(
                    radius: 36,
                    backgroundColor: AppTheme.darkOverlay,
                    child: Icon(Icons.person, size: 40, color: AppTheme.primaryOrange),
                  ),
                  const SizedBox(height: 12),
                  Text(
                    username,
                    style: const TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                      color: Colors.white,
                    ),
                  ),
                  const SizedBox(height: 8),
                  ContributorBadge(level: level, isFounding: isFounding),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),

          // Trust Score Meter Card
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text(
                        'Trust Score',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                      Text(
                        '$trustScore / 100',
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
                          color: AppTheme.primaryOrange,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 10),
                  ClipRRect(
                    borderRadius: BorderRadius.circular(8),
                    child: LinearProgressIndicator(
                      value: trustScore / 100,
                      minHeight: 8,
                      backgroundColor: AppTheme.darkOverlay,
                      valueColor: const AlwaysStoppedAnimation<Color>(AppTheme.primaryOrange),
                    ),
                  ),
                  const SizedBox(height: 8),
                  const Text(
                    'Higher trust scores unlock instant publishing privileges.',
                    style: TextStyle(fontSize: 12, color: Colors.grey),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),

          // Stats Grid Cards
          Row(
            children: [
              Expanded(
                child: _buildStatCard('Points', '$points', Icons.stars, Colors.orange),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildStatCard('Approved Posts', '$approvedCount', Icons.check_circle, Colors.green),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildStatCard(String title, String value, IconData icon, Color color) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Icon(icon, color: color, size: 28),
            const SizedBox(height: 8),
            Text(
              value,
              style: const TextStyle(
                fontSize: 22,
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              title,
              style: const TextStyle(fontSize: 12, color: Colors.grey),
            ),
          ],
        ),
      ),
    );
  }
}
