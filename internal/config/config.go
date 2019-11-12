package config

import (
	"flag"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var Config *AppConfig

func New() *AppConfig {
	return &AppConfig{
		HTTP: HTTPConfig{
			IP:   "127.0.0.1",
			Port: 8844,
		},
		DB: DBConfig{
			Type: "sqlite3",
			Host: "./musicmash.sqlite3",
			Log:  false,
		},
		Log: LogConfig{
			File:  "./musicmash.log",
			Level: "INFO",
		},
		Fetcher: FetcherConfig{
			Enabled:           true,
			RefetchAfterHours: 1,
			Delay:             8,
		},
		Stores: StoresConfig{
			"itunes": {
				Name:         "Apple Music",
				URL:          "https://api.music.apple.com",
				FetchWorkers: 5,
				ReleaseURL:   "https://itunes.apple.com/us/album/%s",
				Fetch:        true,
			},
		},
		Sentry: SentryConfig{
			Enabled:     false,
			Key:         "https://uuid@sentry.io/123456",
			Environment: "production",
		},
		Notifier: NotifierConfig{
			Enabled:       true,
			TelegramToken: "12345:xxxx_yyy_token",
			Delay:         1,
		},
	}
}

func (c *AppConfig) LoadFromFile(configPath string) error {
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	return c.LoadFromBytes(b)
}

func (c *AppConfig) LoadFromBytes(val []byte) error {
	if err := yaml.Unmarshal(val, c); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (c *AppConfig) FlagSet() {
	flag.StringVar(&c.HTTP.IP, "http-ip", c.HTTP.IP, "API ip address")
	flag.IntVar(&c.HTTP.Port, "http-port", c.HTTP.Port, "API port")

	flag.StringVar(&c.DB.Type, "db-type", c.DB.Type, "Database type: mysql or sqlite3")
	flag.StringVar(&c.DB.Host, "db-host", c.DB.Host, "Database host")
	flag.StringVar(&c.DB.Name, "db-name", c.DB.Name, "Database name")
	flag.StringVar(&c.DB.Login, "db-login", c.DB.Login, "Database user login")
	flag.StringVar(&c.DB.Pass, "db-pass", c.DB.Pass, "Database user password")
	flag.BoolVar(&c.DB.Log, "db-log", c.DB.Log, "Echo database queries")

	flag.StringVar(&c.Log.Level, "log-level", c.Log.Level, "log level")
	flag.StringVar(&c.Log.File, "log-file", c.Log.File, "path to log file")

	flag.BoolVar(&c.Fetcher.Enabled, "fetcher-enabled", c.Fetcher.Enabled, "Is fetcher enabled")
	flag.Float64Var(&c.Fetcher.Delay, "fetcher-delay", c.Fetcher.Delay, "Delay between fetches")

	flag.BoolVar(&c.Sentry.Enabled, "sentry-enabled", c.Sentry.Enabled, "Is Sentry enabled")
	flag.StringVar(&c.Sentry.Key, "sentry-key", c.Sentry.Key, "Sentry dsn")
	flag.StringVar(&c.Sentry.Environment, "sentry-environment", c.Sentry.Environment, "Sentry environment")

	flag.BoolVar(&c.Notifier.Enabled, "notifier-enabled", c.Notifier.Enabled, "Is notifier enabled")
	flag.Float64Var(&c.Notifier.Delay, "notifier-delay", c.Notifier.Delay, "Delay between notifies")
	flag.StringVar(&c.Notifier.TelegramToken, "notifier-telegram-token", c.Notifier.TelegramToken, "Telegram bot token")
}

func (c *AppConfig) FlagReload() {
	_ = flag.Set("http-port", fmt.Sprintf("%d", c.HTTP.Port))

	_ = flag.Set("db-type", c.DB.Type)
	_ = flag.Set("db-host", c.DB.Host)
	_ = flag.Set("db-name", c.DB.Name)
	_ = flag.Set("db-login", c.DB.Login)
	_ = flag.Set("db-pass", c.DB.Pass)
	_ = flag.Set("db-log", fmt.Sprintf("%t", c.DB.Log))

	_ = flag.Set("log-level", c.Log.Level)
	_ = flag.Set("log-file", c.Log.File)

	_ = flag.Set("fetcher-enabled", fmt.Sprintf("%t", c.Fetcher.Enabled))
	_ = flag.Set("fetcher-delay", fmt.Sprintf("%v", c.Fetcher.Delay))

	_ = flag.Set("sentry", fmt.Sprintf("%t", c.Sentry.Enabled))
	_ = flag.Set("sentry-key", c.Sentry.Key)
	_ = flag.Set("sentry-environment", c.Sentry.Environment)

	_ = flag.Set("notifier-enabled", fmt.Sprintf("%t", c.Notifier.Enabled))
	_ = flag.Set("notifier-delay", fmt.Sprintf("%v", c.Notifier.Delay))
	_ = flag.Set("notifier-telegram-token", c.Notifier.TelegramToken)
}

func (db *DBConfig) GetConnString() (dialect, connString string) {
	const mysqlConnectionFormat = "%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=UTC"
	switch db.Type {
	case "mysql":
		return db.Type, fmt.Sprintf(mysqlConnectionFormat, db.Login, db.Pass, db.Host, db.Name)
	case "sqlite3":
		return db.Type, db.Host
	default:
		panic("Only mysql/sqlite3 is currently supported")
	}
}
