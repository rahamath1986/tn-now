import 'package:flutter/material.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/widgets/design_system.dart';

class GrievanceScreen extends StatefulWidget {
  const GrievanceScreen({super.key});

  @override
  State<GrievanceScreen> createState() => _GrievanceScreenState();
}

class _GrievanceScreenState extends State<GrievanceScreen> {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _contentIdController = TextEditingController();
  final _descriptionController = TextEditingController();

  String? _selectedCategory;
  bool _isSubmitting = false;

  final List<Map<String, String>> _categories = const [
    {'id': 'COPYRIGHT', 'name': 'Copyright / IP Infringement'},
    {'id': 'DEFAMATION', 'name': 'Defamation / Misinformation'},
    {'id': 'PRIVACY_VIOLATION', 'name': 'Privacy / PII Leak'},
    {'id': 'OBJECTIONABLE_CONTENT', 'name': 'Objectionable Content'},
    {'id': 'IT_ACT_VIOLATION', 'name': 'IT Act Violation'},
  ];

  @override
  void dispose() {
    _emailController.dispose();
    _contentIdController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }

  void _handleSubmit() async {
    if (!_formKey.currentState!.validate() || _selectedCategory == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please complete all required fields')),
      );
      return;
    }

    setState(() => _isSubmitting = true);
    await Future.delayed(const Duration(seconds: 1));

    if (!mounted) return;
    setState(() => _isSubmitting = false);

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Row(
          children: [
            Icon(Icons.verified, color: Colors.green),
            SizedBox(width: 8),
            Text('Grievance Registered'),
          ],
        ),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Your grievance has been formally registered under Indian IT Rules 2021.',
              style: TextStyle(fontSize: 14),
            ),
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: AppTheme.darkOverlay,
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('Ticket Ref: GRV-2026-8942', style: TextStyle(fontWeight: FontWeight.bold, color: AppTheme.primaryOrange)),
                  SizedBox(height: 4),
                  Text('• Acknowledgment SLA: 24 Hours', style: TextStyle(fontSize: 12, color: Colors.grey)),
                  Text('• Resolution SLA: 15 Days', style: TextStyle(fontSize: 12, color: Colors.grey)),
                ],
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              _emailController.clear();
              _contentIdController.clear();
              _descriptionController.clear();
            },
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(
          'Grievance Redressal',
          style: TextStyle(fontWeight: FontWeight.bold, color: Colors.white),
        ),
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            // Mandatory Disclosure Card
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Row(
                      children: [
                        Icon(Icons.gavel, color: AppTheme.primaryOrange),
                        SizedBox(width: 8),
                        Text(
                          'IT Rules 2021 Compliance Officer',
                          style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Colors.white),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    const Text(
                      'Nodal Grievance Redressal Officer\nTN NOW Media Platform, Tech Park, Madurai, TN 625020\nEmail: grievance@tnnow.in',
                      style: TextStyle(fontSize: 13, color: Colors.grey, height: 1.4),
                    ),
                    const SizedBox(height: 8),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                      decoration: BoxDecoration(
                        color: Colors.blue.withValues(alpha: 0.2),
                        borderRadius: BorderRadius.circular(4),
                      ),
                      child: const Text(
                        '⏱️ 24h Receipt Acknowledgment | 15d Final Resolution SLA',
                        style: TextStyle(fontSize: 11, color: Colors.blue, fontWeight: FontWeight.bold),
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 20),

            // Email Input
            TextFormField(
              controller: _emailController,
              decoration: const InputDecoration(
                labelText: 'Complainant Email *',
                hintText: 'name@example.com',
                prefixIcon: Icon(Icons.email_outlined),
              ),
              validator: (v) => v == null || v.trim().isEmpty ? 'Email is required' : null,
            ),
            const SizedBox(height: 16),

            // Category Picker
            DropdownButtonFormField<String>(
              initialValue: _selectedCategory,
              decoration: const InputDecoration(
                labelText: 'Violation Category *',
                prefixIcon: Icon(Icons.category_outlined),
              ),
              items: _categories
                  .map((c) => DropdownMenuItem(
                        value: c['id'],
                        child: Text(c['name']!),
                      ))
                  .toList(),
              onChanged: (val) => setState(() => _selectedCategory = val),
              validator: (v) => v == null ? 'Required' : null,
            ),
            const SizedBox(height: 16),

            // Content ID / Link Input
            TextFormField(
              controller: _contentIdController,
              decoration: const InputDecoration(
                labelText: 'Content Link or ID',
                hintText: 'https://tnnow.in/content/...',
                prefixIcon: Icon(Icons.link_outlined),
              ),
            ),
            const SizedBox(height: 16),

            // Description Input
            TextFormField(
              controller: _descriptionController,
              maxLines: 4,
              decoration: const InputDecoration(
                labelText: 'Detailed Complaint Description *',
                hintText: 'Explain the nature of the violation or grievance...',
              ),
              validator: (v) => v == null || v.trim().isEmpty ? 'Description is required' : null,
            ),
            const SizedBox(height: 24),

            // Submit Button
            TNPrimaryButton(
              label: 'Register Grievance',
              isLoading: _isSubmitting,
              onPressed: _handleSubmit,
            ),
          ],
        ),
      ),
    );
  }
}
