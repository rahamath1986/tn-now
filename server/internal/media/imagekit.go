package media

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"tn-now/server/internal/config"
)

type ImageKitConfig struct {
	URLEndpoint  string `json:"urlEndpoint"`
	PublicKey    string `json:"publicKey"`
	PrivateKey   string `json:"privateKey"`
	Folder       string `json:"folder"`
	AutoOptimize bool   `json:"autoOptimize"`
}

type ImageKitUploadResponse struct {
	FileID       string `json:"fileId"`
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	FilePath     string `json:"filePath"`
	URL          string `json:"url"`
	FileType     string `json:"fileType"`
	Height       int    `json:"height"`
	Width        int    `json:"width"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type UploadResult struct {
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
	FileID       string `json:"fileId"`
	Name         string `json:"name"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Size         int64  `json:"size"`
}

// GetEffectiveImageKitConfig retrieves ImageKit credentials from system_settings DB table, falling back to config
func GetEffectiveImageKitConfig(ctx context.Context, conn *pgxpool.Pool, baseCfg *config.Config) *ImageKitConfig {
	cfg := &ImageKitConfig{
		Folder:       "/articles",
		AutoOptimize: true,
	}

	if baseCfg != nil {
		cfg.URLEndpoint = baseCfg.ImageKitURLEndpoint
		cfg.PublicKey = baseCfg.ImageKitPublicKey
		cfg.PrivateKey = baseCfg.ImageKitPrivateKey
		if baseCfg.ImageKitFolder != "" {
			cfg.Folder = baseCfg.ImageKitFolder
		}
	}

	if conn == nil {
		return cfg
	}

	var endpoint, pubKey, privKey, folder, autoOpt string
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'imagekit_url_endpoint'").Scan(&endpoint)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'imagekit_public_key'").Scan(&pubKey)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'imagekit_private_key'").Scan(&privKey)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'imagekit_folder'").Scan(&folder)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'imagekit_auto_optimize'").Scan(&autoOpt)

	if strings.TrimSpace(endpoint) != "" {
		cfg.URLEndpoint = strings.TrimSpace(endpoint)
	}
	if strings.TrimSpace(pubKey) != "" {
		cfg.PublicKey = strings.TrimSpace(pubKey)
	}
	if strings.TrimSpace(privKey) != "" {
		cfg.PrivateKey = strings.TrimSpace(privKey)
	}
	if strings.TrimSpace(folder) != "" {
		cfg.Folder = strings.TrimSpace(folder)
	}
	if strings.TrimSpace(autoOpt) != "" {
		cfg.AutoOptimize = (autoOpt == "true" || autoOpt == "1")
	}

	return cfg
}

// UploadToImageKit uploads an image to ImageKit.io via their REST API
func UploadToImageKit(ctx context.Context, cfg *ImageKitConfig, fileBytes []byte, originalFilename string) (*UploadResult, error) {
	if cfg == nil || strings.TrimSpace(cfg.PrivateKey) == "" || strings.TrimSpace(cfg.URLEndpoint) == "" {
		return nil, errors.New("ImageKit is not configured. Please add your credentials in the Admin Panel.")
	}

	if len(fileBytes) == 0 {
		return nil, errors.New("empty file payload")
	}

	// Validate maximum size (10 MB limit)
	if len(fileBytes) > 10*1024*1024 {
		return nil, errors.New("image size exceeds 10MB limit")
	}

	ext := strings.ToLower(filepath.Ext(originalFilename))
	if ext == "" {
		ext = ".jpg"
	}
	cleanName := fmt.Sprintf("article_%d%s", time.Now().UnixNano()/1e6, ext)

	// Prepare multipart form body
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// File part
	part, err := writer.CreateFormFile("file", cleanName)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(fileBytes); err != nil {
		return nil, err
	}

	// fileName field
	_ = writer.WriteField("fileName", cleanName)

	// folder field
	folder := cfg.Folder
	if folder == "" {
		folder = "/articles"
	}
	_ = writer.WriteField("folder", folder)
	_ = writer.WriteField("useUniqueFileName", "true")
	_ = writer.WriteField("tags", "tn24,article,citizen_journalism")

	if err := writer.Close(); err != nil {
		return nil, err
	}

	// Create HTTP request to ImageKit Upload API
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://upload.imagekit.io/api/v1/files/upload", &body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Basic Auth using PrivateKey as username with blank password
	auth := base64.StdEncoding.EncodeToString([]byte(cfg.PrivateKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ImageKit request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp struct {
			Message string `json:"message"`
			Help    string `json:"help"`
		}
		_ = json.Unmarshal(respBody, &errResp)
		msg := errResp.Message
		if msg == "" {
			msg = string(respBody)
		}
		return nil, fmt.Errorf("ImageKit upload error (%d): %s", resp.StatusCode, msg)
	}

	var ikResp ImageKitUploadResponse
	if err := json.Unmarshal(respBody, &ikResp); err != nil {
		return nil, fmt.Errorf("failed to parse ImageKit response: %w", err)
	}

	finalURL := ikResp.URL
	// If auto-optimize is enabled, inject ImageKit real-time WebP & responsive width transformation
	if cfg.AutoOptimize && !strings.Contains(finalURL, "tr=") {
		endpoint := strings.TrimRight(cfg.URLEndpoint, "/")
		if strings.HasPrefix(finalURL, endpoint) {
			relPath := strings.TrimPrefix(finalURL, endpoint)
			finalURL = fmt.Sprintf("%s/tr:w-1200,q-85,f-webp%s", endpoint, relPath)
		}
	}

	return &UploadResult{
		URL:          finalURL,
		ThumbnailURL: ikResp.ThumbnailURL,
		FileID:       ikResp.FileID,
		Name:         ikResp.Name,
		Width:        ikResp.Width,
		Height:       ikResp.Height,
		Size:         ikResp.Size,
	}, nil
}

// TestImageKitConnection verifies that the PrivateKey and URLEndpoint can successfully communicate with ImageKit API
func TestImageKitConnection(ctx context.Context, cfg *ImageKitConfig) error {
	if cfg == nil || strings.TrimSpace(cfg.PrivateKey) == "" {
		return errors.New("Private Key is missing")
	}
	if strings.TrimSpace(cfg.URLEndpoint) == "" {
		return errors.New("URL Endpoint is missing")
	}

	// 1x1 transparent GIF pixel for ping test
	testPixel := []byte{
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00, 0x01, 0x00,
		0x80, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x21,
		0xf9, 0x04, 0x01, 0x00, 0x00, 0x00, 0x00, 0x2c, 0x00, 0x00,
		0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x02, 0x44,
		0x01, 0x00, 0x3b,
	}

	res, err := UploadToImageKit(ctx, cfg, testPixel, "connection_test.gif")
	if err != nil {
		return err
	}

	// Clean up the test file if possible
	if res != nil && res.FileID != "" {
		go func(fileID, privKey string) {
			delCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			delReq, _ := http.NewRequestWithContext(delCtx, http.MethodDelete, "https://api.imagekit.io/v1/files/"+fileID, nil)
			if delReq != nil {
				auth := base64.StdEncoding.EncodeToString([]byte(privKey + ":"))
				delReq.Header.Set("Authorization", "Basic "+auth)
				_, _ = http.DefaultClient.Do(delReq)
			}
		}(res.FileID, cfg.PrivateKey)
	}

	return nil
}

