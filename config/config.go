package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config holds all application settings loaded from environment.
type Config struct {
	BindAddr  string   // address for HTTP server binding
	WhiteList []string // allowed IPs for access control
	Interval  int      // minimum interval between SMS sends
	Period    int      // TOTP period
	Digits    int      // number of digits in TOTP code
	Algorithm string   // hash algorithm for TOTP
	Skew      int      // allowed clock skew in periods

	SMSC struct {
		Login    string `mapstructure:"login" env:"SMSC_LOGIN" validate:"required"`
		Password string `mapstructure:"password" env:"SMSC_PASSWORD" validate:"required"`
		URL      string `mapstructure:"url" env:"SMSC_URL" default:"https://smsc.ru/sys/send.php"`
	} `mapstructure:"smsc"`
}

// LoadConfig reads environment variables and returns Config.
// It falls back to sensible defaults if variables are not set.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		BindAddr:  getEnv("TOTP_BIND", ":8080"),
		WhiteList: splitEnv("TOTP_WHITELIST", ","),
		Interval:  getEnvInt("TOTP_INTERVAL", 30),
		Period:    getEnvInt("TOTP_PERIOD", 60),
		Digits:    getEnvInt("TOTP_DIGITS", 6),
		Algorithm: strings.ToUpper(getEnv("TOTP_ALGO", "SHA1")),
		Skew:      getEnvInt("TOTP_SKEW", 1),
	}
	return cfg, nil
}

// getEnv returns value or default.
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// splitEnv splits env var by sep or returns empty slice.
func splitEnv(key, sep string) []string {
	val := os.Getenv(key)
	if val == "" {
		return nil
	}
	parts := strings.Split(val, sep)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// getEnvInt parses int or returns default.
func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		var i int
		if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
			return i
		}
	}
	return def
}

// getEnvDuration parses seconds or returns default.
func getEnvDuration(key string, defSec int) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v + "s"); err == nil {
			return d
		}
	}
	return time.Duration(defSec) * time.Second
}
