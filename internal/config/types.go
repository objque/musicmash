package config

type AppConfig struct {
	HTTP     HTTPConfig     `yaml:"http"`
	DB       DBConfig       `yaml:"db"`
	Log      LogConfig      `yaml:"log"`
	Fetcher  FetcherConfig  `yaml:"fetcher"`
	Stores   StoresConfig   `yaml:"stores"`
	Sentry   SentryConfig   `yaml:"sentry"`
	Notifier NotifierConfig `yaml:"notifier"`
	Proxy    ProxyConfig    `yaml:"proxy"`
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
	Type          string `yaml:"type"`
	Host          string `yaml:"host"`
	Name          string `yaml:"name"`
	Login         string `yaml:"login"`
	Pass          string `yaml:"pass"`
	Log           bool   `yaml:"log"`
	AutoMigrate   bool   `yaml:"auto_migrate"`
	MigrationsDir string `yaml:"migrations_dir"`
}

type FetcherConfig struct {
	Enabled bool    `yaml:"enabled"`
	Delay   float64 `yaml:"delay"`
}

type NotifierConfig struct {
	Enabled       bool    `yaml:"enabled"`
	TelegramToken string  `yaml:"telegram_token"`
	Delay         float64 `yaml:"delay"`
}

type StoreConfig struct {
	URL          string `yaml:"url"`
	FetchWorkers int    `yaml:"fetch_workers"`
	SaveWorkers  int    `yaml:"save_workers"`
	Meta         Meta   `yaml:"meta"`
	ReleaseURL   string `yaml:"release_url"`
	Name         string `yaml:"name"`
	Fetch        bool   `yaml:"fetch"`
}
type StoresConfig map[string]*StoreConfig

type Meta map[string]string

type SentryConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Key         string `yaml:"key"`
	Environment string `yaml:"environment"`
}

type ProxyConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
}
