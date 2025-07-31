package config

import (
	"os"
	"testing"
)

func TestLoadConfigDefaults(t *testing.T) {
	os.Clearenv()
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if cfg.BindAddr != ":8080" {
		t.Errorf("BindAddr = %q; want \":8080\"", cfg.BindAddr)
	}
	if len(cfg.WhiteList) != 0 {
		t.Errorf("WhiteList = %v; want empty list", cfg.WhiteList)
	}
	if cfg.Interval != 30 {
		t.Errorf("Interval = %v; want 30", cfg.Interval)
	}
	if cfg.Period != 60 {
		t.Errorf("Period = %v; want 60", cfg.Period)
	}
	if cfg.Digits != 6 {
		t.Errorf("Digits = %d; want 6", cfg.Digits)
	}
	if cfg.Algorithm != "SHA1" {
		t.Errorf("Algorithm = %q; want \"SHA1\"", cfg.Algorithm)
	}
	if cfg.Skew != 1 {
		t.Errorf("Skew = %d; want 1", cfg.Skew)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("TOTP_BIND_ADDR", ":9090")
	os.Setenv("TOTP_INTERVAL", "10")
	os.Setenv("TOTP_PERIOD", "45")
	os.Setenv("TOTP_DIGITS", "8")
	os.Setenv("TOTP_ALGORITHM", "SHA256")
	os.Setenv("TOTP_SKEW", "3")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if cfg.BindAddr != ":9090" {
		t.Errorf("BindAddr = %q; want \":9090\"", cfg.BindAddr)
	}
	if cfg.Interval != 10 {
		t.Errorf("Interval = %v; want 10", cfg.Interval)
	}
	if cfg.Period != 45 {
		t.Errorf("Period = %v; want 45", cfg.Period)
	}
	if cfg.Digits != 8 {
		t.Errorf("Digits = %d; want 8", cfg.Digits)
	}
	if cfg.Algorithm != "SHA256" {
		t.Errorf("Algorithm = %q; want \"SHA256\"", cfg.Algorithm)
	}
	if cfg.Skew != 3 {
		t.Errorf("Skew = %d; want 3", cfg.Skew)
	}
}
