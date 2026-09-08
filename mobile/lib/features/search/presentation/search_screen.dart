import 'package:flutter/material.dart';
import '../../../core/theme/app_theme.dart';

class SearchScreen extends StatefulWidget {
  const SearchScreen({super.key});

  @override
  State<SearchScreen> createState() => _SearchScreenState();
}

class _SearchScreenState extends State<SearchScreen> {
  final _searchController = TextEditingController();
  String _selectedDistrict = 'All';
  List<Map<String, String>> _results = [];
  bool _isSearching = false;

  final List<String> _districts = const [
    'All',
    'Madurai',
    'Chennai',
    'Coimbatore',
    'Trichy',
    'Salem',
  ];

  final List<Map<String, String>> _dummyData = const [
    {
      'id': 'res-1',
      'title': 'Madurai Meenakshi Temple Annual Chithirai Festival 2026',
      'description': 'Grand procession and cultural gathering scheduled for next week in Madurai city center.',
      'contentType': 'EVENT',
      'district': 'Madurai',
    },
    {
      'id': 'res-2',
      'title': 'New Highway Expansion Project Update - Chennai to Trichy',
      'description': 'Infrastructure milestone announced by state highways department today.',
      'contentType': 'TEXT_STORY',
      'district': 'Chennai',
    },
    {
      'id': 'res-3',
      'title': 'Coimbatore Tech Park Opening Highlights',
      'description': 'Viral video walkthrough of the newly inaugurated software hub in Coimbatore.',
      'contentType': 'VIDEO_LINK',
      'district': 'Coimbatore',
    },
  ];

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _onSearchChanged(String query) {
    if (query.trim().isEmpty) {
      setState(() {
        _results = [];
        _isSearching = false;
      });
      return;
    }

    setState(() => _isSearching = true);

    final filtered = _dummyData.where((item) {
      final matchesQuery = item['title']!.toLowerCase().contains(query.toLowerCase()) ||
          item['description']!.toLowerCase().contains(query.toLowerCase());
      final matchesDistrict = _selectedDistrict == 'All' || item['district'] == _selectedDistrict;
      return matchesQuery && matchesDistrict;
    }).toList();

    setState(() {
      _results = filtered;
      _isSearching = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: TextField(
          controller: _searchController,
          autofocus: true,
          style: const TextStyle(color: Colors.white),
          decoration: InputDecoration(
            hintText: 'Search updates across Tamil Nadu...',
            hintStyle: const TextStyle(color: Colors.grey),
            border: InputBorder.none,
            suffixIcon: _searchController.text.isNotEmpty
                ? IconButton(
                    icon: const Icon(Icons.clear, color: Colors.grey),
                    onPressed: () {
                      _searchController.clear();
                      _onSearchChanged('');
                    },
                  )
                : null,
          ),
          onChanged: _onSearchChanged,
        ),
      ),
      body: Column(
        children: [
          // District Filter Chips Bar
          SizedBox(
            height: 48,
            child: ListView.builder(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 16),
              itemCount: _districts.length,
              itemBuilder: (context, index) {
                final d = _districts[index];
                final isSelected = d == _selectedDistrict;
                return Padding(
                  padding: const EdgeInsets.only(right: 8),
                  child: ChoiceChip(
                    label: Text(d),
                    selected: isSelected,
                    selectedColor: AppTheme.primaryOrange,
                    onSelected: (val) {
                      if (val) {
                        setState(() => _selectedDistrict = d);
                        _onSearchChanged(_searchController.text);
                      }
                    },
                  ),
                );
              },
            ),
          ),
          const Divider(color: Colors.grey),

          // Search Results
          Expanded(
            child: _isSearching
                ? const Center(child: CircularProgressIndicator(color: AppTheme.primaryOrange))
                : _searchController.text.isEmpty
                    ? const Center(
                        child: Text(
                          'Type a keyword to discover Tamil Nadu updates 🔍',
                          style: TextStyle(color: Colors.grey),
                        ),
                      )
                    : _results.isEmpty
                        ? const Center(
                            child: Text(
                              'No matching updates found.',
                              style: TextStyle(color: Colors.grey),
                            ),
                          )
                        : ListView.builder(
                            padding: const EdgeInsets.all(16),
                            itemCount: _results.length,
                            itemBuilder: (context, index) {
                              final item = _results[index];
                              return Card(
                                margin: const EdgeInsets.only(bottom: 12),
                                child: Padding(
                                  padding: const EdgeInsets.all(16),
                                  child: Column(
                                    crossAxisAlignment: CrossAxisAlignment.start,
                                    children: [
                                      Row(
                                        children: [
                                          Container(
                                            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                                            decoration: BoxDecoration(
                                              color: AppTheme.primaryOrange.withValues(alpha: 0.2),
                                              borderRadius: BorderRadius.circular(4),
                                            ),
                                            child: Text(
                                              item['contentType']!,
                                              style: const TextStyle(
                                                fontSize: 10,
                                                color: AppTheme.primaryOrange,
                                                fontWeight: FontWeight.bold,
                                              ),
                                            ),
                                          ),
                                          const Spacer(),
                                          Text(
                                            item['district']!,
                                            style: const TextStyle(fontSize: 12, color: Colors.grey),
                                          ),
                                        ],
                                      ),
                                      const SizedBox(height: 8),
                                      Text(
                                        item['title']!,
                                        style: const TextStyle(
                                          fontSize: 16,
                                          fontWeight: FontWeight.bold,
                                          color: Colors.white,
                                        ),
                                      ),
                                      const SizedBox(height: 4),
                                      Text(
                                        item['description']!,
                                        style: const TextStyle(fontSize: 13, color: Colors.grey),
                                      ),
                                    ],
                                  ),
                                ),
                              );
                            },
                          ),
          ),
        ],
      ),
    );
  }
}
