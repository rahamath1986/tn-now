import 'package:flutter/material.dart';
import '../../../core/theme/app_theme.dart';
import 'comments_sheet.dart';

class EngagementBar extends StatefulWidget {
  final String contentId;
  final int initialLikeCount;
  final int initialCommentCount;
  final bool initialIsLiked;
  final bool initialIsBookmarked;

  const EngagementBar({
    super.key,
    required this.contentId,
    this.initialLikeCount = 24,
    this.initialCommentCount = 5,
    this.initialIsLiked = false,
    this.initialIsBookmarked = false,
  });

  @override
  State<EngagementBar> createState() => _EngagementBarState();
}

class _EngagementBarState extends State<EngagementBar> {
  late int _likeCount;
  late int _commentCount;
  late bool _isLiked;
  late bool _isBookmarked;

  @override
  void initState() {
    super.initState();
    _likeCount = widget.initialLikeCount;
    _commentCount = widget.initialCommentCount;
    _isLiked = widget.initialIsLiked;
    _isBookmarked = widget.initialIsBookmarked;
  }

  void _toggleLike() {
    setState(() {
      _isLiked = !_isLiked;
      if (_isLiked) {
        _likeCount++;
      } else {
        _likeCount--;
      }
    });
  }

  void _toggleBookmark() {
    setState(() {
      _isBookmarked = !_isBookmarked;
    });

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(_isBookmarked ? 'Saved to Bookmarks' : 'Removed from Bookmarks'),
        duration: const Duration(seconds: 1),
      ),
    );
  }

  void _openComments() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => CommentsSheet(
        contentId: widget.contentId,
        onCommentAdded: () {
          setState(() {
            _commentCount++;
          });
        },
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        // Like Button
        IconButton(
          icon: Icon(
            _isLiked ? Icons.favorite : Icons.favorite_border,
            color: _isLiked ? Colors.red : Colors.grey,
          ),
          onPressed: _toggleLike,
        ),
        Text(
          '$_likeCount',
          style: const TextStyle(color: Colors.white, fontSize: 13),
        ),
        const SizedBox(width: 16),

        // Comment Button
        IconButton(
          icon: const Icon(Icons.mode_comment_outlined, color: Colors.grey),
          onPressed: _openComments,
        ),
        Text(
          '$_commentCount',
          style: const TextStyle(color: Colors.white, fontSize: 13),
        ),
        const Spacer(),

        // Bookmark Button
        IconButton(
          icon: Icon(
            _isBookmarked ? Icons.bookmark : Icons.bookmark_border,
            color: _isBookmarked ? AppTheme.primaryOrange : Colors.grey,
          ),
          onPressed: _toggleBookmark,
        ),
      ],
    );
  }
}
