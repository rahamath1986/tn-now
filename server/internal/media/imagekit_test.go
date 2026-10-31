package media

import (
	"context"
	"strings"
	"testing"

	"tn-now/server/internal/config"
)

func TestGetEffectiveImageKitConfig(t *testing.T) {
	baseCfg := &config.Config{
		ImageKitURLEndpoint: "https://ik.imagekit.io/tn24test",
		ImageKitPublicKey:   "public_test123",
		ImageKitPrivateKey:  "private_test456",
		ImageKitFolder:      "/news_articles",
	}

	cfg := GetEffectiveImageKitConfig(context.Background(), nil, baseCfg)
	if cfg.URLEndpoint != "https://ik.imagekit.io/tn24test" {
		t.Errorf("Expected endpoint %s, got %s", "https://ik.imagekit.io/tn24test", cfg.URLEndpoint)
	}
	if cfg.PublicKey != "public_test123" {
		t.Errorf("Expected public key %s, got %s", "public_test123", cfg.PublicKey)
	}
	if cfg.PrivateKey != "private_test456" {
		t.Errorf("Expected private key %s, got %s", "private_test456", cfg.PrivateKey)
	}
	if cfg.Folder != "/news_articles" {
		t.Errorf("Expected folder %s, got %s", "/news_articles", cfg.Folder)
	}
	if !cfg.AutoOptimize {
		t.Error("Expected AutoOptimize default to be true")
	}
}

func TestUploadValidation(t *testing.T) {
	ctx := context.Background()

	// Missing config
	_, err := UploadToImageKit(ctx, nil, []byte("data"), "test.jpg")
	if err == nil || !strings.Contains(err.Error(), "not configured") {
		t.Errorf("Expected not configured error, got: %v", err)
	}

	cfg := &ImageKitConfig{
		URLEndpoint: "https://ik.imagekit.io/demo",
		PrivateKey:  "private_key",
	}

	// Empty payload
	_, err = UploadToImageKit(ctx, cfg, []byte{}, "test.jpg")
	if err == nil || !strings.Contains(err.Error(), "empty") {
		t.Errorf("Expected empty payload error, got: %v", err)
	}

	// Payload over 10MB
	hugeFile := make([]byte, 11*1024*1024)
	_, err = UploadToImageKit(ctx, cfg, hugeFile, "large.jpg")
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Errorf("Expected size exceeds error, got: %v", err)
	}
}

func TestTestConnectionValidation(t *testing.T) {
	ctx := context.Background()

	err := TestImageKitConnection(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil config")
	}

	err = TestImageKitConnection(ctx, &ImageKitConfig{PrivateKey: ""})
	if err == nil || !strings.Contains(err.Error(), "Private Key") {
		t.Errorf("Expected Private Key error, got: %v", err)
	}

	err = TestImageKitConnection(ctx, &ImageKitConfig{PrivateKey: "priv", URLEndpoint: ""})
	if err == nil || !strings.Contains(err.Error(), "URL Endpoint") {
		t.Errorf("Expected URL Endpoint error, got: %v", err)
	}
}

