import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/widgets/design_system.dart';

class SubmissionScreen extends ConsumerStatefulWidget {
  const SubmissionScreen({super.key});

  @override
  ConsumerState<SubmissionScreen> createState() => _SubmissionScreenState();
}

class _SubmissionScreenState extends ConsumerState<SubmissionScreen>
    with SingleTickerProviderStateMixin {
  late final TabController _tabController;
  final _formKey = GlobalKey<FormState>();

  // Shared form controllers
  final _titleController = TextEditingController();
  final _descriptionController = TextEditingController();
  String? _selectedCategory;
  String? _selectedDistrict;

  // Video specific controllers
  final _videoUrlController = TextEditingController();

  // Photo specific controllers
  final _photoUrlController = TextEditingController();

  // Story specific controllers
  final _headlineController = TextEditingController();
  final _storyBodyController = TextEditingController();

  // Event specific controllers
  final _eventNameController = TextEditingController();
  final _organizerNameController = TextEditingController();

  bool _isSubmitting = false;

  final List<Map<String, String>> _categories = const [
    {'id': 'cat-viral', 'name': 'Viral Videos'},
    {'id': 'cat-news', 'name': 'News'},
    {'id': 'cat-politics', 'name': 'Politics'},
    {'id': 'cat-cinema', 'name': 'Cinema'},
    {'id': 'cat-sports', 'name': 'Sports'},
    {'id': 'cat-events', 'name': 'Events'},
  ];

  final List<Map<String, String>> _districts = const [
    {'id': 'dist-madurai', 'name': 'Madurai'},
    {'id': 'dist-chennai', 'name': 'Chennai'},
    {'id': 'dist-coimbatore', 'name': 'Coimbatore'},
    {'id': 'dist-trichy', 'name': 'Trichy'},
    {'id': 'dist-salem', 'name': 'Salem'},
    {'id': 'dist-tn', 'name': 'Tamil Nadu'},
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    _titleController.dispose();
    _descriptionController.dispose();
    _videoUrlController.dispose();
    _photoUrlController.dispose();
    _headlineController.dispose();
    _storyBodyController.dispose();
    _eventNameController.dispose();
    _organizerNameController.dispose();
    super.dispose();
  }

  Future<void> _handleSubmit() async {
    if (!_formKey.currentState!.validate()) return;
    if (_selectedCategory == null || _selectedDistrict == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please select Category and District')),
      );
      return;
    }

    setState(() => _isSubmitting = true);

    // Simulate API submission submission delay
    await Future.delayed(const Duration(seconds: 1));

    if (!mounted) return;
    setState(() => _isSubmitting = false);

    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Post submitted successfully! Queued for moderation.'),
        backgroundColor: Colors.green,
      ),
    );

    // Clear form
    _titleController.clear();
    _descriptionController.clear();
    _videoUrlController.clear();
    _photoUrlController.clear();
    _storyBodyController.clear();
    _eventNameController.clear();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(
          'Submit Update',
          style: TextStyle(fontWeight: FontWeight.bold, color: Colors.white),
        ),
        bottom: TabBar(
          controller: _tabController,
          indicatorColor: AppTheme.primaryOrange,
          labelColor: AppTheme.primaryOrange,
          unselectedLabelColor: Colors.grey,
          tabs: const [
            Tab(icon: Icon(Icons.video_library), text: 'Video'),
            Tab(icon: Icon(Icons.photo_library), text: 'Photo'),
            Tab(icon: Icon(Icons.article), text: 'Story'),
            Tab(icon: Icon(Icons.event), text: 'Event'),
          ],
        ),
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            // Title Input
            TextFormField(
              controller: _titleController,
              decoration: const InputDecoration(
                labelText: 'Title *',
                hintText: 'Enter update headline...',
              ),
              validator: (v) =>
                  v == null || v.trim().isEmpty ? 'Title is required' : null,
            ),
            const SizedBox(height: 16),

            // Description Input
            TextFormField(
              controller: _descriptionController,
              maxLines: 2,
              decoration: const InputDecoration(
                labelText: 'Description',
                hintText: 'Add context or details...',
              ),
            ),
            const SizedBox(height: 16),

            // Category & District Dropdowns
            Row(
              children: [
                Expanded(
                  child: DropdownButtonFormField<String>(
                    initialValue: _selectedCategory,
                    decoration: const InputDecoration(labelText: 'Category *'),
                    items: _categories
                        .map((c) => DropdownMenuItem(
                              value: c['id'],
                              child: Text(c['name']!),
                            ))
                        .toList(),
                    onChanged: (val) => setState(() => _selectedCategory = val),
                    validator: (v) => v == null ? 'Required' : null,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: DropdownButtonFormField<String>(
                    initialValue: _selectedDistrict,
                    decoration: const InputDecoration(labelText: 'District *'),
                    items: _districts
                        .map((d) => DropdownMenuItem(
                              value: d['id'],
                              child: Text(d['name']!),
                            ))
                        .toList(),
                    onChanged: (val) => setState(() => _selectedDistrict = val),
                    validator: (v) => v == null ? 'Required' : null,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 24),

            // Format Specific Fields
            SizedBox(
              height: 250,
              child: TabBarView(
                controller: _tabController,
                children: [
                  // Tab 1: Video Link
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      TextFormField(
                        controller: _videoUrlController,
                        decoration: const InputDecoration(
                          labelText: 'Video URL *',
                          hintText: 'YouTube, Instagram Reel, Facebook, or X link',
                          prefixIcon: Icon(Icons.link, color: AppTheme.primaryOrange),
                        ),
                        validator: (v) {
                          if (_tabController.index == 0 && (v == null || v.trim().isEmpty)) {
                            return 'Video URL is required';
                          }
                          return null;
                        },
                      ),
                      const SizedBox(height: 12),
                      const Text(
                        '💡 Supported hosts: YouTube, Instagram Reels, Facebook Videos, and X Media posts.',
                        style: TextStyle(fontSize: 12, color: Colors.grey),
                      ),
                    ],
                  ),

                  // Tab 2: Photo Post
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      TextFormField(
                        controller: _photoUrlController,
                        decoration: const InputDecoration(
                          labelText: 'Photo Image URL *',
                          hintText: 'https://...',
                          prefixIcon: Icon(Icons.image, color: AppTheme.primaryOrange),
                        ),
                        validator: (v) {
                          if (_tabController.index == 1 && (v == null || v.trim().isEmpty)) {
                            return 'Photo URL is required';
                          }
                          return null;
                        },
                      ),
                    ],
                  ),

                  // Tab 3: Text Story
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      TextFormField(
                        controller: _storyBodyController,
                        maxLines: 4,
                        decoration: const InputDecoration(
                          labelText: 'Story Body Content *',
                          hintText: 'Write complete news report or story details...',
                        ),
                        validator: (v) {
                          if (_tabController.index == 2 && (v == null || v.trim().isEmpty)) {
                            return 'Story body content is required';
                          }
                          return null;
                        },
                      ),
                    ],
                  ),

                  // Tab 4: Event
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      TextFormField(
                        controller: _eventNameController,
                        decoration: const InputDecoration(
                          labelText: 'Event Name *',
                          hintText: 'e.g. Madurai Annual Cultural Meet',
                        ),
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: _organizerNameController,
                        decoration: const InputDecoration(
                          labelText: 'Organizer Name *',
                          hintText: 'e.g. City Cultural Association',
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),

            // Submit Button
            TNPrimaryButton(
              label: 'Submit for Moderation',
              isLoading: _isSubmitting,
              onPressed: _handleSubmit,
            ),
          ],
        ),
      ),
    );
  }
}
