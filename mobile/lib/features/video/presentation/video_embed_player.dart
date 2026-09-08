import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:webview_flutter/webview_flutter.dart';
import '../../../core/theme/app_theme.dart';

class VideoEmbedPlayer extends StatefulWidget {
  final String canonicalUrl;
  final String platform; // 'youtube' | 'instagram' | 'facebook' | 'x_twitter'
  final String? embedHtml;
  final String? thumbnailUrl;

  const VideoEmbedPlayer({
    super.key,
    required this.canonicalUrl,
    required this.platform,
    this.embedHtml,
    this.thumbnailUrl,
  });

  @override
  State<VideoEmbedPlayer> createState() => _VideoEmbedPlayerState();
}

class _VideoEmbedPlayerState extends State<VideoEmbedPlayer> {
  bool _isPlaying = false;
  WebViewController? _webViewController;

  @override
  void initState() {
    super.initState();
  }

  void _startPlayback() {
    if (widget.embedHtml != null && widget.embedHtml!.isNotEmpty) {
      final controller = WebViewController()
        ..setJavaScriptMode(JavaScriptMode.unrestricted)
        ..setBackgroundColor(Colors.black)
        ..loadHtmlString('''
          <!DOCTYPE html>
          <html>
            <head>
              <meta name="viewport" content="width=device-width, initial-scale=1.0">
              <style>
                body { margin: 0; padding: 0; background-color: #000; display: flex; justify-content: center; align-items: center; height: 100vh; }
                iframe { width: 100%; height: 100%; border: none; }
              </style>
            </head>
            <body>
              ${widget.embedHtml}
            </body>
          </html>
        ''');
      setState(() {
        _webViewController = controller;
        _isPlaying = true;
      });
    } else {
      _openExternalUrl();
    }
  }

  Future<void> _openExternalUrl() async {
    final uri = Uri.parse(widget.canonicalUrl);
    if (await canLaunchUrl(uri)) {
      await launchUrl(uri, mode: LaunchMode.externalApplication);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      height: 220,
      decoration: BoxDecoration(
        color: AppTheme.darkSurface,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: AppTheme.darkOverlay),
      ),
      clipBehavior: Clip.antiAlias,
      child: _isPlaying && _webViewController != null
          ? Stack(
              children: [
                WebViewWidget(controller: _webViewController!),
                Positioned(
                  top: 8,
                  right: 8,
                  child: IconButton(
                    icon: const Icon(Icons.close, color: Colors.white),
                    onPressed: () {
                      setState(() {
                        _isPlaying = false;
                        _webViewController = null;
                      });
                    },
                  ),
                ),
              ],
            )
          : Stack(
              fit: StackFit.expand,
              children: [
                // Thumbnail image or platform gradient placeholder
                if (widget.thumbnailUrl != null && widget.thumbnailUrl!.isNotEmpty)
                  CachedNetworkImage(
                    imageUrl: widget.thumbnailUrl!,
                    fit: BoxFit.cover,
                    placeholder: (context, url) => Container(color: AppTheme.darkSurface),
                    errorWidget: (context, url, error) => _buildFallbackCover(),
                  )
                else
                  _buildFallbackCover(),

                // Dark overlay gradient
                Container(
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      begin: Alignment.topCenter,
                      end: Alignment.bottomCenter,
                      colors: [
                        Colors.transparent,
                        Colors.black.withValues(alpha: 0.7),
                      ],
                    ),
                  ),
                ),

                // Play Button Center Overlay
                Center(
                  child: GestureDetector(
                    onTap: _startPlayback,
                    child: Container(
                      padding: const EdgeInsets.all(12),
                      decoration: BoxDecoration(
                        color: AppTheme.primaryOrange.withValues(alpha: 0.9),
                        shape: BoxShape.circle,
                        boxShadow: [
                          BoxShadow(
                            color: AppTheme.primaryOrange.withValues(alpha: 0.4),
                            blurRadius: 16,
                            spreadRadius: 2,
                          ),
                        ],
                      ),
                      child: const Icon(
                        Icons.play_arrow_rounded,
                        color: Colors.white,
                        size: 36,
                      ),
                    ),
                  ),
                ),

                // Platform Badge and External Link Button Bottom Bar
                Positioned(
                  left: 12,
                  bottom: 12,
                  right: 12,
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        decoration: BoxDecoration(
                          color: Colors.black.withValues(alpha: 0.7),
                          borderRadius: BorderRadius.circular(6),
                        ),
                        child: Row(
                          children: [
                            _getPlatformIcon(widget.platform),
                            const SizedBox(width: 6),
                            Text(
                              widget.platform.toUpperCase(),
                              style: const TextStyle(
                                color: Colors.white,
                                fontSize: 10,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ],
                        ),
                      ),
                      GestureDetector(
                        onTap: _openExternalUrl,
                        child: Container(
                          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                          decoration: BoxDecoration(
                            color: Colors.black.withValues(alpha: 0.7),
                            borderRadius: BorderRadius.circular(6),
                          ),
                          child: const Row(
                            children: [
                              Text(
                                'Watch on App',
                                style: TextStyle(
                                  color: Colors.white70,
                                  fontSize: 10,
                                ),
                              ),
                              SizedBox(width: 4),
                              Icon(Icons.open_in_new, size: 10, color: Colors.white70),
                            ],
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
    );
  }

  Widget _buildFallbackCover() {
    return Container(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            AppTheme.primaryOrange.withValues(alpha: 0.3),
            AppTheme.darkSurface,
          ],
        ),
      ),
      child: Center(
        child: _getPlatformIcon(widget.platform, size: 48),
      ),
    );
  }

  Widget _getPlatformIcon(String platform, {double size = 16}) {
    switch (platform.toLowerCase()) {
      case 'youtube':
        return Icon(Icons.play_circle_fill, color: Colors.red, size: size);
      case 'instagram':
        return Icon(Icons.camera_alt, color: Colors.purpleAccent, size: size);
      case 'facebook':
        return Icon(Icons.facebook, color: Colors.blue, size: size);
      case 'x_twitter':
        return Icon(Icons.alternate_email, color: Colors.white, size: size);
      default:
        return Icon(Icons.video_library, color: AppTheme.primaryOrange, size: size);
    }
  }
}
