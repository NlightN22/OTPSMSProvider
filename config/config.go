package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application settings loaded from environment.
type Config struct {
	BindAddr   string   `mapstructure:"bind_addr" default:":8080"`                                    // address for HTTP server binding
	WhiteList  []string `mapstructure:"white_list"`                                                   // allowed IPs for access control
	Interval   int      `mapstructure:"interval" validate:"gte=0" default:"30"`                       // minimum interval between SMS sends
	Period     int      `mapstructure:"period"  validate:"gt=0" default:"60"`                         // TOTP period
	Digits     int      `mapstructure:"digits"  validate:"gt=0,lte=10" default:"6"`                   // number of digits in TOTP code
	Algorithm  string   `mapstructure:"algorithm" validate:"oneof=SHA1 SHA256 SHA512" default:"SHA1"` // hash algorithm for TOTP
	Skew       int      `mapstructure:"skew" default:"1"`                                             // allowed clock skew in periods
	Debug      bool     `mapstructure:"debug" env:"TOTP_DEBUG"`
	LogLevel   string   `mapstructure:"log_level" env:"TOTP_LOG_LEVEL" default:"info"`
	PrefixText string   `mapstructure:"prefix_text" env:"TOTP_LOG_LEVEL" default:"Your code is:"`

	SMSC struct {
		Login    string `mapstructure:"login"  validate:"required"`
		Password string `mapstructure:"password"  validate:"required"`
	} `mapstructure:"smsc"`
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// LoadConfig reads environment variables and returns Config.
// It falls back to sensible defaults if variables are not set.
func LoadConfig() (*Config, error) {

	v := viper.New()

	v.SetDefault("bind_addr", ":8080")
	v.SetDefault("interval", 30)
	v.SetDefault("period", 60)
	v.SetDefault("digits", 6)
	v.SetDefault("algorithm", "SHA1")
	v.SetDefault("skew", 1)

	v.SetEnvPrefix("TOTP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var usedFile string
	if fileExists("config.yaml") {
		usedFile = "config.yaml"
	}

	if usedFile != "" {
		v.SetConfigFile(usedFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config file %q: %w", usedFile, err)
		}
	} else {
		v.AutomaticEnv()
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
