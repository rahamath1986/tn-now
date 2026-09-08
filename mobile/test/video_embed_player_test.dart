import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/features/video/presentation/video_embed_player.dart';

void main() {
  testWidgets('VideoEmbedPlayer renders thumbnail cover, play button, and platform badge', (WidgetTester tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: Scaffold(
          body: VideoEmbedPlayer(
            canonicalUrl: 'https://www.youtube.com/watch?v=dQw4w9WgXcQ',
            platform: 'youtube',
            embedHtml: '<iframe src="https://www.youtube.com/embed/dQw4w9WgXcQ"></iframe>',
            thumbnailUrl: 'https://img.youtube.com/vi/dQw4w9WgXcQ/hqdefault.jpg',
          ),
        ),
      ),
    );

    // Verify platform badge text
    expect(find.text('YOUTUBE'), findsOneWidget);

    // Verify external watch link button
    expect(find.text('Watch on App'), findsOneWidget);

    // Verify play button icon
    expect(find.byIcon(Icons.play_arrow_rounded), findsOneWidget);
  });
}
