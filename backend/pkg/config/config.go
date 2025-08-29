package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type LogConfig struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

type JWTConfig struct {
	Issuer   string `yaml:"issuer"`
	Audience string `yaml:"audience"`
	Alg      string `yaml:"alg"`
	Secret   string `yaml:"secret"` // 推荐用环境变量 JWT_SECRET 覆盖
}

type Config struct {
	HTTPAddr string    `yaml:"http_addr"`
	JWT      JWTConfig `yaml:"jwt"`
	Log      LogConfig `yaml:"log"`
}

// Load reads config/app.yaml then overlays limited secrets from env.
func Load() Config {
	cfg := Config{
		HTTPAddr: ":8080",
		JWT: JWTConfig{
			Issuer:   "neoneuro",
			Audience: "neoneuro-client",
			Alg:      "HS256",
		},
		Log: LogConfig{
			Level:      "info",
			File:       "log/app.log",
			MaxSize:    100,
			MaxBackups: 7,
			MaxAge:     14,
			Compress:   true,
		},
	}

	path := "config/app.yaml"
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		path = p
	}
	if b, err := os.ReadFile(path); err == nil {
		_ = yaml.Unmarshal(b, &cfg)
	}

	// Overlay secrets via env (limited)
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}

	// Normalize relative paths to be under current working dir
	if !filepath.IsAbs(cfg.Log.File) {
		cfg.Log.File = filepath.Clean(cfg.Log.File)
	}
	return cfg
}
