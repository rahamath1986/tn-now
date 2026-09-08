package config

import (
	"os"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	// Set custom env variables
	os.Setenv("PORT", "9999")
	os.Setenv("APP_ENV", "testing")
	os.Setenv("DATABASE_URL", "postgres://test_user:pass@localhost:5432/test_db")
	os.Setenv("ALLOWED_VIDEO_DOMAINS", "domain1.com, domain2.com")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("APP_ENV")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("ALLOWED_VIDEO_DOMAINS")
	}()

	cfg := Load()

	if cfg.Port != "9999" {
		t.Errorf("Expected Port to be 9999, got %s", cfg.Port)
	}

	if cfg.AppEnv != "testing" {
		t.Errorf("Expected AppEnv to be testing, got %s", cfg.AppEnv)
	}

	if cfg.DatabaseURL != "postgres://test_user:pass@localhost:5432/test_db" {
		t.Errorf("Expected DatabaseURL to match, got %s", cfg.DatabaseURL)
	}

	if len(cfg.AllowedVideoDomains) != 2 || cfg.AllowedVideoDomains[0] != "domain1.com" || cfg.AllowedVideoDomains[1] != "domain2.com" {
		t.Errorf("Expected AllowedVideoDomains to be [domain1.com, domain2.com], got %v", cfg.AllowedVideoDomains)
	}
}

func TestConfigDefaults(t *testing.T) {
	// Clear env variables
	os.Unsetenv("PORT")
	os.Unsetenv("APP_ENV")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Expected default Port to be 8080, got %s", cfg.Port)
	}

	if cfg.AppEnv != "development" {
		t.Errorf("Expected default AppEnv to be development, got %s", cfg.AppEnv)
	}
}
