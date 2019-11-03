package config

type AppConfig struct {
	HTTP     HTTPConfig     `yaml:"http"`
	DB       DBConfig       `yaml:"db"`
	Log      LogConfig      `yaml:"log"`
	Fetching FetchingConfig `yaml:"fetching"`
	Stores   StoresConfig   `yaml:"stores"`
	Sentry   SentryConfig   `yaml:"sentry"`
	Notifier NotifierConfig `yaml:"notifier"`
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
	Type  string `yaml:"type"`
	Host  string `yaml:"host"`
	Name  string `yaml:"name"`
	Login string `yaml:"login"`
	Pass  string `yaml:"pass"`
	Log   bool   `yaml:"log"`
}

type FetchingConfig struct {
	RefetchAfterHours   float64 `yaml:"refetch_after_hours"`
	CountOfSkippedHours float64 `yaml:"count_of_skipped_hours"`
}

type NotifierConfig struct {
	TelegramToken       string  `yaml:"telegram_token"`
	CountOfSkippedHours float64 `yaml:"count_of_skipped_hours"`
}

type StoreConfig struct {
	URL          string `yaml:"url"`
	FetchWorkers int    `yaml:"fetch_workers"`
	Meta         Meta   `yaml:"meta"`
	ReleaseURL   string `yaml:"release_url"`
	ArtistURL    string `yaml:"artist_url"`
	Name         string `yaml:"name"`
	Fetch        bool   `json:"fetch"`
}
type StoresConfig map[string]*StoreConfig

type Meta map[string]string

type SentryConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Key         string `yaml:"key"`
	Environment string `yaml:"environment"`
}
