package video

import (
	"testing"
)

func TestParseVideoURL(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		wantPlatform string
		wantID       string
		wantErr      bool
	}{
		{
			name:         "YouTube Watch URL",
			url:          "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			wantPlatform: "youtube",
			wantID:       "dQw4w9WgXcQ",
			wantErr:      false,
		},
		{
			name:         "YouTube Short URL",
			url:          "https://youtu.be/dQw4w9WgXcQ",
			wantPlatform: "youtube",
			wantID:       "dQw4w9WgXcQ",
			wantErr:      false,
		},
		{
			name:         "YouTube Shorts Native URL",
			url:          "https://www.youtube.com/shorts/dQw4w9WgXcQ",
			wantPlatform: "youtube",
			wantID:       "dQw4w9WgXcQ",
			wantErr:      false,
		},
		{
			name:         "YouTube Live URL",
			url:          "https://www.youtube.com/live/dQw4w9WgXcQ",
			wantPlatform: "youtube",
			wantID:       "dQw4w9WgXcQ",
			wantErr:      false,
		},
		{
			name:         "YouTube Embed URL",
			url:          "https://www.youtube-nocookie.com/embed/dQw4w9WgXcQ",
			wantPlatform: "youtube",
			wantID:       "dQw4w9WgXcQ",
			wantErr:      false,
		},
		{
			name:         "Instagram Reel",
			url:          "https://www.instagram.com/reel/C123456789/",
			wantPlatform: "instagram",
			wantID:       "C123456789",
			wantErr:      false,
		},
		{
			name:         "Facebook Video",
			url:          "https://www.facebook.com/watch/?v=9876543210",
			wantPlatform: "facebook",
			wantID:       "9876543210",
			wantErr:      false,
		},
		{
			name:         "X Twitter Status",
			url:          "https://x.com/user/status/112233445566",
			wantPlatform: "x_twitter",
			wantID:       "112233445566",
			wantErr:      false,
		},
		{
			name:         "Unsupported Platform",
			url:          "https://vimeo.com/123456",
			wantPlatform: "",
			wantID:       "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, err := ParseVideoURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVideoURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if meta.Platform != tt.wantPlatform {
					t.Errorf("Platform = %v, want %v", meta.Platform, tt.wantPlatform)
				}
				if meta.ExternalVideoID != tt.wantID {
					t.Errorf("ExternalVideoID = %v, want %v", meta.ExternalVideoID, tt.wantID)
				}
				if meta.EmbedHTML == "" {
					t.Error("Expected non-empty EmbedHTML")
				}
			}
		})
	}
}
