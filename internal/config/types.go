package config

import "time"

type AppConfig struct {
	HTTP     HTTPConfig   `yaml:"http"`
	DB       DBConfig     `yaml:"db"`
	Log      LogConfig    `yaml:"log"`
	Notifier NotifyConfig `yaml:"notifier"`
	Sentry   SentryConfig `yaml:"sentry"`
}

type HTTPConfig struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type LogConfig struct {
	File  string `yaml:"file"`
	Level string `yaml:"level"`
}

type DBConfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	Name          string `yaml:"name"`
	Login         string `yaml:"login"`
	Pass          string `yaml:"pass"`
	Log           bool   `yaml:"log"`
	AutoMigrate   bool   `yaml:"auto_migrate"`
	MigrationsDir string `yaml:"migrations_dir"`
}

type NotifyConfig struct {
	Enabled bool          `yaml:"enabled"`
	Delay   time.Duration `yaml:"delay"`
	URL     string        `yaml:"url"`
}

type SentryConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Key         string `yaml:"key"`
	Environment string `yaml:"environment"`
}
