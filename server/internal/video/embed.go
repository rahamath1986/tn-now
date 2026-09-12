package video

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type VideoMetadata struct {
	Platform        string `json:"platform"`
	ExternalVideoID string `json:"externalVideoId"`
	CanonicalURL    string `json:"canonicalUrl"`
	EmbedHTML       string `json:"embedHtml"`
	ThumbnailURL    string `json:"thumbnailUrl"`
}

var (
	youtubeRegex   = regexp.MustCompile(`(?i)(?:youtube(?:-nocookie)?\.com/(?:watch\?.*v=|embed/|v/|shorts/|live/)|youtu\.be/)([\w-]{11})`)
	instagramRegex = regexp.MustCompile(`instagram\.com/(?:p|reel)/([\w-]+)`)
	facebookRegex  = regexp.MustCompile(`facebook\.com/(?:watch/\?v=|.*/videos/|v/)([\d]+)`)
	xTwitterRegex  = regexp.MustCompile(`(?:twitter\.com|x\.com)/.+/status/([\d]+)`)
)

// ParseVideoURL validates and extracts metadata from external video host URLs
func ParseVideoURL(rawURL string) (*VideoMetadata, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil, errors.New("URL cannot be empty")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return nil, errors.New("invalid URL format")
	}

	// 1. YouTube
	if matches := youtubeRegex.FindStringSubmatch(rawURL); len(matches) > 1 {
		videoID := matches[1]
		embedHTML := fmt.Sprintf(`<iframe width="100%%" height="100%%" src="https://www.youtube.com/embed/%s" frameborder="0" allowfullscreen></iframe>`, videoID)
		thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", videoID)

		return &VideoMetadata{
			Platform:        "youtube",
			ExternalVideoID: videoID,
			CanonicalURL:    fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
			EmbedHTML:       embedHTML,
			ThumbnailURL:    thumbnailURL,
		}, nil
	}

	// 2. Instagram
	if matches := instagramRegex.FindStringSubmatch(rawURL); len(matches) > 1 {
		postID := matches[1]
		embedHTML := fmt.Sprintf(`<iframe src="https://www.instagram.com/p/%s/embed" width="100%%" height="100%%" frameborder="0" scrolling="no" allowtransparency="true"></iframe>`, postID)

		return &VideoMetadata{
			Platform:        "instagram",
			ExternalVideoID: postID,
			CanonicalURL:    fmt.Sprintf("https://www.instagram.com/p/%s/", postID),
			EmbedHTML:       embedHTML,
			ThumbnailURL:    "",
		}, nil
	}

	// 3. Facebook
	if matches := facebookRegex.FindStringSubmatch(rawURL); len(matches) > 1 {
		videoID := matches[1]
		encodedURL := url.QueryEscape(rawURL)
		embedHTML := fmt.Sprintf(`<iframe src="https://www.facebook.com/plugins/video.php?href=%s&show_text=false" width="100%%" height="100%%" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowfullscreen="true"></iframe>`, encodedURL)

		return &VideoMetadata{
			Platform:        "facebook",
			ExternalVideoID: videoID,
			CanonicalURL:    rawURL,
			EmbedHTML:       embedHTML,
			ThumbnailURL:    "",
		}, nil
	}

	// 4. X / Twitter
	if matches := xTwitterRegex.FindStringSubmatch(rawURL); len(matches) > 1 {
		statusID := matches[1]
		embedHTML := fmt.Sprintf(`<blockquote class="twitter-tweet"><a href="https://twitter.com/x/status/%s"></a></blockquote><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>`, statusID)

		return &VideoMetadata{
			Platform:        "x_twitter",
			ExternalVideoID: statusID,
			CanonicalURL:    fmt.Sprintf("https://x.com/i/status/%s", statusID),
			EmbedHTML:       embedHTML,
			ThumbnailURL:    "",
		}, nil
	}

	return nil, errors.New("unsupported video platform. Only YouTube, Instagram, Facebook, and X video links are allowed")
}
